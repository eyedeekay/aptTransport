
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