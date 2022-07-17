package main

import (
	"flag"
	"io"
	"os"
)

type Arguments map[string]string

func Perform(args Arguments, writer io.Writer) error {
	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}

func parseArgs() Arguments {
	id := flag.String("id", "", "for identifying")
	item := flag.String("item", "", "for transferring")
	op := flag.String("operation", "", "for operating")
	file := flag.String("fileName", "", "for performing with the file")

	flag.Parse()

	return Arguments{
		"id":        *id,
		"item":      *item,
		"operation": *op,
		"fileName":  *file,
	}
}
