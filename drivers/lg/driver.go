/*
 * Copyright (C) 2019
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: MIT
 */

package lg

import (
	"time"

	"github.com/OSSystems/tv-serial-control/driverapi"
	"github.com/tarm/serial"
)

func init() {
	driverapi.RegisterDriver("lg", &driver{})
}

type driver struct {
	port *serial.Port
}

func (l *driver) Initialize(device string) error {
	port, err := serial.OpenPort(&serial.Config{
		Name:        device,
		Baud:        9600,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    serial.Stop1,
		ReadTimeout: time.Second,
	})

	l.port = port

	return err
}

func (l *driver) AvailableCommands() []string {
	keys := make([]string, 0, len(cmds))

	for cmd := range cmds {
		keys = append(keys, string(cmd))
	}

	return keys
}

func (l *driver) SendCommand(name string, args ...interface{}) (map[string]interface{}, error) {
	if _, ok := cmds[Command(name)]; ok {
		res, err := Command(name).Send(l.port, args...)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	return nil, UnknownCommandErr
}
