package driverapi

var drivers = make(map[string]Driver)

type Driver interface {
	Initialize(name string) error
	AvailableCommands() []string
	SendCommand(name string, args ...interface{}) (map[string]interface{}, error)
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
