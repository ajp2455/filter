package main

import (
	"io"
	"log"
	"os"
	"runtime"

	"github.com/ajp2455/filter"
	"github.com/pborman/getopt/v2"
)

func main() {
	token := getopt.StringLong("replace", 'r', "", "token to replace in command")
	threads := getopt.IntLong("threads", 't', runtime.NumCPU(), "number of threads")
	help := getopt.BoolLong("help", 'h', "display help")

	getopt.Parse()

	if *help {
		getopt.PrintUsage(os.Stderr)
		return
	}

	command := getopt.Args()

	if len(command) == 0 {
		getopt.PrintUsage(os.Stderr)
		return
	}

	f, err := filter.Filter(os.Stdin, *threads, filter.CmdPredicate(command, *token))
	if err != nil {
		log.Println(err)
		return
	}

	if _, err := io.Copy(os.Stdout, f); err != nil {
		log.Println(err)
		return
	}
}
