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
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(file)

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

		var user User
		err := json.Unmarshal([]byte(args["item"]), &user)
		if err != nil {
			return err
		}

		foundId := false

		for _, v := range users {
			if v.Id == user.Id {
				foundId = true
			}
		}
		if foundId {
			_, err = writer.Write([]byte(fmt.Sprintf("Item with id %s already exists", user.Id)))
			if err != nil {
				return err
			}
		} else {
			users = append(users, user)
			err = WriteToFile(users, *file)
			if err != nil {
				return err
			}
		}
	case "remove":
		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")
		}

		foundId := false

		var newUsers []User
		for _, v := range users {
			if args["id"] == v.Id {
				foundId = true
			} else {
				newUsers = append(newUsers, v)
			}
		}

		if foundId {
			err = WriteToFile(newUsers, *file)
			if err != nil {
				return err
			}
		} else {
			_, err = writer.Write([]byte(fmt.Sprintf("Item with id %s not found", args["id"])))
			if err != nil {
				return err
			}
		}

	case "findById":
		if args["id"] == "" {
			return fmt.Errorf("-id flag has to be specified")
		}

		foundId := false
		var foundUser User

		for _, v := range users {
			if args["id"] == v.Id {
				foundId = true
				foundUser = v
			}
		}

		if foundId {
			userToWrite, err := json.Marshal(foundUser)
			if err != nil {
				return err
			}

			_, err = writer.Write(userToWrite)
			if err != nil {
				return err
			}
		} else {
			_, err = writer.Write([]byte(""))
			if err != nil {
				return err
			}
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

func WriteToFile(users []User, file os.File) error {
	_, err := file.Seek(0, 0)
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
	if err != nil {
		return err
	}

	return nil
}
