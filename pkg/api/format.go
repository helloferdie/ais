package api

import (
	"ais/lib/libresponse"
	"time"

	"github.com/labstack/echo/v4"
)

// Header -
type Header struct {
	AcceptTimezone string
}

// GetFormatOutput -
func GetFormatOutput(c echo.Context) map[string]interface{} {
	tmp := new(Header)
	tmp.AcceptTimezone = c.Request().Header.Get("Accept-Timezone")

	m := libresponse.GetFormatOutput()
	m["header"] = tmp

	loc, err := time.LoadLocation(tmp.AcceptTimezone)
	if err == nil {
		m["timezone"] = loc.String()
	}
	return m
}
