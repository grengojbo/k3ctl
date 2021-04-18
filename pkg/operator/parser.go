package operator

import (
	"fmt"
	"strconv"
	"strings"

	"errors"

	log "github.com/sirupsen/logrus"
)

// ParseDfOutput парсим строку комманды df -P
// TODO: translate
func ParseDfOutput(out string) (float64, error) {
	outlines := strings.Split(out, "\n")
	l := len(outlines)
	var total, used float64 = 0, 0
	for _, line := range outlines[1 : l-1] {
		parsedLine := strings.Fields(line)
		t, err := strconv.ParseFloat(parsedLine[1], 64)
		if err != nil {
			return 0, err
		}
		u, err := strconv.ParseFloat(parsedLine[2], 64)
		if err != nil {
			return 0, err
		}
		total += t
		used += u
	}
	return used / total, nil
}

// ParseInt64Output парсим одну строку с числом
// TODO: translate
func ParseInt64Output(out string) (val int64, err error) {
	outlines := strings.Split(out, "\n")
	l := len(outlines)
	if l > 2 {
		return 0, errors.New(fmt.Sprintf("Output line count %d > 1", l-1))
	}
	u, err := strconv.ParseInt(outlines[0], 10, 64)
	if err != nil {
		return 0, err
	}
	// log.Debugf("line => %v", u)
	return u, nil
}

func ParseOutput(out string) (val int64, err error) {
	outlines := strings.Split(out, "\n")
	l := len(outlines)
	// log.Errorln("[OUT] ", out, "len=", l)
	// for i, line := range outlines[1 : l-1] {
	for i, line := range outlines[0 : l-1] {
		parsedLine := strings.Fields(line)
		log.Debugf("line [%d] => %s", i, parsedLine)
	}
	return 0, nil
}
