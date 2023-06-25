package handler

import (
	"ais/lib/libecho"
	"ais/pkg/api"
	"ais/pkg/article/request"
	"ais/pkg/article/service"

	"github.com/labstack/echo/v4"
)

// ArticleList - List articles
func ArticleList(c echo.Context) (err error) {
	r := new(request.List)
	if err = c.Bind(r); err != nil {
		return
	}

	format := api.GetFormatOutput(c)
	resp := service.List(r, format)
	return libecho.ParseResponse(c, resp)
}

// ArticleView - View article
func ArticleView(c echo.Context) (err error) {
	r := new(request.View)
	if err = c.Bind(r); err != nil {
		return
	}

	format := api.GetFormatOutput(c)
	resp := service.View(r, format)
	return libecho.ParseResponse(c, resp)
}

// ArticlePost - Post new article
func ArticlePost(c echo.Context) (err error) {
	r := new(request.Post)
	if err = c.Bind(r); err != nil {
		return
	}

	format := api.GetFormatOutput(c)
	resp := service.Post(r, format)
	return libecho.ParseResponse(c, resp)
}

// ArticleUpdate - Update article
func ArticleUpdate(c echo.Context) (err error) {
	r := new(request.Update)
	if err = c.Bind(r); err != nil {
		return
	}

	format := api.GetFormatOutput(c)
	resp := service.Update(r, format)
	return libecho.ParseResponse(c, resp)
}

// ArticleDelete - Delete article
func ArticleDelete(c echo.Context) (err error) {
	r := new(request.Delete)
	if err = c.Bind(r); err != nil {
		return
	}

	format := api.GetFormatOutput(c)
	resp := service.Delete(r, format)
	return libecho.ParseResponse(c, resp)
}
