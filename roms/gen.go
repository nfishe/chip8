//go:build ignore
// +build ignore

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
)

var outFilename = flag.String("output", "roms.go", "output file name")

func main() {
	flag.Parse()

	var buf bytes.Buffer

	fmt.Fprintln(&buf, "// Code generated by go run gen.go -output palette.go; DO NOT EDIT.")
	fmt.Fprintln(&buf)
	fmt.Fprintln(&buf, "package roms")
	fmt.Fprintln(&buf)

	b, err := ioutil.ReadFile("../testdata/test_blinky")
	if err != nil {
		panic(err)
	}

	fmt.Fprintf(&buf, "var Blinky = %#v", b)

	data, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(*outFilename, data, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
