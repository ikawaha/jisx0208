package main

import (
	"bufio"
	"log"
	"os"

	"github.com/ikawaha/encoding/unicode"
)

const csvFilePath = "../../testdata/golden_jyouyou_H22-11-30.csv"

func main() {
	runes, err := loadRunesFromCSV(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}
	table := unicode.RangeTable(runes)
	unicode.DumpRangeTable(os.Stdout, table)
}

func loadRunesFromCSV(path string) ([]rune, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	s := bufio.NewScanner(f)

	var ret []rune
	i := 0
	for s.Scan() {
		row := []rune(s.Text())
		i++
		if len(row) < 1 {
			log.Println("empty line, line no:", i)
			continue
		}
		ret = append(ret, row[0])
	}
	return ret, s.Err()
}
