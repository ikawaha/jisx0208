package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	jisx0208htmlURL = "https://www.asahi-net.or.jp/~ax2s-kmtn/ref/jisx0208.html"
)

func OpenGoldenSrc(path string) (io.ReadCloser, error) {
	if !strings.HasPrefix(path, "https://") {
		return os.Open(path)
	}
	res, err := http.Get(path)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewBuffer(b)), nil
}

func MakeGolden(w io.Writer, werr io.Writer) error {
	r, err := OpenGoldenSrc(jisx0208htmlURL)
	if err != nil {
		return fmt.Errorf("cannot open golden src: %w", err)
	}
	defer r.Close()
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return fmt.Errorf("invalid document src: %w", err)
	}
	doc.Find("table.basic").Each(func(_ int, s *goquery.Selection) {
		blockTitle, ok := s.Attr("summary")
		if !ok {
			blockTitle = "ブロックタイトル不明"
		}
		fmt.Println("!!!! " + blockTitle)
		s.Find("td").Each(func(_ int, s *goquery.Selection) {
			switch {
			case s.HasClass("cha"):
				fmt.Fprintf(w, "%s\n", s.Text())
			case s.HasClass("cha_jis83"), s.HasClass("cha_jis90"), s.HasClass("cha_jis8390"):
				fmt.Fprintf(w, "%s\t%s\n", s.Text(), s.AttrOr("title", "???"))
			case s.HasClass("v"), s.HasClass("blnk"):
				//nop
			default:
				val := "UNKNOWN"
				if s.Nodes != nil && s.Nodes[0] != nil && s.Nodes[0].Attr != nil {
					val = s.Nodes[0].Attr[0].Val
				}
				fmt.Fprintf(werr, "%q\t%s\n", s.Text(), val)
			}
		})
	})
	return nil
}
