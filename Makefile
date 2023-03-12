
all: lib default onion garlic

lib:
	go build -o lib.a *.go

default:
	go build -o apt-transport-default/apt-transport-default apt-transport-default/main.go

onion:
	go build -o apt-transport-onion/apt-transport-onion apt-transport-onion/main.go

garlic:
	go build -o apt-transport-garlic/apt-transport-garlic apt-transport-garlic/main.go