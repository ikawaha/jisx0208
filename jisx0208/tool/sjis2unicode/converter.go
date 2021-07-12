package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"unicode"

	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

const MaxUint16 = 1<<16 - 1

type SJISCode []byte

func (c SJISCode) String() string {
	if len(c) == 0 {
		return "<nil>"
	}
	ret := "0x"
	for i := range c {
		ret += fmt.Sprintf("%X", c[i])
	}
	return ret
}

type CodeRange struct {
	Start uint
	End   uint
}

// https://www.asahi-net.or.jp/~ax2s-kmtn/ref/jisx0213/index.html
var JISX0208SJIS = []CodeRange{
	// 1区
	{
		Start: 0x813F + 0x1,
		End:   0x8190 + 0xE,
	},
	// 2区
	{
		Start: 0x819E + 0x1,
		End:   0x819E + 0xE,
	},
	{
		Start: 0x81AE + 0xA,
		End:   0x81BE + 0x1,
	},
	{
		Start: 0x81BE + 0xA,
		End:   0x81CE + 0x0,
	},
	{
		Start: 0x81CE + 0xC,
		End:   0x81DE + 0xA,
	},
	{
		Start: 0x81EE + 0x2,
		End:   0x81EE + 0x9,
	},
	{
		Start: 0x81EE + 0xE,
		End:   0x81EE + 0xE,
	},
	// 3区
	{
		Start: 0x824F + 0x0,
		End:   0x824F + 0x9,
	},
	{
		Start: 0x825F + 0x1,
		End:   0x826F + 0xA,
	},
	{
		Start: 0x8280 + 0x1,
		End:   0x8290 + 0xA,
	},
	// 4区
	{
		Start: 0x829E + 0x1,
		End:   0x82EE + 0x3,
	},
	// 5区
	{
		Start: 0x833F + 0x1,
		End:   0x8390 + 0x6,
	},
	// 6区
	{
		Start: 0x839E + 0x1,
		End:   0x83AE + 0x8,
	},
	{
		Start: 0x83BE + 0x1,
		End:   0x83CE + 0x8,
	},
	// 7区
	{
		Start: 0x843F + 0x1,
		End:   0x845F + 0x1,
	},
	{
		Start: 0x846F + 0x1,
		End:   0x8490 + 0x1,
	},
	// 8区
	{
		Start: 0x849E + 0x1,
		End:   0x84BE + 0x0,
	},
	// 16区
	{
		Start: 0x889E + 0x1,
		End:   0x88EE + 0xE,
	},
	// 17区
	{
		Start: 0x893F + 0x1,
		End:   0x8990 + 0xE,
	},
	// 18区
	{
		Start: 0x899E + 0x1,
		End:   0x89EE + 0xE,
	},
	// 19区
	{
		Start: 0x8A3F + 0x1,
		End:   0x8A90 + 0xE,
	},
	// 20区
	{
		Start: 0x8A9E + 0x1,
		End:   0x8AEE + 0xE,
	},
	// 21区
	{
		Start: 0x8B3F + 0x1,
		End:   0x8B90 + 0xE,
	},
	// 22区
	{
		Start: 0x8B9E + 0x1,
		End:   0x8BEE + 0xE,
	},
	// 23区
	{
		Start: 0x8C3F + 0x1,
		End:   0x8C90 + 0xE,
	},
	// 24区
	{
		Start: 0x8C9E + 0x1,
		End:   0x8CEE + 0xE,
	},
	// 25区
	{
		Start: 0x8D3F + 0x1,
		End:   0x8D90 + 0xE,
	},
	// 26区
	{
		Start: 0x8D9E + 0x1,
		End:   0x8DEE + 0xE,
	},
	// 27区
	{
		Start: 0x8E3F + 0x1,
		End:   0x8E90 + 0xE,
	},
	// 28区
	{
		Start: 0x8E9E + 0x1,
		End:   0x8EEE + 0xE,
	},
	// 29区
	{
		Start: 0x8F3F + 0x1,
		End:   0x8F90 + 0xE,
	},
	// 30区
	{
		Start: 0x8F9E + 0x1,
		End:   0x8FEE + 0xE,
	},
	// 31区
	{
		Start: 0x903F + 0x1,
		End:   0x9090 + 0xE,
	},
	// 32区
	{
		Start: 0x909E + 0x1,
		End:   0x90EE + 0xE,
	},
	// 33区
	{
		Start: 0x913F + 0x1,
		End:   0x9190 + 0xE,
	},
	// 34区
	{
		Start: 0x919E + 0x1,
		End:   0x91EE + 0xE,
	},
	// 35区
	{
		Start: 0x923F + 0x1,
		End:   0x9290 + 0xE,
	},
	// 36区
	{
		Start: 0x929E + 0x1,
		End:   0x92EE + 0xE,
	},
	// 37区
	{
		Start: 0x933F + 0x1,
		End:   0x9390 + 0xE,
	},
	// 38区
	{
		Start: 0x939E + 0x1,
		End:   0x93EE + 0xE,
	},
	// 39区
	{
		Start: 0x943F + 0x1,
		End:   0x9490 + 0xE,
	},
	// 40区
	{
		Start: 0x949E + 0x1,
		End:   0x94EE + 0xE,
	},
	// 41区
	{
		Start: 0x953F + 0x1,
		End:   0x9590 + 0xE,
	},
	// 42区
	{
		Start: 0x959E + 0x1,
		End:   0x95EE + 0xE,
	},
	// 43区
	{
		Start: 0x963F + 0x1,
		End:   0x9690 + 0xE,
	},
	// 44区
	{
		Start: 0x969E + 0x1,
		End:   0x96EE + 0xE,
	},
	// 45区
	{
		Start: 0x973F + 0x1,
		End:   0x9790 + 0xE,
	},
	// 46区
	{
		Start: 0x979E + 0x1,
		End:   0x97EE + 0xE,
	},
	// 47区
	{
		Start: 0x983F + 0x1,
		End:   0x986F + 0x3,
	},
	// 48区
	{
		Start: 0x989E + 0x1,
		End:   0x98EE + 0xE,
	},
	// 49区
	{
		Start: 0x993F + 0x1,
		End:   0x9990 + 0xE,
	},
	// 50区
	{
		Start: 0x999E + 0x1,
		End:   0x99EE + 0xE,
	},
	// 51区
	{
		Start: 0x9A3F + 0x1,
		End:   0x9A90 + 0xE,
	},
	// 52区
	{
		Start: 0x9A9E + 0x1,
		End:   0x9AEE + 0xE,
	},
	// 53区
	{
		Start: 0x9B3F + 0x1,
		End:   0x9B90 + 0xE,
	},
	// 54区
	{
		Start: 0x9B9E + 0x1,
		End:   0x9BEE + 0xE,
	},
	// 55区
	{
		Start: 0x9C3F + 0x1,
		End:   0x9C90 + 0xE,
	},
	// 56区
	{
		Start: 0x9C9E + 0x1,
		End:   0x9CEE + 0xE,
	},
	// 57区
	{
		Start: 0x9D3F + 0x1,
		End:   0x9D90 + 0xE,
	},
	// 58区
	{
		Start: 0x9D9E + 0x1,
		End:   0x9DEE + 0xE,
	},
	// 59区
	{
		Start: 0x9E3F + 0x1,
		End:   0x9E90 + 0xE,
	},
	// 60区
	{
		Start: 0x9E9E + 0x1,
		End:   0x9EEE + 0xE,
	},
	// 61区
	{
		Start: 0x9F3F + 0x1,
		End:   0x9F90 + 0xE,
	},
	// 62区
	{
		Start: 0x9F9E + 0x1,
		End:   0x9FEE + 0xE,
	},
	// 63区
	{
		Start: 0xE03F + 0x1,
		End:   0xE090 + 0xE,
	},
	// 64区
	{
		Start: 0xE09E + 0x1,
		End:   0xE0EE + 0xE,
	},
	// 65区
	{
		Start: 0xE13F + 0x1,
		End:   0xE190 + 0xE,
	},
	// 66区
	{
		Start: 0xE19E + 0x1,
		End:   0xE1EE + 0xE,
	},
	// 67区
	{
		Start: 0xE23F + 0x1,
		End:   0xE290 + 0xE,
	},
	// 68区
	{
		Start: 0xE29E + 0x1,
		End:   0xE2EE + 0xE,
	},
	// 69区
	{
		Start: 0xE33F + 0x1,
		End:   0xE390 + 0xE,
	},
	// 70区
	{
		Start: 0xE39E + 0x1,
		End:   0xE3EE + 0xE,
	},
	// 71区
	{
		Start: 0xE43F + 0x1,
		End:   0xE490 + 0xE,
	},
	// 72区
	{
		Start: 0xE49E + 0x1,
		End:   0xE4EE + 0xE,
	},
	// 73区
	{
		Start: 0xE53F + 0x1,
		End:   0xE590 + 0xE,
	},
	// 74区
	{
		Start: 0xE59E + 0x1,
		End:   0xE5EE + 0xE,
	},
	// 75区
	{
		Start: 0xE63F + 0x1,
		End:   0xE690 + 0xE,
	},
	// 76区
	{
		Start: 0xE69E + 0x1,
		End:   0xE6EE + 0xE,
	},
	// 77区
	{
		Start: 0xE73F + 0x1,
		End:   0xE790 + 0xE,
	},
	// 78区
	{
		Start: 0xE79E + 0x1,
		End:   0xE7EE + 0xE,
	},
	// 79区
	{
		Start: 0xE83F + 0x1,
		End:   0xE890 + 0xE,
	},
	// 80区
	{
		Start: 0xE89E + 0x1,
		End:   0xE8EE + 0xE,
	},
	// 81区
	{
		Start: 0xE93F + 0x1,
		End:   0xE990 + 0xE,
	},
	// 82区
	{
		Start: 0xE99E + 0x1,
		End:   0xE9EE + 0xE,
	},
	// 83区
	{
		Start: 0xEA3F + 0x1,
		End:   0xEA90 + 0xE,
	},
	// 84区
	{
		Start: 0xEA9E + 0x1,
		End:   0xEA9E + 0x6,
	},
}

