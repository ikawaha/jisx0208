package jisx0208

import (
	"bufio"
	"os"
	"strings"
	"testing"
)

func TestIs(t *testing.T) {
	tests := []struct {
		name  string
		runes []rune
		want  bool
	}{
		{
			name: "JISX0208",
			runes: []rune{
				'a', 'A', 'ａ', 'Ａ', 'ア', 'ガ', 'ン', 'ッ',
				'高', '崎', '間', '塚', '鴎', 'な', 'ぎ', 'け', 'い',
				'れ', 'ん', '嵜', 'よ', '□', '濱', '渺', '裴', '祐', '元',
				'桑', '￢', '吉', '祇', '曉', '蒋', '☆', '★', '○', '●',
				'◎', '△', '▲', '▽', '▼', '■', '◇', '◆', '緒', '桜',
				'緒', '曖', '/', '\\',
			},
			want: true,
		},
		{
			name: "not JISX0208",
			runes: []rune{
				'髙', '﨑', '閒', '\uFA01', '德', '鷗', '彅', '栁', '琰', '炅',
				'璉', '㟢', '妤', '濵', '裵' /*'祥',*/, '璟', '淼', '姝', '祏',
				'橳', '沅', '吴', '珌', '婧', '桒', '刘', '雯', '昹', '趟',
				'𠮷', '衹', '˘', 'ż', '丂', '龥', '晓', '蔣', '璂', '♦',
				'∮', '緖', '樱', '绪', '韬', '邹', '‼', '⁇', '⁈', '⁉',
				'୨', '୧', '说', '你', '说', '过', '窦', '丽', 'ụ',
				'🙅', '剝',
				// SJIS FA40〜
				'纊', '褜', '鍈', '銈', '蓜', '俉', '炻', '昱', '棈', '鋹', '曻', '彅',
				'丨', '仡', '仼', '伀', '伃', '伹', '佖', '侒', '侊', '侚', '侔', '俍',
				'偀', '倢', '俿', '倞', '偆', '偰', '偂', '傔', '僴', '僘', '兊', '兤',
				'冝', '冾', '凬', '刕', '劜', '劦', '勀', '勛', '匀', '匇', '匤', '卲',
				'厓', '厲', '叝', '﨎', '咜', '咊', '咩', '哿', '喆', '坙', '坥', '垬',
				'埈', '埇', '﨏', '增', '墲', '夋', '奓', '奛', '奝', '奣', '妤',
				'妺', '孖', '寀', '甯', '寘', '寬', '尞', '岦', '岺', '峵', '崧', '嵓',
				'﨑', '嵂', '嵭', '嶸', '嶹', '巐', '弡', '弴', '彧', '德', '忞', '恝',
				'悅', '悊', '惞', '惕', '愠', '惲', '愑', '愷', '愰', '憘', '戓', '抦',
				'揵', '摠', '撝', '擎', '敎', '昀', '昕', '昻', '昉', '昮', '昞', '昤',
				'晥', '晗', '晙', '晳', '暙', '暠', '暲', '暿', '曺', '朎',
				'杦', '枻', '桒', '柀', '栁', '桄', '棏', '﨓', '楨', '﨔', '榘', '槢',
				'樰', '橫', '橆', '橳', '橾', '櫢', '櫤', '毖', '氿', '汜', '沆', '汯',
				'泚', '洄', '涇', '浯', '涖', '涬', '淏', '淸', '淲', '淼', '渹', '湜',
				'渧', '渼', '溿', '澈', '澵', '濵', '瀅', '瀇', '瀨', '炅', '炫', '焏',
				'焄', '煜', '煆', '煇', '凞', '燁', '燾', '犱', '犾', '猤', '獷',
				'玽', '珉', '珖', '珣', '珒', '琇', '珵', '琦', '琪', '琩', '琮', '瑢',
				'璉', '璟', '甁', '畯', '皂', '皜', '皞', '皛', '皦', '睆', '劯',
				'砡', '硎', '硤', '硺', '礰', '禔', '禛', '竑',
				'竧', '竫', '箞', '絈', '絜', '綷', '綠', '緖', '繒', '罇',
				'羡', '茁', '荢', '荿', '菇', '菶', '葈', '蒴', '蕓', '蕙', '蕫',
				'﨟', '薰', '蘒', '﨡', '蠇', '裵', '訒', '訷', '詹', '誧', '誾', '諟',
				'諶', '譓', '譿', '賰', '賴', '贒', '赶', '﨣', '軏', '﨤',
				'遧', '郞', '鄕', '鄧', '釚', '釗', '釞', '釭', '釮', '釤', '釥',
				'鈆', '鈐', '鈊', '鈺', '鉀', '鈼', '鉎', '鉙', '鉑', '鈹', '鉧', '銧',
				'鉷', '鉸', '鋧', '鋗', '鋙', '鋐', '﨧', '鋕', '鋠', '鋓', '錥', '錡',
				'鋻', '﨨', '錞', '鋿', '錝', '錂', '鍰', '鍗', '鎤', '鏆', '鏞', '鏸',
				'鐱', '鑅', '鑈', '閒', '﨩', '隝', '隯', '霳', '霻', '靃', '靍',
				'靏', '靑', '靕', '顗', '顥', '餧', '馞', '驎', '髙',
				'髜', '魵', '魲', '鮏', '鮱', '鮻', '鰀', '鵰', '鵫', '鸙', '黑',
			},
			want: false,
		},
		{
			name:  "Halfwidth and Fullwidth Forms (U+FF00〜U+FFEF)",
			runes: []rune{'￤'},
			want:  false,
		},
		{
			name: "Number Forms",
			runes: []rune{
				'ⅰ', 'ⅱ', 'ⅲ', 'ⅳ', 'ⅴ', 'ⅵ', 'ⅶ', 'ⅷ', 'ⅸ', 'ⅹ',
				'Ⅰ', 'Ⅱ', 'Ⅲ', 'Ⅳ', 'Ⅴ', 'Ⅵ', 'Ⅶ', 'Ⅷ', 'Ⅸ', 'Ⅹ',
			},
			want: false,
		},
		{
			name: "Katakana Phonetic Extensions",
			runes: []rune{
				'ㇰ', 'ㇱ', 'ㇲ', 'ㇳ', 'ㇴ', 'ㇵ', 'ㇶ', 'ㇷ', 'ㇸ', 'ㇹ', 'ㇺ', 'ㇻ', 'ㇼ', 'ㇽ', 'ㇾ', 'ㇿ',
			},
			want: false,
		},
		{
			name:  "Letterlike Symbols",
			runes: []rune{'™', '℠', '©', '№', '℡'},
			want:  false,
		},
		{
			name: "Enclosed Alphanumerics",
			runes: []rune{
				'⓪', '①', '②', '③', '④', '⑤', '⑥', '⑦', '⑧', '⑨',
				'ⓐ', 'ⓑ', 'ⓒ', 'ⓓ', 'ⓔ', 'ⓕ', 'ⓖ',
				'Ⓗ', 'Ⓘ', 'Ⓙ', 'Ⓚ', 'Ⓛ', 'Ⓜ', 'Ⓝ',
				'🄾', '🄿', '🅀', '🅁', '🅂', '🅃', '🅄',
				'🅥', '🅦', '🅧', '🅨', '🅩',
				'🅰', '🅱', '🅲', '🅳', '🅴', '🅵', '🅶',
				'⒣', '⒤', '⒥', '⒦', '⒧', '⒨', '⒩',
				'🄞', '🄟', '🄠', '🄡', '🄢', '🄣', '🄤',
				'🇻', '🇼', '🇽', '🇾', '🇿',
			},
			want: false,
		},
		{
			name: "Enclosed CJK Letters and Months",
			runes: []rune{
				'㋀', '㋁', '㋂', '㋃', '㋄', '㋅', '㋆', '㋇', '㋈', '㋉', '㋊', '㋋',
				'㈠', '㈡', '㈢', '㈣', '㈤', '㈥', '㈦', '㈧', '㈨', '㈩',
				'㊀', '㊁', '㊂', '㊃', '㊄', '㊅', '㊆', '㊇', '㊈', '㊉',
				'㋐', '㋑', '㋒', '㋓', '㋔', '㋕', '㋖', '㋗', '㋘', '㋙',
				'㋿', '㍻', '㍼', '㍽', '㍾',
				'㈪', '㈫', '㈬', '㈭', '㈮', '㈯', '㈰',
				'㊊', '㊋', '㊌', '㊍', '㊎', '㊏', '㊐',
				'㈱', '㈲', '㈳', '㈴', '㈵', '㈶', '㈷', '㈸', '㈹', '㈺',
				'㈻', '㈼', '㈽', '㈾', '㈿', '㉀', '㉁', '㉂', '㉃',
				'㉄', '㉅', '㉆', '㉇', '㊑', '㊒', '㊓', '㊔', '㊕', '㊖',
				'㊗', '㊘', '㊙', '㊚', '㊛', '㊜', '㊝', '㊞', '㊟', '㊠',
				'㊡', '㊢', '㊣', '㊩', '㊪', '㊫', '㊬', '㊭', '㊮', '㊯', '㊰',
				'㊤', '㊥', '㊦', '㊧', '㊨',
			},
			want: false,
		},
		{
			name: "CJK Compatibility",
			runes: []rune{
				'㍘', '㍙', '㍚', '㍛', '㍜', '㍝', '㍞', '㍟', '㍠', '㍡',
				'㏠', '㏡', '㏢', '㏣', '㏤', '㏥', '㏦', '㏧', '㏨', '㏩',
				'㌀', '㌁', '㌂', '㌃', '㌄', '㌅', '㌆', '㌇', '㌈', '㌉', '㌊',
				'㍿',
			},
			want: false,
		},
		{
			name: "Enclosed Ideographic Supplemente",
			runes: []rune{
				'🈂', '🈁', '🉐', '🉑',
				'🉀', '🉁', '🉂', '🉃', '🉄', '🉅', '🉆', '🉇', '🉈',
			},
			want: false,
		},
	}
	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			for _, r := range v.runes {
				if got := Is(r); got != v.want {
					t.Errorf("got %v, want %v, rune=[%c]", got, v.want, r)
				}
			}
		})
	}
}

