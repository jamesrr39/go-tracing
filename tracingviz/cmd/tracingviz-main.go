package main

import (
	"log"

	"github.com/jamesrr39/go-tracing/tracingviz"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	dataFilePath := kingpin.Arg("data file path", "file path to the profile data").Required().String()
	outFilePath := kingpin.Arg("out file path", "file path to the out file").Required().String()
	kingpin.Parse()

	err := tracingviz.Generate(*dataFilePath, *outFilePath)
	if err != nil {
		log.Fatalf("error: %s\nStack trace:\n%s\n", err.Error(), err.Stack())
	}
}
