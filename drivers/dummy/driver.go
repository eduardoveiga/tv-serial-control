/*
 * Copyright (C) 2019
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: MIT
 */

package dummy

import (
	"errors"

	"github.com/OSSystems/tv-serial-control/driverapi"
)

var ErrCommandNotFound = errors.New("Command not found")

func init() {
	driverapi.RegisterDriver("dummy", &driver{})

}

type driver struct {
}

type TV struct {
}

func (d *driver) Initialize(device string) error {
	return nil
}

func (d *driver) AvailableCommands() []string {
	keys := make([]string, 0, len(cmds))

	for cmd := range cmds {
		keys = append(keys, string(cmd))
	}

	return keys

}



func (d *driver) SendCommand(name string, tv *driverapi.TV, args ...interface{}) (map[string]interface{}, error) {
	for _, cmd := range d.AvailableCommands() {
		if cmd == name {


			res, err := Command(name).Send(tv, args...)
			if err != nil {
				return nil, err
			}

			return res, nil

		}
	}

	return nil, ErrCommandNotFound
}

