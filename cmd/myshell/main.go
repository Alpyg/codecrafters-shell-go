package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
prompt:
	for {
		fmt.Fprint(os.Stdout, "$ ")

		cmd, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			os.Exit(1)
		}
		cmd = strings.TrimSpace(cmd)

		cmds := strings.Split(cmd, " ")
		switch cmds[0] {
		case "exit":
			code, err := strconv.Atoi(cmds[1])
			if err != nil {
				os.Exit(1)
			}
			os.Exit(code)
		case "echo":
			fmt.Fprintf(os.Stdout, "%s\n", strings.Join(cmds[1:], " "))
		case "type":
			switch cmds[1] {
			case "exit", "echo", "type", "pwd", "cd":
				fmt.Printf("%s is a shell builtin\n", cmds[1])
			default:
				for _, path := range strings.Split(os.Getenv("PATH"), ":") {
					fp := filepath.Join(path, cmds[1])
					if _, err := os.Stat(fp); err == nil {
						fmt.Printf("%s is %s\n", cmds[1], fp)
						continue prompt
					}
				}
				fmt.Printf("%s: not found\n", cmds[1])
			}
		case "pwd":
			wd, _ := os.Getwd()
			fmt.Printf("%s\n", wd)
		case "cd":
			path := cmds[1]
			if cmds[1] == "~" {
				path = os.Getenv("HOME")
			}
			if err := os.Chdir(path); err != nil {
				fmt.Printf("cd: %s: No such file or directory\n", path)
			}
		default:
			cmd := exec.Command(cmds[0], cmds[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr

			err := cmd.Run()
			if err != nil {
				fmt.Printf("%s: command not found\n", cmd)
			}
		}
	}
}
