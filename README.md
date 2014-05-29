# sphere-leds

Listens for MQTT messages and toggles the leds on the sphere to various colors.

# Usage

# Messages

Set the PWM brightness to 100.

```
$ mosquitto_pub -m '{"brightness": 100}' -t '$hardware/status/reset'
```

Set the PWM brightness to 0.

```
$ mosquitto_pub -m '{"brightness": 0}' -t '$hardware/status/reset'
```

Change the "power" led to "blue".

```
$ mosquitto_pub -m '{"color": "blue", "flash": true}' -t '$hardware/status/power'
```

Valid led and color names are listed below:

```go

var Colors = map[string][]int{
  "black":   {0, 0, 0},
  "red":     {1, 0, 0},
  "green":   {0, 1, 0},
  "blue":    {0, 0, 1},
  "cyan":    {0, 1, 1},
  "magenta": {1, 0, 1},
  "yellow":  {1, 1, 0},
  "white":   {1, 1, 1},
}

var LedNames = []string{
  "power",
  "wired_internet",
  "wireless",
  "pairing",
  "radio",
}

````


