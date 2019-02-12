# Configuration

* `DRIVER` driver name
* `DEVICE` device
* `PORT` listen port

# Available Drivers

* `dummy` Default dummy driver
* `lg` LG monitors

> The dummy driver is available by default. To enable other supported driver
> pass the driver name as build tag to the build system.

# Endpoints

## `GET /driver`

Get driver info.

Response:

```json
{
  "driver": "lg",
  "device": "/dev/ttyS0",
  "commands": [
    "power_on",
    "power_off"
  ]
}
```

## `POST /commands/:cmd`

Send command `cmd`
