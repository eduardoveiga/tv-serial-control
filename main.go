package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gustavosbarreto/tv-control/driverapi"
	_ "github.com/gustavosbarreto/tv-control/drivers/lg"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

func main() {
	e := echo.New()

	driver := driverapi.GetDriver("lg")
	if driver == nil {
		logrus.Panic("Driver not found")
	}

	if err := driver.Initialize(os.Args[1]); err != nil {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Panic("Failed to initialize driver")
	}

	logrus.WithFields(logrus.Fields{
		"driver":   "lg",
		"commands": driver.AvailableCommands(),
	}).Info("Driver loaded")

	e.GET("/driver", func(c echo.Context) error {
		return c.JSON(http.StatusOK, echo.Map{
			"driver":   "lg",
			"device":   os.Args[1],
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

	log.Fatal(e.Start(":8080"))
}
