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

func (d *driver) Initialize(device string) error {
	return nil
}

func (l *driver) AvailableCommands() []string {
	return []string{"echo"}
}

func (l *driver) SendCommand(name string, args ...interface{}) (map[string]interface{}, error) {
	for _, cmd := range l.AvailableCommands() {
		if cmd == name {
			return map[string]interface{}{
				"command": name,
				"args":    args,
			}, nil
		}
	}

	return nil, ErrCommandNotFound
}
