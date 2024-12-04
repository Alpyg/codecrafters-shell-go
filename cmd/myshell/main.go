package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

type Builtiner interface {
	Builtin(args ...string) int
}

var builtins = map[string]Builtiner{
	"exit": ExitBuiltin{},
	"echo": EchoBuiltin{},
	"type": TypeBuiltin{},
	"pwd":  PrintWorkingDirectoryBuiltin{},
	"cd":   ChangeDirectoryBuiltin{},
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")

		input, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			os.Exit(1)
		}
		input = strings.TrimSpace(input)

		cmds := processInput(input)
		var builtin Builtiner
		builtin, ok := builtins[cmds[0]]

		if ok {
			builtin.Builtin(cmds[1:]...)
			continue
		}

		cmd := exec.Command(cmds[0], cmds[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err = cmd.Run()
		if err != nil {
			fmt.Printf("%s: command not found\n", cmd)
		}
	}
}

func processInput(input string) []string {
	re := regexp.MustCompile(`"([^"]*)"|'([^']*)'|(\S+)`)
	matches := re.FindAllStringSubmatch(input, -1)

	var cmds []string
	for _, match := range matches {
		if match[1] != "" {
			cmds = append(cmds, match[1])
		} else if match[2] != "" {
			cmds = append(cmds, match[2])
		} else {
			word := match[3]
			cleanWord := regexp.MustCompile(`\\(.)`).ReplaceAllString(word, `$1`)
			cmds = append(cmds, cleanWord)
		}
	}

	return cmds
}

type ExitBuiltin struct{}

func (b ExitBuiltin) Builtin(args ...string) int {
	if len(args) == 0 {
		os.Exit(0)
	}
	code, err := strconv.Atoi(args[0])
	if err != nil {
		os.Exit(1)
	}
	os.Exit(code)
	return 0
}

type EchoBuiltin struct{}

func (b EchoBuiltin) Builtin(args ...string) int {
	fmt.Fprintf(os.Stdout, "%s\n", strings.Join(args, " "))

	return 0
}

type TypeBuiltin struct{}

func (b TypeBuiltin) Builtin(args ...string) int {
	_, ok := builtins[args[0]]
	if ok {
		fmt.Printf("%s is a shell builtin\n", args[0])
		return 0
	}

	for _, path := range strings.Split(os.Getenv("PATH"), ":") {
		fp := filepath.Join(path, args[0])
		if _, err := os.Stat(fp); err == nil {
			fmt.Printf("%s is %s\n", args[0], fp)
			return 0
		}
	}
	fmt.Printf("%s: not found\n", args[0])
	return 1
}

type PrintWorkingDirectoryBuiltin struct{}

func (b PrintWorkingDirectoryBuiltin) Builtin(args ...string) int {
	wd, err := os.Getwd()
	if err != nil {
		return 1
	}
	fmt.Printf("%s\n", wd)

	return 0
}

type ChangeDirectoryBuiltin struct{}

func (b ChangeDirectoryBuiltin) Builtin(args ...string) int {
	path := args[0]
	if args[0] == "~" {
		path = os.Getenv("HOME")
	}
	if err := os.Chdir(path); err != nil {
		fmt.Printf("cd: %s: No such file or directory\n", path)
		return 1
	}

	return 0
}
