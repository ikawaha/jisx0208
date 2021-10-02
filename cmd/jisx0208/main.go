package main

import (
	"fmt"
	"os"

	"github.com/ikawaha/jisx0208"
)

func main() {
	if err := run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: jisx0208 <string>")
	}
	fmt.Println(jisx0208.ToValid(args[1], "□"))
	return nil
}