type RuneMapper map[rune][]SJISCode

func (m RuneMapper) AddCode(code SJISCode) error {
	r := bufio.NewReader(transform.NewReader(bytes.NewReader(code), japanese.ShiftJIS.NewDecoder()))
	c, _, err := r.ReadRune()
	if err != nil {
		return fmt.Errorf("read rune error: %q, %v", code, err)
	}
	//fmt.Printf("%c, %s\n", c, code)
	m[c] = append(m[c], code)
	return nil
}

func (m RuneMapper) AddCodeRange(r CodeRange) error {
	for i := r.Start; i <= r.End; i++ {
		upper := byte(i >> 8)
		lower := byte(i & 0xFF)
		if lower == 0x7F { // sjis には 7Fは現れない
			continue
		}
		code := SJISCode{upper, lower}
		if err := m.AddCode(code); err != nil {
			return err
		}
	}
	return nil
}

func (m RuneMapper) Runes() []rune {
	ret := make([]rune, 0, len(m))
	for k := range m {
		if k == 0xFFFD { // replacement character
			continue
		}
		ret = append(ret, k)
	}
	return ret
}

func NewJIS0213RuneMapper() (*RuneMapper, error) {
	ret := RuneMapper{}
	// latin
	for i := byte(0x20); i <= 0x7e; i++ {
		ret.AddCode(SJISCode{i})
	}
	// ideographic
	for _, v := range JISX0208SJIS {
		if err := ret.AddCodeRange(v); err != nil {
			return nil, err
		}
	}
	return &ret, nil

}

