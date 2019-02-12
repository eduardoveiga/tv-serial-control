/*
 * Copyright (C) 2019
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: MIT
 */

package lg

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

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
	UnknownCommandErr    = errors.New("Unknown command")
)

func (c Command) Send(port *serial.Port, args ...interface{}) (map[string]interface{}, error) {
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

	return map[string]interface{}{
		"cmd2": cmd2,
		"id":   id,
		"ack":  result[0][3],
		"data": data,
	}, nil
}
