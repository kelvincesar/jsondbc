package canconv

import (
	"os"
	"strconv"
)

func FormatFloat(val float64) string {
	return strconv.FormatFloat(val, 'f', -1, 64)
}

func FormatString(val string) string {
	return "\"" + val + "\""
}

func FormatUint(val uint32) string {
	return strconv.FormatUint(uint64(val), 10)
}

type file struct {
	f *os.File
}

func newFile(f *os.File) *file {
	return &file{f: f}
}

func (f *file) print(str ...string) {
	tmp := ""
	for i, s := range str {
		if i == 0 {
			tmp += s
			continue
		}
		tmp += " " + s
	}
	_, err := f.f.WriteString(tmp + "\n")
	if err != nil {
		panic(err)
	}
}
