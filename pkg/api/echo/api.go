package echo

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/labstack/gommon/log"
	echoSwagger "github.com/swaggo/echo-swagger"
)

type RouteHandler interface {
	Bind(e *echo.Echo)
}

type App struct {
	echo *echo.Echo
	Addr string
}

func NewApp(echo *echo.Echo, addr string) *App {
	a := App{
		echo,
		addr,
	}

	a.configureEcho()

	return &a
}

func (a *App) BindRoutes(h RouteHandler) {
	h.Bind(a.echo)
}

func (a *App) configureEcho() {
	e := a.echo

	e.Logger.SetLevel(echolog.DEBUG)
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/", Redirect("/swagger/index.html"))
	e.GET("/swagger*", echoSwagger.WrapHandler)
}

// Serve start listening at addr
// @title           Stock-Data API
// @description     api to display stock data of imported exchanges and symbols
// @BasePath  /
func (a *App) Serve() error {
	return a.echo.Start(a.Addr)
}

func Redirect(url string) func(ctx echo.Context) error {
	return func(ctx echo.Context) error { return ctx.Redirect(http.StatusPermanentRedirect, url) }
}
