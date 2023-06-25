package libecho

import (
	"ais/lib/libresponse"
	"context"
	"encoding/json"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/helloferdie/golib/liblogger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Initialize - Initialize required middleware before assign route
func Initialize(e *echo.Echo) {
	e.HTTPErrorHandler = ErrorHandler
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(Logger)
}

// StartHTTP - Start as HTTP server
func StartHTTP(e *echo.Echo) {
	// Start server
	go func() {
		if err := e.Start(":" + os.Getenv("app_port")); err != nil {
			liblogger.Log(nil, false).Error(err)
			liblogger.Log(nil, false).Error("Fail start HTTP server")
			liblogger.Log(nil, false).Error("Shutting down the server")
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		liblogger.Log(nil, false).Error("Fail shutting down server")
		os.Exit(1)
	} else {
		liblogger.Log(nil, false).Info("Shutdown HTTP server - done")
	}
}

// Logger - Log every incoming request
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := map[string]interface{}{
			"method":   c.Request().Method,
			"uri":      c.Request().URL.String(),
			"ip":       c.Request().RemoteAddr,
			"ip_real":  c.Request().Header.Get("X-Real-Ip"),
			"ip_proxy": c.Request().Header.Get("X-Proxy-Ip"),
		}
		liblogger.Log(log, false).Info("incoming request")
		return next(c)
	}
}

// ErrorHandler - Custom error handler
func ErrorHandler(err error, c echo.Context) {
	resp := libresponse.GetDefault()
	report, ok := err.(*echo.HTTPError)
	if !ok {
		report = echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	resp.Code = report.Code
	switch resp.Code {
	case 400:
		resp.Message = "common.error.request.default"
		reportMsg := report.Message.(string)
		if reportMsg == "missing or malformed jwt" {
			resp.Error = "common.error.request.jwt.required"
		} else {
			if e1, ok := report.Internal.(*json.UnmarshalTypeError); ok {
				resp.Error = "common.error.request.unmarshal_var"
				resp.ErrorVar = map[string]interface{}{
					"f": e1.Field,
					"e": e1.Type,
					"v": e1.Value,
				}
			}
		}

		if resp.Error == "" {
			resp.Error = "common.error.request.bad"
		}
	case 401:
		resp.Message = "common.error.request.default"
		reportMsg := report.Message.(string)
		if reportMsg == "invalid or expired jwt" {
			reportInternal := report.Internal.Error()
			if reportInternal == "Token is expired" {
				resp.Error = "common.error.request.jwt.expired"
			} else {
				resp.Error = "common.error.request.jwt.invalid"
			}
		}

		if resp.Error == "" {
			resp.Error = "common.error.request.unauthorized"
		}
	case 403:
		resp.Message = "common.error.request.default"
		resp.Error = "common.error.request.forbidden"
	case 404:
		resp.Message = "common.error.request.default"
		if resp.Error == "" {
			resp.Error = "common.error.request.route.not_found"
		}
	case 405:
		resp.Message = "common.error.request.default"
		resp.Error = "common.error.request.method"
	case 413:
		resp.Message = "common.error.request.default"
		resp.Error = "common.error.request.size.large"
	case 415:
		resp.Message = "common.error.request.default"
		resp.Error = "common.error.request.media_type"
	case 422:
		resp.Message = "validation.error.default"
	case 500:
		resp.Message = "common.error.server.internal"
		if resp.Error == "" {
			resp.Error = "common.error.server.internal"
		}
	case 502:
		resp.Message = "common.error.server.gateway"
		resp.Error = "common.error.service.unreachable"
	}

	ParseResponse(c, resp)
}

// ParseResponse - Parse return response in JSON format
func ParseResponse(c echo.Context, resp *libresponse.Response) (err error) {
	return c.JSON(resp.Code, resp)
}
