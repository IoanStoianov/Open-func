package main

import (
	"encoding/json"
	"io"

	adapter "github.com/IoanStoianov/Open-func/pkg/adapters/golang"
)

type expectedInput struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func readInput(input io.ReadCloser) (expectedInput, error) {
	var payload expectedInput

	decoder := json.NewDecoder(input)
	if err := decoder.Decode(&payload); err != nil {
		return payload, err
	}
	defer input.Close()
	return payload, nil
}

func exampleFunc(input io.ReadCloser) string {
	payload, err := readInput(input)
	if err != nil {
		return "Invelid Input"
	}
	return payload.Name
}

func main() {
	adapter.TriggerListener(exampleFunc)
}
