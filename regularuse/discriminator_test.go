package regularuse

import (
	"bufio"
	"os"
	"testing"
)

func TestIs_Golden(t *testing.T) {
	f, err := os.Open("./testdata/golden_jyouyou_H22-11-30.csv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	var line int
	for s.Scan() {
		line++
		txt := s.Text()
		if txt == "" {
			t.Errorf("invalid golden data, line=%d, %s", line, txt)
			continue
		}
		v := []rune(txt)[0]
		if !Is(v) {
			t.Errorf("line=%d, want Is(%s)=true, got false", line, string(v))
		}
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIsNotRegularIdeographic(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		input := "常用漢字表にのっている 漢字　だけが入っている文章はOKとしています。123Abc!"
		for _, r := range input {
			if IsNotRegularIdeographic(r) {
				t.Errorf("got IsNotRegularIdeographic(%c) = false, want true", r)
			}
		}
	})

	t.Run("ng", func(t *testing.T) {
		input := "勺錘銑脹匁"
		for _, r := range input {
			if !IsNotRegularIdeographic(r) {
				t.Errorf("got IsNotRegularIdeographic(%c) = true, want false", r)
			}
		}
	})
}
