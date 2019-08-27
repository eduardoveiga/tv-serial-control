/*
 * Copyright (C) 2019
 * O.S. Systems Sofware LTDA: contato@ossystems.com.br
 *
 * SPDX-License-Identifier: MIT
 */

package driverapi

var drivers = make(map[string]Driver)

type Driver interface {
	Initialize(name string) error
	AvailableCommands() []string
	SendCommand(name string, tv *TV, args ...interface{}) (map[string]interface{}, error)
}

func RegisterDriver(name string, driver Driver) {
	drivers[name] = driver
}

func GetDriver(name string) Driver {
	if d, ok := drivers[name]; ok {
		return d
	}

	return nil
}

type TV struct {
	//TV := make(map[string]int)
	//volume
	//power_on_off
	//power  int
	//volume int
	Data map[string]int64
}

func (tv *TV) GetData(cmd string) int64 {
	return tv.Data[cmd]
}

func (tv *TV) SetData(cmd string, data int64) {

	tv.Data[cmd] = data

}
