package main

import (
	"encoding/json"
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

	var users []User
	file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if len(data) != 0 {
		err = json.Unmarshal(data, &users)
		if err != nil {
			return err
		}
	}

	switch args["operation"] {
	case "list":
		_, err = writer.Write(data)
		if err != nil {
			return err
		}
	case "add":
		if args["item"] == "" {
			return fmt.Errorf("-item flag has to be specified")
		}

		if !json.Valid([]byte(args["item"])) {
			return fmt.Errorf("incorrect format of an item")
		}

		var user User
		err := json.Unmarshal([]byte(args["item"]), &user)
		if err != nil {
			return err
		}

		users = append(users, user)

		_, err = file.Seek(0, 0)
		if err != nil {
			return err
		}
		err = file.Truncate(0)
		if err != nil {
			return err
		}
		usersToWrite, err := json.Marshal(users)
		if err != nil {
			return err
		}
		_, err = file.Write(usersToWrite)
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
