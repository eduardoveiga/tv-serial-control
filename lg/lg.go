package lg

import (
	"time"

	"github.com/tarm/serial"
)

type LG struct {
	port *serial.Port
}

func (l *LG) Initialize(device string) error {
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

func (l *LG) AvailableCommands() []string {
	return []string{}
}

func (l *LG) SendCommand(name string, args ...interface{}) (map[string]interface{}, error) {
	if _, ok := cmds[Command(name)]; ok {
		res, err := Command(name).Send(l.port, args...)
		if err != nil {
			return nil, err
		}

		return res, nil
	}

	return nil, UnknownCommandErr
}
