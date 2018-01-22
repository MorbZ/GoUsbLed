## GoUsbLed
[![Go Report Card](https://goreportcard.com/badge/github.com/MorbZ/GoUsbLed)](https://goreportcard.com/report/github.com/MorbZ/GoUsbLed) [![Docker Build Status](https://img.shields.io/docker/build/morbz/gousbled.svg)](https://hub.docker.com/r/morbz/gousbled/)

Bitcoin price ticker for the USB Message Board from Dream Cheeky. Shows the current Bitcoin USD price from Bitstamp.

![USB Message Board](https://raw.githubusercontent.com/MorbZ/GoUsbLed/master/img/board.jpg)

## Installation
### Using Go
    $ go get github.com/MorbZ/GoUsbLed
    $ [sudo] GoUsbLed

### Using Docker
    $ docker run --rm --privileged -v /dev/bus/usb:/dev/bus/usb morbz/gousbled

Or using Docker compose:

    $ git clone https://github.com/MorbZ/GoUsbLed
    $ cd GoUsbLed
    $ docker-compose up -d

### Using binaries
Download the binary for your OS and architecture from the [releases page](https://github.com/MorbZ/GoUsbLed/releases).  
Unpack and execute it. On Linux the program needs to be run as superuser.
