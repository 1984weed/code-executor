package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/urfave/cli/v2"
)

type responsePayload struct {
	Result             string `json:"result"`
	ExecTime           int    `json:"execTime,omitempty"`
	LastInput          string `json:"lastInput,omitempty"`
	LastOutput         string `json:"lastOutput,omitempty"`
	LastExpectedOutput string `json:"lastExpectedOutput,omitempty"`
}
type payload struct {
	Language      string     `json:"language"`
	Content       string     `json:"content"`
	InputCount    int        `json:"inputCount"`
	Inputs        [][]string `json:"inputs"`
	ArgumentCount int        `json:"argumentCount"`
	ExpectOutputs []string   `json:"expectOutputs"`
}

func main() {
	app := &cli.App{
		Name:  "code executor",
		Usage: "make an explosive entrance",
	}

	app.Action = func(c *cli.Context) error {
		reader := bufio.NewReader(os.Stdin)
		payloadBytes, _ := reader.ReadBytes('\n')

		var payload payload

		if err := json.Unmarshal(payloadBytes, &payload); err != nil {
			return err
		}

		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		id := uuid.New()

		workDir := fmt.Sprintf("%s/workspace/%s-sandboxfiles", path, id)
		sand := SandBoxRunner{
			Dir:           workDir,
			FileName:      "main.py",
			Inputs:        payload.Inputs,
			ExpectOutputs: payload.ExpectOutputs,
			Content:       payload.Content,
			Timeout:       1 * time.Second,
		}

		sand.prepare()

		response, err := sand.run()

		if err != nil {
			return err
		}

		res, err := json.Marshal(response)
		if err != nil {
			return err
		}
		fmt.Fprint(os.Stdout, string(res))

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

type SandBoxRunner struct {
	Dir           string
	FileName      string
	path          string
	Content       string
	Inputs        [][]string
	ExpectOutputs []string
	Command       string
	Timeout       time.Duration
}

func (s *SandBoxRunner) prepare() error {
	err := os.MkdirAll(s.Dir, os.ModePerm)

	if err != nil {
		return err
	}

	s.path = fmt.Sprintf("%s/%s", s.Dir, s.FileName)

	file := NewFileUtils(s.path)
	file.WriteCreateFile(s.Content)

	return nil
}

func (s *SandBoxRunner) run() (*responsePayload, error) {
	defer os.RemoveAll(s.Dir)

	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout)
	defer cancel()

	// result states are "success", "fail", "tle", "error"
	result := "error"

	for i, inputs := range s.Inputs {
		cmd := exec.CommandContext(ctx, "python3", s.path)

		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		input := ""
		for _, argInput := range inputs {
			input += fmt.Sprintf("%s\n", argInput)
		}

		cmd.Stdin = strings.NewReader(input + s.ExpectOutputs[i] + "\n")

		if err := cmd.Run(); err != nil {
			return nil, err
		}

		r := csv.NewReader(strings.NewReader(out.String()))

		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}

			if record[0] == "success" {
				result = "success"
			}
			if record[0] == "fail" {
				result = "fail"
			}
		}

		if result != "success" {
			break
		}
	}

	return &responsePayload{
		Result: result,
	}, nil
}

type FileUtil struct {
	path string
}

func NewFileUtils(path string) *FileUtil {
	return &FileUtil{
		path: path,
	}
}

func (f *FileUtil) Create() error {
	var _, err = os.Stat(f.path)

	if os.IsNotExist(err) {
		var file, err = os.Create(f.path)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	return nil
}

func (f *FileUtil) Write(content string) error {
	var file, err = os.OpenFile(f.path, os.O_RDWR, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	err = file.Sync()
	if err != nil {
		return err
	}

	return nil
}

func (f *FileUtil) WriteCreateFile(content string) error {
	err := f.Create()
	if err != nil {
		return err
	}
	err = f.Write(content)
	if err != nil {
		return err
	}

	return nil
}
