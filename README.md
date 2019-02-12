# Configuration

* `DRIVER` driver name
* `DEVICE` device
* `PORT` listen port

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