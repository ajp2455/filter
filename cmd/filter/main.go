package main

import (
	"io"
	"log"
	"os"

	"github.com/ajp2455/filter"
)

func main() {
	f, err := filter.NewFilter(os.Stdin, []string{"test", "-x"})
	if err != nil {
		log.Println(err)
		return
	}

	if _, err := io.Copy(os.Stdout, f); err != nil {
		log.Println(err)
		return
	}
}
