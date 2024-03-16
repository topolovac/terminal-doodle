package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "tm"
	app.Usage = "A simple notes app for the terminal"

	hd, err := os.UserHomeDir()
	if err != nil {
		log.Panic(err)
	}
	fs := &NoteService{
		directory_path: hd + "/terminal-doodle",
	}

	app.Commands = []cli.Command{
		{
			Name:    "status",
			Aliases: []string{"s"},
			Usage:   "general status",
			Action: func(c *cli.Context) error {
				path := fs.getFilePath()
				fmt.Println("File path: ", path)
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a note",
			Action: func(c *cli.Context) error {
				note := c.Args().First()
				for i := 1; i < c.NArg(); i++ {
					note += " " + c.Args().Get(i)
				}
				err := fs.AddNote(note)
				if err != nil {
					return err
				}
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
		{
			Name:    "nano",
			Aliases: []string{"n"},
			Usage:   "open today's notes in nano",
			Action: func(c *cli.Context) error {
				file_path := fs.getFilePath()
				cmd := exec.Command("nano", file_path)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		{
			Name:    "vim",
			Aliases: []string{"v"},
			Usage:   "open today's notes in vim",
			Action: func(c *cli.Context) error {
				file_path := fs.getFilePath()
				cmd := exec.Command("vim", file_path)
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
				}
				return nil
			},
		},
		{
			Name:    "editor",
			Aliases: []string{"e"},
			Usage:   "open today's notes in default text editor",
			Action: func(c *cli.Context) error {
				file_path := fs.getFilePath()
				fmt.Println(file_path)
				var cmd *exec.Cmd
				switch runtime.GOOS {
				case "windows":
					cmd = exec.Command("start", file_path)
				case "darwin":
					cmd = exec.Command("open", "-t", file_path)
				case "linux":
					cmd = exec.Command("xdg-open", file_path)
				default:
					log.Fatalf("Unsupported OS: %s", runtime.GOOS)
				}
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil {
					log.Fatal(err)
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
	if _, err := os.Stat(n.directory_path); os.IsNotExist(err) {
		err := os.MkdirAll(n.directory_path, 0755)
		if err != nil {
			log.Panic(err)
		}
		fmt.Println("Directory created")
	}

	today := time.Now()
	file_name := today.Format("01-02-2006") + ".txt"
	return n.directory_path + "/" + file_name
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
