package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

const csvFilePath = "../../testdata/golden_jyouyou_H22-11-30.csv"

func main() {
	runes, err := loadRunesFromCSV(csvFilePath)
	if err != nil {
		log.Fatal(err)
	}
	table := RangeTable(runes)
	DumpRangeTable(os.Stdout, table)
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
		row := strings.Split(s.Text(), "\t")
		i++
		if len(row) < 1 {
			log.Println("empty line, line no:", i)
			continue
		}
		for _, v := range []rune(row[0]) {
			if v == '［' || v == '］' || v == '（' || v == '）' { // 餅［餅］（餠）
				continue
			}
			ret = append(ret, v)
		}
	}
	return ret, s.Err()
}
