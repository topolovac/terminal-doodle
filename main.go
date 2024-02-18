package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	fs := &NoteService{
		directory_path: "./notes/",
	}

	app.Commands = []cli.Command{
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "general status",
			Action: func(c *cli.Context) error {
				fmt.Println("Generate status")
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a note",
			Action: func(c *cli.Context) error {
				note := c.Args().First()
				for i := 0; i < c.NArg(); i++ {
					note += " " + c.Args().Get(i)
				}
				fs.AddNote(note)
				return nil
			},
		},
		{
			Name:    "today",
			Aliases: []string{"t"},
			Usage:   "get today's notes",
			Action: func(c *cli.Context) error {
				notes, err := fs.GetNotes()
				if err != nil {
					return err
				}

				if notes == "" {
					fmt.Println("No notes for today")
				} else {
					fmt.Println(notes)
				}
				return nil
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type NoteService struct {
	directory_path string
}

func (n *NoteService) GetActiveFile() (*os.File, error) {
	file_path := n.getFilePath()
	file, err := os.OpenFile(file_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (n *NoteService) getFilePath() string {
	today := time.Now()
	file_name := today.Format("01-02-2006") + ".txt"
	return n.directory_path + file_name
}

func (n *NoteService) AddNote(note string) error {
	file, err := n.GetActiveFile()
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(note + "\n")
	file.Close()
	if err != nil {
		return err
	}
	return nil
}

func (n *NoteService) GetNotes() (string, error) {
	file_path := n.getFilePath()

	content, err := os.ReadFile(file_path)

	if err != nil {
		return "", err
	}

	return string(content), nil
}
