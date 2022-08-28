package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/ikawaha/jisx0208"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) >= 1 {
		fmt.Println(jisx0208.ToValid(args[0], "□"))
		return nil
	}
	fp := os.Stdin
	s := bufio.NewScanner(fp)
	for s.Scan() {
		fmt.Println(jisx0208.ToValid(s.Text(), "□"))
	}
	return s.Err()
}
