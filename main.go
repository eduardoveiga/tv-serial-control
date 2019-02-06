package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gustavosbarreto/tv-control/lg"
	"github.com/labstack/echo"
)

func main() {
	e := echo.New()

	e.POST("/:cmd", func(c echo.Context) error {
		cmd := c.Param("cmd")

		d := lg.LG{}
		err := d.Initialize(os.Args[1])
		if err != nil {
			return err
		}

		var req struct {
			Args []interface{} `json:"args,omitempty"`
		}

		if c.Request().ContentLength > 0 {
			err := c.Bind(&req)
			if err != nil {
				return err
			}
		}

		res, err := d.SendCommand(cmd, req.Args...)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	})

	log.Fatal(e.Start(":8080"))
}
