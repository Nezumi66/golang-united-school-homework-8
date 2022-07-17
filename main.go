package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string
type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {
	if args["operation"] == "" {
		return errors.New("-operation flag has to be specified")
	}
	if args["fileName"] == "" {
		return errors.New("-fileName flag has to be specified")
	}

	switch args["operation"] {
	case "list":
		data, err := ioutil.ReadFile(args["fileName"])
		if len(data) == 0 {
			break
		}
		if err != nil {
			return err
		}

		_, err = writer.Write(data)
		if err != nil {
			return err
		}
	case "add":
		if args["item"] == "" {
			return fmt.Errorf("-item flag has to be specified")
		}

	case "remove":
		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")
		}
	case "findById":
		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")
		}
	default:
		return fmt.Errorf("Operation %v not allowed!", args["operation"])
	}
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
