package template

import (
	"container/list"
	"fmt"
	"html/template"
	"runtime"
	"strings"
	"time"

	"github.com/credli/castor/modules/base"
	"github.com/credli/hcsg/settings"
)

var Funcs template.FuncMap = map[string]interface{}{
	"GoVer": func() string {
		return strings.Title(runtime.Version())
	},
	"UseHTTPS": func() bool {
		return strings.HasPrefix(settings.AppURL, "https")
	},
	"AppName": func() string {
		return settings.AppName
	},
	"AppSubURL": func() string {
		return settings.AppSubURL
	},
	"AppVer": func() string {
		return settings.AppVer
	},
	"loadTimes": func(startTime time.Time) string {
		return fmt.Sprint(time.Since(startTime).Nanoseconds()/1e6) + "ms"
	},
	"Safe":     Safe,
	"Str2html": Str2html,
	"DateFmtLong": func(t time.Time) string {
		return t.Format(time.RFC1123Z)
	},
	"DateFmtShort": func(t time.Time) string {
		return t.Format("Jan 02, 2006")
	},
	"Add": func(a, b int) int {
		return a + b
	},
	"List": List,
	"SubStr": func(str string, start, length int) string {
		if len(str) == 0 {
			return ""
		}
		end := start + length
		if length == -1 {
			end = len(str)
		}
		if len(str) < end {
			return str
		}
		return str[start:end]
	},
	"Sha1": Sha1,
	"EscapePound": func(str string) string {
		return strings.Replace(strings.Replace(str, "%", "%25", -1), "#", "%23", -1)
	},
}

func Safe(raw string) template.HTML {
	return template.HTML(raw)
}

func Str2html(raw string) template.HTML {
	return template.HTML(base.Sanitizer.Sanitize(raw))
}

func Range(l int) []int {
	return make([]int, l)
}

func List(l *list.List) chan interface{} {
	e := l.Front()
	c := make(chan interface{})
	go func() {
		for e != nil {
			c <- e.Value
			e = e.Next()
		}
		close(c)
	}()
	return c
}

func Sha1(str string) string {
	return base.EncodeSha1(str)
}

// Replaces all prefixes 'old' in 's' with 'new'.
func ReplaceLeft(s, old, new string) string {
	old_len, new_len, i, n := len(old), len(new), 0, 0
	for ; i < len(s) && strings.HasPrefix(s[i:], old); n += 1 {
		i += old_len
	}

	// simple optimization
	if n == 0 {
		return s
	}

	// allocating space for the new string
	newLen := n*new_len + len(s[i:])
	replacement := make([]byte, newLen, newLen)

	j := 0
	for ; j < n*new_len; j += new_len {
		copy(replacement[j:j+new_len], new)
	}

	copy(replacement[j:], s[i:])
	return string(replacement)
}
