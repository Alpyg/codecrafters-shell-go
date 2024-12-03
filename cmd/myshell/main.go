package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
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
			case "exit", "echo", "type":
				fmt.Printf("%s is a shell builtin\n", cmds[1])
			default:
				fmt.Printf("%s: not found\n", cmds[1])
			}
		default:
			fmt.Printf("%s: command not found\n", cmd)
		}
	}
}
