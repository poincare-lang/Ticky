package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Tasks struct {
	Tasks []string `json:"tasks"`
}

func createFiles() {
	homeDir, err := os.UserHomeDir()
	err = os.MkdirAll(filepath.Join(homeDir, ".tick"), os.ModePerm)

	if err != nil {
		fmt.Println("Alas, theres been an error creating the tick directory")
	}
}

func write(t []string) {
	homeDir, err := os.UserHomeDir()
	tasks := Tasks{t}
	j, err := json.Marshal(tasks)
	err = os.WriteFile(filepath.Join(homeDir, ".tick/tasks.json"), j, os.ModePerm)
	if err != nil {
		fmt.Println("Alas, there's been a json writing error", err)
	}
}

func read() []string {

	homeDir, err := os.UserHomeDir()

	file, err := os.Open(filepath.Join(homeDir, ".tick/tasks.json"))
	if os.IsNotExist(err) {
		return nil
	}

	defer file.Close()
	byteValue, _ := io.ReadAll(file)
	var tasks Tasks
	err = json.Unmarshal(byteValue, &tasks)

	if err != nil {
		fmt.Println("Alas, theres been an error reading the json data")
	}

	return tasks.Tasks

}