func TestIs_Golden(t *testing.T) {
	f, err := os.Open("./testdata/golden.txt")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	defer f.Close() //nolint:gosec
	s := bufio.NewScanner(f)
	var line int
	for s.Scan() {
		line++
		txt := s.Text()
		if txt == "" {
			t.Errorf("invalid golden data, line=%d, %s", line, txt)
			continue
		}
		if strings.HasPrefix(txt, "!!!!") {
			t.Log(txt)
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

func TestToValid(t *testing.T) {
	type args struct {
		s           string
		replacement string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ascii",
			args: args{
				s:           "abc1234",
				replacement: " ",
			},
			want: "abc1234",
		},
		{
			name: "invalid utf8",
			args: args{
				s:           "abc\xFF1234",
				replacement: " ",
			},
			want: "abc 1234",
		},
		{
			name: "not JIS X 0213",
			args: args{
				s:           "髙﨑閒\uFA01德鷗彅栁琰炅璉㟢妤",
				replacement: "□",
			},
			want: "□□□□□□□□□□□□□",
		},
		{
			name: "JIS X 0213",
			args: args{
				s:           "人魚は、南の方の海にばかり棲んでいるのではありません。",
				replacement: "□",
			},
			want: "人魚は、南の方の海にばかり棲んでいるのではありません。",
		},
		{
			name: "replace",
			args: args{
				s:           "人魚は、\xFF南の方の海にばかり髙棲んでいるのではありません。",
				replacement: "□",
			},
			want: "人魚は、□南の方の海にばかり□棲んでいるのではありません。",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToValid(tt.args.s, tt.args.replacement); got != tt.want {
				t.Errorf("ToValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_AmbiguousRunes(t *testing.T) {
	tests := []struct {
		rune rune
		want bool
	}{
		{
			rune: '¥',
			want: false,
		},
	}
	for _, v := range tests {
		if got := Is(v.rune); got != v.want {
			t.Errorf("Is(%c) = %v, want %v", v.rune, got, v.want)
		}
	}
}

func TestDiscriminator_Is(t *testing.T) {
	t.Run("allow", func(t *testing.T) {
		allow := []rune{'¥'}
		disallow := []rune(nil)
		d := NewDiscriminator(Allow(allow...))
		tests := []struct {
			rune rune
			want bool
		}{
			{
				rune: '¥',
				want: true,
			},
		}
		for _, v := range tests {
			if got := d.Is(v.rune); got != v.want {
				t.Errorf("d.Is(%c) = %v, allow=%+v, disallow=%+v, want %v", v.rune, got, allow, disallow, v.want)
			}
		}
	})
	t.Run("disallow", func(t *testing.T) {
		allow := []rune(nil)
		disallow := []rune{'あ'}
		d := NewDiscriminator(Disallow(disallow...))
		tests := []struct {
			rune rune
			want bool
		}{
			{
				rune: 'あ',
				want: false,
			},
		}
		for _, v := range tests {
			if got := d.Is(v.rune); got != v.want {
				t.Errorf("d.Is(%c) = %v, allow=%+v, disallow=%+v, want %v", v.rune, got, allow, disallow, v.want)
			}
		}
	})
	t.Run("allow and disallow", func(t *testing.T) {
		allow := []rune{'髙'}
		disallow := []rune{'あ'}
		d := NewDiscriminator(Allow(allow...), Disallow(disallow...))
		tests := []struct {
			rune rune
			want bool
		}{
			{
				rune: 'あ',
				want: false,
			},
			{
				rune: '髙',
				want: true,
			},
		}
		for _, v := range tests {
			if got := d.Is(v.rune); got != v.want {
				t.Errorf("d.Is(%c) = %v, allow=%+v, disallow=%+v, want %v", v.rune, got, allow, disallow, v.want)
			}
		}
	})
}

func TestDiscriminator_ToValid(t *testing.T) {
	t.Run("unspecified", func(t *testing.T) {
		d := NewDiscriminator()
		type args struct {
			s           string
			replacement string
		}
		tests := []struct {
			name string
			args args
			want string
		}{
			{
				name: "ascii",
				args: args{
					s:           "abc1234",
					replacement: " ",
				},
				want: "abc1234",
			},
			{
				name: "invalid utf8",
				args: args{
					s:           "abc\xFF1234",
					replacement: " ",
				},
				want: "abc 1234",
			},
			{
				name: "not JIS X 0213",
				args: args{
					s:           "髙﨑閒\uFA01德鷗彅栁琰炅璉㟢妤",
					replacement: "□",
				},
				want: "□□□□□□□□□□□□□",
			},
			{
				name: "JIS X 0213",
				args: args{
					s:           "人魚は、南の方の海にばかり棲んでいるのではありません。",
					replacement: "□",
				},
				want: "人魚は、南の方の海にばかり棲んでいるのではありません。",
			},
			{
				name: "replace",
				args: args{
					s:           "人魚は、\xFF南の方の海にばかり髙棲んでいるのではありません。",
					replacement: "□",
				},
				want: "人魚は、□南の方の海にばかり□棲んでいるのではありません。",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := d.ToValid(tt.args.s, tt.args.replacement); got != tt.want {
					t.Errorf("ToValid() = %v, want %v", got, tt.want)
				}
			})
		}
	})
	t.Run("specified", func(t *testing.T) {
		d := NewDiscriminator(Allow('髙'), Disallow('魚'))
		type args struct {
			s           string
			replacement string
		}
		tests := []struct {
			name string
			args args
			want string
		}{
			{
				name: "ascii",
				args: args{
					s:           "abc1234",
					replacement: " ",
				},
				want: "abc1234",
			},
			{
				name: "invalid utf8",
				args: args{
					s:           "abc\xFF1234",
					replacement: " ",
				},
				want: "abc 1234",
			},
			{
				name: "not JIS X 0213",
				args: args{
					s:           "髙﨑閒\uFA01德鷗彅栁琰炅璉㟢妤",
					replacement: "□",
				},
				want: "髙□□□□□□□□□□□□",
			},
			{
				name: "JIS X 0213",
				args: args{
					s:           "人□は、南の方の海にばかり棲んでいるのではありません。",
					replacement: "□",
				},
				want: "人□は、南の方の海にばかり棲んでいるのではありません。",
			},
			{
				name: "replace",
				args: args{
					s:           "人魚は、\xFF南の方の海にばかり髙棲んでいるのではありません。",
					replacement: "□",
				},
				want: "人□は、□南の方の海にばかり髙棲んでいるのではありません。",
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := d.ToValid(tt.args.s, tt.args.replacement); got != tt.want {
					t.Errorf("ToValid() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}
