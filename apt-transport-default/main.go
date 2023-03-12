package main

import "github.com/eyedeekay/apttransport"

func main() {
	transport := &apttransport.AptMethod{}
	transport.Main = transport.DefaultMain
	transport.Main()
}
