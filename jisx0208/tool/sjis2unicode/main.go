package main

import (
	"fmt"
	"os"

	"github.com/ikawaha/encoding/internal"
)

func main() {
	mapper, err := NewJIS0213RuneMapper()
	if err != nil {
		fmt.Fprintf(os.Stderr, "sjis to unicode mapping construction failed: %v", err)
		os.Exit(1)
	}
	for k, v := range *mapper {
		if len(v) > 1 {
			fmt.Printf("%#0X(%d): %c\n", k, k, k)
			for _, vv := range v {
				fmt.Printf("%v\n", vv)
			}
		}
	}

	//<---
	//r := bufio.NewReader(transform.NewReader(bytes.NewReader(sjis.SJISCode{0x81, 0x80}), japanese.ShiftJIS.NewDecoder()))
	//c, _, err := r.ReadRune()
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("!!!%c\n", c)
	//<---

	runes := mapper.Runes()
	//fmt.Printf("runes: %d\n", len(runes))
	table := internal.RangeTable(runes)
	internal.DumpRangeTable(os.Stdout, table)
}
