package main

import (
	"fmt"
	"os"
)

func raiseIfError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func raiseIfIdOOB(id int, length int) {
	if id >= length {
		fmt.Fprintf(os.Stderr, "id supplied is more than the notes collected. Need to be less than %d.\n", length)
		os.Exit(1)
	}
}
