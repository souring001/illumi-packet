# ILLUMI PACKET

Illuminating Packets on an Ethernet Cable using LED Strip.

## Overview

![illumi-packet](https://user-images.githubusercontent.com/29009733/70907987-8ab66000-204d-11ea-86e2-09a34d7c557a.jpg)

[Movie (YouTube)](https://youtu.be/5yJyNpHeRzg)

A packet is visualized by LED lights as it flows in the transmission direction.

ARP packets glow ORANGE and DHCP packets glow BLUE for example.
8 LED colors were assigned to each packet type.

Therefore, Illumi Packet makes the presence of packets familiar and helps to intuitively understand what kind of packets are generated while operating a computer.

![color](https://user-images.githubusercontent.com/29009733/71455676-786cbc80-27d9-11ea-980c-99a22d31696f.png)

## Hardware Setup

### Requirements

|Materials|Quantity|Remarks|
|:-|-:|:-|
|LED strip (WS281B)|1m|Recommended > Allowing Individual address, full color. 144LEDs/m|
|Raspberry Pi|1|Operation confirmed > Raspberry Pi 2, 3, 4|
|Ethernet cable|1m||
|Jumper wire male-female|3||
|Cable tie|4||

You need also Keyboard, Display, HDMI Cable, Routers etc.

### Assembly

1. Connect an LED strip to the GPIO of a Raspberry Pi.
2. Fix it to an Ethernet cable.
3. Connect one end of the Ethernet cable with the Raspberry Pi and the other end with several wired network access points.

![GPIO_Outline](https://user-images.githubusercontent.com/29009733/71317350-aba20980-24c2-11ea-8a59-47388f5b2d73.png)

![GPIO](https://user-images.githubusercontent.com/29009733/70908199-f7315f00-204d-11ea-9cb0-256967c7ca5e.png)


### Setup on Raspberry Pi

1. Install [golang](https://golang.org/doc/install#install)
2. Install libpcap `sudo apt-get install libpcap-dev`
3. Install SCons `sudo apt-get install scons`
4. Install [rpi_ws281x](https://github.com/jgarff/rpi_ws281x)
5. Run `git clone https://github.com/souring001/illumi-packet.git`
6. Change the parameters in `illumi-packet.go` according to the number of LEDs as follows:

| LEDs/m | count | speed | series |
| ------:| -----:| -----:| ------:|
|60      |    60 |      1 |     6 |
|144     |   144 |      4 |    12 |
|144 (0.5m)    |   72 |      1 |    12 |


#### 1. Install golang

The program to be executed is described in go language.
Install golang according to https://golang.org/doc/install#install.
Execute it in the terminal as shown below, and confirm that the version is displayed at the end.
```sh
$ version=1.13.4
$ wget https://storage.googleapis.com/golang/go${version}.linux-armv6l.tar.gz
$ sudo tar -C /usr/local -xzf go${version}.linux-armv6l.tar.gz

$ echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.profile
$ . ~/.profile

$ go version
go version go1.13.4 linux/arm
```

#### 4. Install rpi_ws281x
An LED strip is controlled by a library called rpi_ws281x.
Install rpi_ws281x according to https://github.com/jgarff/rpi_ws281x.

```sh
$ git clone https://github.com/jgarff/rpi_ws281x.git
$ cd rpi_ws281x
$ scons

$ sudo cp -ai ./ws2811.h ./rpihw.h ./pwm.h /usr/local/include/
$ sudo cp -ai ./libws2811.a /usr/local/lib/
```

## Build

```sh
$ go build illumi-packet.go
```

## Run

```sh
$ sudo ./illumi-packet
```

Press Ctr-C to quit.

|Option||
|:-|:-|
|-h|Help command|
|-debug |Print packet details. (default: `true`)|
|-device [string]|Set network interface (default: `eth0`)|
|-speed [int]|Set speed of flowing packet(default: `1`)|
|-narp|Disable visualizing ARP packets|
|-ntcp|Disable visualizing TCP packets|
|-nudp|Disable visualizing UDP packets|
|-reset|Reset LEDs|
|-ipaddr|Show IP address on LED|

### Examples

Not showing TCP/UDP packets.
```sh
$ sudo ./illumi-packet -nudp -ntcp
```

<br>

Not showing packet details. (recommend on SSH)
```sh
$ sudo ./illumi-packet -debug=false
```

<br>

Visualize packets on Wi-Fi.
```sh
$ sudo ./illumi-packet -device wlan0
```

<br>

Show IP address on LED.
```sh
$ sudo ./illumi-packet -ipaddr
```
![showipaddress](https://user-images.githubusercontent.com/29009733/70908359-5e4f1380-204e-11ea-9187-a2d385c9f300.JPG)

Turn off the LED lights.
```sh
$ sudo ./illumi-packet -reset
```

## License

<a rel="license" href="http://creativecommons.org/licenses/by/4.0/"><img alt="Creative Commons License" style="border-width:0" src="https://i.creativecommons.org/l/by/4.0/88x31.png" /></a><br /><span xmlns:dct="http://purl.org/dc/terms/" property="dct:title">ILLUMI PACKET</span> by <span xmlns:cc="http://creativecommons.org/ns#" property="cc:attributionName">Kohei Aso</span> is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by/4.0/">Creative Commons Attribution 4.0 International License</a>.

- Free to modify and redistribute
- Displaying the credits is required

We may introduce your usage as an examples on our site.  
Please contact us if you do not want to give credit.

## Contact

Twitter: [@souring001](https://twitter.com/souring001)


ILLUMI PACKET
<br />Copyright (c) 2019, Kohei Aso
