/*
 * Copyright (C) 2019
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: MIT
 */

package dummy

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/OSSystems/tv-serial-control/driverapi"
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
)

var (
	MalformattedErr      = errors.New("Malformatted raw response")
	ArgumentsMismatchErr = errors.New("Arguments mismatch")
	InvalidResponseErr   = errors.New("Invalid response")
	UnknownCommandErr    = errors.New("Unknown command")
)

func (c Command) Send(tv *driverapi.TV, args ...interface{}) (map[string]interface{}, error) {

	command_string := fmt.Sprintf(string(cmds[c]) + "\r")

	time.Sleep(time.Second)

	var response string

	re, err := regexp.Compile(`(\w)(\w) ([0-9a-fA-F]{2}) ([0-9a-fA-F]{2})\r`)

	if err != nil {
		return nil, err
	}

	fmt.Printf("command: %s \n", command_string)

	result := re.FindAllStringSubmatch(command_string, -1)

	if len(result) == 0 {
		return nil, MalformattedErr
	}

	if len(result[0]) != 5 {
		return nil, ArgumentsMismatchErr
	}
	cmd1 := result[0][1]
	cmd2 := result[0][2]

	id, err := strconv.ParseInt(result[0][3], 16, 0)
	if err != nil {
		return nil, err
	}

	data, err := strconv.ParseInt(result[0][4], 16, 0)
	if err != nil {
		return nil, err
	}
	var data_return int64

	if data != 255 {
		data_return = data
		tv.SetData(fmt.Sprintf("%s%s", cmd1, cmd2), data_return)

	} else {
		data_return = tv.GetData(fmt.Sprintf("%s%s", cmd1, cmd2))
	}

	response = fmt.Sprintf("%s %02x OK%02xx", cmd2, id, data_return)
	fmt.Printf("ack:%s\n", response)

	return map[string]interface{}{
		"cmd2": cmd2,
		"id":   id,
		"ack":  "OK",
		"data": data_return,
	}, nil
}
