package main

import "log"

func main() {
	// khoi tao struct Demo

	d, err := initService()
	if err != nil {
		panic(err)
	}
	log.Print("hello")
	GRPCServe("9090", d)
}
