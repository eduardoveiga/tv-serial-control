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
	KeyHome:        "mc 00 7c",
	KeyUp:          "mc 00 40",
	KeyLeft:        "mc 00 07",
	KeyEnter:       "mc 00 7c",
	KeyRight:       "mc 00 06",
	KeyDown:        "mc 00 31",
	KeyEsc:         "mc 00 1b",
	KeyA:           "mc 00 41",
	KeyZ:           "mc 00 5a",
	KeyR:           "mc 00 52",
	KeyVolUp:       "mc 00 02",
	keyVolDown:     "mc 00 03",
}

const (
	PowerOnCmd     Command = "power_on"
	PowerOffCmd    Command = "power_off"
	PowerStatusCmd Command = "power_status"
	VolumeGetCmd   Command = "volume_get"
	VolumeSetCmd   Command = "volume_set"
	KeyHome        Command = "key_home"
	KeyUp          Command = "key_up"
	KeyLeft        Command = "key_left"
	KeyEnter       Command = "key_enter"
	KeyRight       Command = "key_right"
	KeyDown        Command = "key_down"
	KeyEsc         Command = "key_esc"
	KeyA           Command = "key_a"
	KeyZ           Command = "key_z"
	KeyR           Command = "key_r"
	KeyVolUp       Command = "volume_up"
	keyVolDown     Command = "volume_down"
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

	re, err := regexp.Compile(`(\w) ([0-9a-fA-F]{2}) (OK|NG)([0-9a-fA-F]{2})x`)
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

	id, err := strconv.ParseInt(result[0][2], 16, 0)
	if err != nil {
		return nil, err
	}

	data, err := strconv.ParseInt(result[0][4], 16, 0)
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
