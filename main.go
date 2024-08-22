package main

import (
	"bufio"
	"errors"
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

const (
	todoFile = "C:\\Users\\admin\\go\\bin\\.todos.json"
)

func main() {
	// flags to perform specific actions.
	add := flag.Bool("add", false, "Add a new todo!")
	complete := flag.Int("complete", 0, "Mark a todo as completed!")
	del := flag.Int("del", 0, "Delete an old todo!")
	list := flag.Bool("list", false, "Print list all todos!")

	flag.Parse()
	todos := Todos{}

	// loading the todos.
	if err := todos.Load(todoFile); err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	switch {
	case *add:
		task, err := getInput(os.Stdin, flag.Args()...)
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}
		todos.Add(task)
		Store(todos)

	case *complete > 0:
		err := todos.Complete(*complete)
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}
		Store(todos)

	case *del > 0:
		err := todos.Delete(*del)
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}
		Store(todos)

	case *list:
		todos.Print()

	default:
		log.Println("Invalid Command.😒")
		os.Exit(0)
	}
}

/*
Store stores the todos after the update.
*/
func Store(todos Todos) {
	err := todos.Store(todoFile)

	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}
}

/*
Input get user input for the add flag.
*/
func getInput(r io.Reader, args ...string) (string, error) {
	if len(args) > 0 {
		return strings.Join(args, " "), nil
	}

	scanner := bufio.NewScanner(r)
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	text := scanner.Text()
	if len(text) == 0 {
		return "", errors.New("empty todo is not allowed!🥱")
	}
	return text, nil
}
