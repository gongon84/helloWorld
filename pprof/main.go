package main

import (
	"net/http"
	"net/http/pprof"
	_ "net/http/pprof"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	pprofGroup := e.Group("/debug/pprof")
	pprofGroup.Any("/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
	pprofGroup.Any("/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
	pprofGroup.Any("/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
	pprofGroup.Any("/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
	pprofGroup.Any("/*", echo.WrapHandler(http.HandlerFunc(pprof.Index)))

	e.GET("/", calcFib)

	e.Logger.Fatal(e.Start(":8080"))
}

func calcFib(c echo.Context) error {
	for i := 0; i < 100; i++ {
		fib(i)
	}
	return c.String(http.StatusOK, "Hello, World!")
}

func fib(n int) int {
	if n < 2 {
		return n
	}
	return fib(n-1) + fib(n-2)
}
