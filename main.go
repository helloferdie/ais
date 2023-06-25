package main

import (
	"ais/handler"
	"ais/lib/libecho"
	"embed"
	"os"

	"github.com/helloferdie/golib/libdb"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/pressly/goose/v3"
)

//go:embed migration/*.sql
var embedMigration embed.FS

func init() {
	godotenv.Load()

	if os.Getenv("app_migrate") == "1" {
		// Perform database migration
		d, _ := libdb.Open("")
		defer d.Close()

		goose.SetBaseFS(embedMigration)
		if err := goose.SetDialect("mysql"); err != nil {
			panic(err)
		}

		if err := goose.Up(d.DB, "migration"); err != nil {
			panic(err)
		}
	}
}

func main() {
	e := echo.New()
	libecho.Initialize(e)

	e.GET("/articles", handler.ArticleList)
	e.GET("/articles/:id", handler.ArticleView)

	e.POST("/articles", handler.ArticlePost)
	e.POST("/articles/update", handler.ArticleUpdate)
	e.POST("/articles/delete", handler.ArticleDelete)

	libecho.StartHTTP(e)
}
