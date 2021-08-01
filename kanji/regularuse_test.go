package kanji

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
		if !IsRegularUse(v) {
			t.Errorf("line=%d, want IsRegularHan(%s)=true, got false", line, string(v))
		}
	}
	if err := s.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestIsNotRegularUseHan(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			name: "OK",
			args: "漢字以外のひらがなやカタカナや😀などもOKとしています!",
			want: false,
		},
		{
			name: "NG",
			args: "勺錘銑脹匁",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.args {
				if got := IsNotRegularUse(v); got != tt.want {
					t.Errorf("IsNotRegularUse(%c) = %v, want %v", v, got, tt.want)
				}
			}
		})
	}
}

func TestIsRegularUse(t *testing.T) {
	tests := []struct {
		name string
		args string
		want bool
	}{
		{
			name: "OK",
			args: "常用漢字挨曖宛嵐畏萎椅彙茨",
			want: true,
		},
		{
			name: "NG",
			args: "ひらがなカタカナ123😀",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, v := range tt.args {
				if got := IsRegularUse(v); got != tt.want {
					t.Errorf("IsRegularUse(%c) = %v, want %v", v, got, tt.want)
				}
			}
		})
	}
}

func TestRegularUseHanDiscriminator_IsNotRegularUseHan(t *testing.T) {
	type fields struct {
		allow    []rune
		disallow []rune
	}
	tests := []struct {
		name   string
		fields fields
		args   string
		want   bool
	}{
		{
			name:   "OK",
			fields: fields{},
			args:   "漢字以外のひらがなやカタカナや😀などもOKとしています!",
			want:   false,
		},
		{
			name: "OK with allow",
			fields: fields{
				allow: []rune{'勺', '錘', '銑', '脹', '匁'},
			},
			args: "勺錘銑脹匁",
			want: false,
		},
		{
			name:   "NG",
			fields: fields{},
			args:   "勺錘銑脹匁",
			want:   true,
		},
		{
			name: "NG with disallow",
			fields: fields{
				disallow: []rune{'漢', '字', '以', '外', 'の', 'ひ', 'ら', 'が', 'な', 'や', 'カ', 'タ', 'カ', 'ナ', 'や', '😀', 'な', 'ど', 'も', 'O', 'K', 'と', 'し', 'て', 'い', 'ま', 'す', '!'},
			},
			args: "漢字以外のひらがなやカタカナや😀などもOKとしています!",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &RegularUseDiscriminator{
				allow:    tt.fields.allow,
				disallow: tt.fields.disallow,
			}
			for _, v := range tt.args {
				if got := d.IsNotRegularUse(v); got != tt.want {
					t.Errorf("IsNotRegularUse(%c) = %v, want %v", v, got, tt.want)
				}
			}
		})
	}
}

func TestRegularUseHanDiscriminator_ReplaceNotRegularUseHanAll(t *testing.T) {
	type fields struct {
		allow    []rune
		disallow []rune
	}
	type args struct {
		s           string
		replacement string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "no matching",
			fields: fields{},
			args: args{
				s:           "漢字以外のひらがなやカタカナや😀などもOKとしています!",
				replacement: "■",
			},
			want: "漢字以外のひらがなやカタカナや😀などもOKとしています!",
		},
		{
			name: "replace",
			fields: fields{
				disallow: []rune{'漢', '😀'},
			},
			args: args{
				s:           "漢字以外のひらがなやカタカナや😀などもOKとしています!",
				replacement: "■",
			},
			want: "■字以外のひらがなやカタカナや■などもOKとしています!",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &RegularUseDiscriminator{
				allow:    tt.fields.allow,
				disallow: tt.fields.disallow,
			}
			if got := d.ReplaceNotRegularUseAll(tt.args.s, tt.args.replacement); got != tt.want {
				t.Errorf("ReplaceNotRegularUseAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
