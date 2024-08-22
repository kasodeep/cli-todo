package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/alexeyco/simpletable"
)

/*
Item represents a basic todo task that we can add to our todo list.
*/
type Item struct {
	Task        string
	Done        bool
	CreatedAt   time.Time
	CompletedAt time.Time
}

type Todos []Item

/*
Add adds a new todo item to our todo list.
*/
func (t *Todos) Add(task string) {
	todo := Item{
		Task:        task,
		Done:        false,
		CreatedAt:   time.Now(),
		CompletedAt: time.Time{},
	}

	*t = append(*t, todo)
}

/*
Complete marks a todo item as completed.
*/
func (t *Todos) Complete(index int) error {
	if index <= 0 || index > len(*t) {
		return errors.New("invalid index!💀")
	}

	(*t)[index-1].Done = true
	(*t)[index-1].CompletedAt = time.Now()
	return nil
}

/*
Delete removes a todo item from our todo list.
*/
func (t *Todos) Delete(index int) error {
	if index <= 0 || index > len(*t) {
		return errors.New("invalid index!💀")
	}

	*t = append((*t)[:index-1], (*t)[index:]...)
	return nil
}

/*
Load loads our todo list from a json file.
*/
func (t *Todos) Load(filename string) error {
	file, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	if len(file) == 0 {
		return err
	}

	// decoding the file contents into our todo list.
	err = json.Unmarshal(file, t)
	if err != nil {
		return err
	}

	return nil
}

/*
Store stores our todo list into a json file.
*/
func (t *Todos) Store(filename string) error {
	data, err := json.Marshal(t)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

/*
Print prints our todo list.
*/
func (t *Todos) Print() {
	table := simpletable.New()

	// adding the headers.
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: "#"},
			{Align: simpletable.AlignCenter, Text: "Task"},
			{Align: simpletable.AlignCenter, Text: "Done?"},
			{Align: simpletable.AlignRight, Text: "CreatedAt"},
			{Align: simpletable.AlignRight, Text: "CompletedAt"},
		},
	}

	// adding our todos to the table.
	var cells [][]*simpletable.Cell
	for idx, item := range *t {
		idx++
		task := blue(item.Task)
		done := red("no")

		if item.Done {
			task = green(fmt.Sprintf("\u2705 %s", item.Task))
			done = green("yes")
		}

		cells = append(cells, []*simpletable.Cell{
			{Text: fmt.Sprintf("%d", idx)},
			{Text: task},
			{Text: done},
			{Text: item.CreatedAt.Format(time.RFC822)},
			{Text: item.CompletedAt.Format(time.RFC822)},
		})
	}
	table.Body = &simpletable.Body{Cells: cells}

	table.Footer = &simpletable.Footer{Cells: []*simpletable.Cell{
		{Align: simpletable.AlignCenter, Span: 5, Text: blue(fmt.Sprintf("You have %d pending todos🙄", t.countPending()))},
	}}

	table.SetStyle(simpletable.StyleUnicode)
	table.Println()
}

/*
Helper function.
*/
func (t *Todos) countPending() int {
	total := 0
	for _, item := range *t {
		if !item.Done {
			total++
		}
	}

	return total
}
