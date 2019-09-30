# tv-serial-control

Control TV Monitor via serial port

## Configuration

* `DRIVER` driver name
* `DEVICE` device
* `PORT` listen port

## Available Drivers

* `dummy` Default dummy driver
* `lg` LG monitors

> The dummy driver is available by default. To enable other supported driver
> pass the driver name as build tag to the build system.

## Endpoints

### `GET /driver`

Get driver info.

Response:

```json
{
  "driver": "lg",
  "device": "/dev/ttyS0",
  "commands": [
    "power_on",
    "power_off",
    "power_status",
    "volume_get",
    "volume_set",
    "key_home",
    "key_up",
    "key_left",
    "key_right",
    "key_down",
    "key_esc",
    "key_a",
    "key_z",
    "key_r"
  ]
}
```


### `POST /commands/:cmd`

Send command `cmd`


### Dummy Driver Command list
```    
    echo          
```

### LG Driver Command list and translation to  transmission protocol
```	
  power_on:        "ka 00 01"
  power_off:       "ka 00 00"
  power_status:    "ka 00 ff"
  volume_get:      "kf 00 ff"
  volume_set:      "kf 00 VAL"
  key_home:        "mc 00 7c"
  key_up:          "mc 00 40"
  ley_left:        "mc 00 07"
  key_enter:       "mc 00 7c"
  key_right:       "mc 00 06"
  key_down:        "mc 00 31"
  key_esc:         "mc 00 1b"
  key_a:           "mc 00 41"
  key_z:           "mc 00 5a"
  key_r:           "mc 00 52"
  volume_up:       "mc 00 02"
  volume_down:     "mc 00 03"
```
where  ` VAL`  is  a numeric value from 00 to 64 sent as a  string parameter"

#### Example with curl
```curl  -X POST -H "Content-Type: application/json" -d '{"args":["64"]}'   "http://localhost:8080/commands/volume_set"```
   

## License

Licensed under MIT, ([LICENSE](LICENSE) or https://opensource.org/licenses/MIT).

## Contribution

Unless you explicitly state otherwise, any contribution intentionally
submitted for inclusion in the work by you, as defined in the MIT
license, shall be licensed as above, without any additional terms or
conditions.
