/*
 * Copyright (C) 2019
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: MIT
 */

package dummy

import (
	"github.com/gustavosbarreto/tv-control/driverapi"
)

func init() {
	driverapi.RegisterDriver("dummy", &driver{})
}

type driver struct {
}

func (d *driver) Initialize(device string) error {
	return nil
}

func (l *driver) AvailableCommands() []string {
	return []string{}
}

func (l *driver) SendCommand(name string, args ...interface{}) (map[string]interface{}, error) {
	return nil, nil
}
