package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/tarm/serial"
)

type Command string

var cmds = map[Command]string{
	PowerOnCmd:     "ka 00 01",
	PowerOffCmd:    "ka 00 00",
	PowerStatusCmd: "ka 00 ff",
	VolumeGetCmd:   "kf 00 ff",
	VolumeSetCmd:   "kf 00 %g",
}

const (
	PowerOnCmd     Command = "power_on"
	PowerOffCmd    Command = "power_off"
	PowerStatusCmd Command = "power_status"
	VolumeGetCmd   Command = "volume_get"
	VolumeSetCmd   Command = "volume_set"
)

var (
	MalformattedErr      = errors.New("Malformatted raw response")
	ArgumentsMismatchErr = errors.New("Arguments mismatch")
	InvalidResponseErr   = errors.New("Invalid response")
)

func (c Command) Send(port *serial.Port, args ...interface{}) (*Response, error) {
	fmt.Println(fmt.Sprintf(string(cmds[c])+"\r", args...))
	_, err := port.Write([]byte(fmt.Sprintf(string(cmds[c])+"\r", args...)))
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second)

	var buf bytes.Buffer

	for {
		c := make([]byte, 1)
		n, err := port.Read(c)
		if err == io.EOF {
			break
		}

		buf.WriteString(string(c[:n]))
	}

	re, err := regexp.Compile(`(\w) (\d+) (OK|NG)(\d+)x`)
	if err != nil {
		return nil, err
	}

	fmt.Println(buf.String())

	result := re.FindAllStringSubmatch(buf.String(), -1)
	if len(result) == 0 {
		return nil, MalformattedErr
	}

	if len(result[0]) != 5 {
		return nil, ArgumentsMismatchErr
	}

	cmd2 := result[0][1]

	if cmd2[0] != cmds[c][1] {
		return nil, InvalidResponseErr
	}

	id, err := strconv.Atoi(result[0][2])
	if err != nil {
		return nil, err
	}

	data, err := strconv.Atoi(result[0][4])
	if err != nil {
		return nil, err
	}

	return &Response{
		Cmd2: cmd2,
		ID:   id,
		ACK:  result[0][3],
		Data: data,
	}, nil
}

type Response struct {
	Cmd2 string `json:"cmd2"`
	ID   int    `json:"id"`
	ACK  string `json:"ack"`
	Data int    `json:"data"`
}

func main() {
	s, err := serial.OpenPort(&serial.Config{
		Name:        os.Args[1],
		Baud:        9600,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    serial.Stop1,
		ReadTimeout: time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.POST("/:cmd", func(c echo.Context) error {
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

		if _, ok := cmds[Command(cmd)]; ok {
			res, err := Command(cmd).Send(s, req.Args...)
			if err != nil {
				e.Logger.Error(err)
				return err
			}

			return c.JSON(http.StatusOK, res)
		}

		return echo.ErrBadRequest
	})

	log.Fatal(e.Start(":8080"))
}
