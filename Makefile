
CGO_ENABLED=0

GOTAGS=--tags="netgo,osusergo"

all: lib default onion garlic

lib:
	go build "${GOTAGS}" -o lib.a *.go

default:
	go build "${GOTAGS}" -o apt-transport-default/apt-transport-default apt-transport-default/main.go

onion:
	go build "${GOTAGS}" -o apt-transport-onion/apt-transport-onion apt-transport-onion/main.go

garlic:
	go build "${GOTAGS}" -o apt-transport-garlic/apt-transport-garlic apt-transport-garlic/main.go

install:
	install -m755 apt-transport-default/apt-transport-default /usr/lib/apt/methods/https
	install -m755 apt-transport-onion/apt-transport-onion /usr/lib/apt/methods/onion
	install -m755 apt-transport-garlic/apt-transport-garlic /usr/lib/apt/methods/garlic