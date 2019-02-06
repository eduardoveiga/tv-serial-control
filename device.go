package main

type Device interface {
	Initialize() error
	AvailableCommands() []string
	SendCommand(name string, args ...interface{}) (map[string]interface{}, error)
}
