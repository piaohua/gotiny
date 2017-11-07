# gotiny

This sample is a tiny game server

## Installation

```
cd bin
GOPATH=$(cd ../"$(dirname "$0")"; pwd)
go build -o server -ldflags "-w -s" ../src/gotiny/server.go
go build -o robot -ldflags "-w -s" ../src/gotiny/robot.go
```

## Usage:

```
./server -log_dir="log" > /dev/null 2>&1 &
./robot -log_dir="log" > /dev/null 2>&1 &
```

## TODO