// UnicodeRangeTable returns a range table for specified runes.
func UnicodeRangeTable(runes []rune) *unicode.RangeTable {
	if len(runes) == 0 {
		return nil
	}
	sort.Slice(runes, func(i, j int) bool {
		return runes[i] < runes[j]
	})
	var ret unicode.RangeTable
	for start := 0; start < len(runes); {
		end := start
		for j := start + 1; j < len(runes); j++ {
			if runes[j-1]+1 == runes[j] {
				end = j
				continue
			}
			break
		}
		if runes[end] <= MaxUint16 {
			ret.R16 = append(ret.R16, unicode.Range16{
				Lo:     uint16(runes[start]),
				Hi:     uint16(runes[end]),
				Stride: 1,
			})
		} else {
			ret.R32 = append(ret.R32, unicode.Range32{
				Lo:     uint32(runes[start]),
				Hi:     uint32(runes[end]),
				Stride: 1,
			})
		}
		start = end + 1
	}
	return &ret
}

func DumpRangeTable(w io.Writer, table *unicode.RangeTable) {
	fmt.Fprintln(w, "var RangeTable = &unicode.RangeTable{")
	if len(table.R16) > 0 {
		fmt.Println("\tR16: []unicode.Range16{")
		for _, v := range table.R16 {
			fmt.Fprintf(w, "\t\t{Lo: 0x%X, Hi: 0x%X, Stride: 1},\n", v.Lo, v.Hi)
		}
		fmt.Fprintln(w, "\t},")
	}
	if len(table.R32) > 0 {
		fmt.Println("\tR32: []unicode.Range32{")
		for _, v := range table.R32 {
			fmt.Fprintf(w, "\t\t{Lo: 0x%X, Hi: 0x%X, Stride: 1},\n", v.Lo, v.Hi)
		}
		fmt.Fprintln(w, "\t},")
	}
	fmt.Fprintln(w, "}")
}
