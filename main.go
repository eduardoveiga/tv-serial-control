package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gustavosbarreto/tv-control/driverapi"
	_ "github.com/gustavosbarreto/tv-control/drivers/lg"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

type ConfigOptions struct {
	Driver string `envconfig:"driver" required:"true"`
	Device string `envconfig:"device" required:"true"`
	Port   int    `envconfig:"port" default:"8080"`
}

func main() {
	e := echo.New()

	opts := ConfigOptions{}

	err := envconfig.Process("", &opts)
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	driver := driverapi.GetDriver(opts.Driver)
	if driver == nil {
		logrus.Panic("Driver not found")
	}

	if err := driver.Initialize(opts.Device); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Panic("Failed to initialize driver")
	}

	logrus.WithFields(logrus.Fields{
		"driver":   opts.Driver,
		"commands": driver.AvailableCommands(),
	}).Info("Driver loaded")

	e.GET("/driver", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"driver":   opts.Driver,
			"device":   opts.Device,
			"commands": driver.AvailableCommands(),
		})
	})

	e.POST("/commands/:cmd", func(c echo.Context) error {
		cmd := c.Param("cmd")

		var req struct {
			Args []interface{} `json:"args,omitempty"`
		}

		if c.Request().ContentLength > 0 {
			err := c.Bind(&req)
			if err != nil {
				return err
			}
		}

		res, err := driver.SendCommand(cmd, req.Args...)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, res)
	})

	log.Fatal(e.Start(fmt.Sprintf(":%d", opts.Port)))
}
