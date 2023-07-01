package app

import (
	"log"
	"minly-backend/internal/app/database"
	"minly-backend/internal/app/endpoint"
	"minly-backend/internal/app/service"

	"github.com/labstack/echo/v4"
)

type App struct {
	e    *endpoint.Endpoint
	s    *service.Service
	db   *database.DataBase
	echo *echo.Echo
}

func New() (*App, error) {
	a := &App{}

	a.s = service.New()

	DataBase, err := database.New("minly", "links")
	if err != nil {
		panic(err)
	}

	a.db = DataBase

	a.e = endpoint.New(a.s, a.db)

	a.echo = echo.New()

	a.echo.GET("/getResult", a.e.GetResult)
	a.echo.GET("/getLink", a.e.GetLink)

	return a, nil

}

func (a *App) Run() error {

	err := a.echo.Start(":8080")
	if err != nil {
		log.Fatal(err)
	}

	return nil
}
