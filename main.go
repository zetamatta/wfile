package wfile

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func guess(fname string, bin []byte, features []*Feature) string {
	for _, f := range features {
		if f.Offset+len(f.Magic) < len(bin) && bytes.Equal(bin[f.Offset:f.Offset+len(f.Magic)], f.Magic) {
			if f.Func != nil {
				return f.Func(fname, bin)
			} else {
				return f.Desc
			}
		}
	}
	return ""
}

func Report(fname string) (string, error) {
	fd, err := os.Open(fname)
	if err != nil {
		return "", err
	}
	defer fd.Close()

	if stat, err := fd.Stat(); err != nil {
		return "", err
	} else if stat.IsDir() {
		return "", fmt.Errorf("%s: Directory", fname)
	}

	bin := make([]byte, 1024)
	n, err := fd.Read(bin)
	if err == io.EOF {
		return "", fmt.Errorf("%s: zero byte file\n", fname)
	}
	if err != nil {
		return "", err
	}
	bin = bin[:n]

	suffix := strings.TrimPrefix(strings.ToLower(filepath.Ext(fname)), ".")
	if features, ok := suffixTable[suffix]; ok {
		if result := guess(fname, bin, features); result != "" {
			return result, nil
		}
	}
	if result := guess(fname, bin, flatTable); result != "" {
		return result, nil
	}
	for _, features := range suffixTable {
		if result := guess(fname, bin, features); result != "" {
			return result, nil
		}
	}
	return TryText(fname, bin), nil
}
