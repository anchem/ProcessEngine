package util

import (
	"runtime"
	"strings"
)

const (
	CUR_OS = runtime.GOOS
)

func GetSp() string {
	if strings.EqualFold(CUR_OS, "windows") {
		return "\\"
	} else {
		return "/"
	}
}
