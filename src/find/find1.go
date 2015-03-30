package bench

import (
	"bufio"
	"bytes"
	"errors"
	//"io"
	"os"
	"regexp"
	"strconv"
	//"strings"
)

func Find1(path, s string) (string, error) {
	if s == "" {
		return "", errors.New("")
	}
	f, err := os.Open(path)
	defer f.Close()
	if err == nil {
		var result bytes.Buffer
		r := bufio.NewReader(f)
		re := regexp.MustCompile("\n")
		sLoc := re.FindReaderIndex(r)
		for sLoc != nil {
			result.WriteString(strconv.Itoa(sLoc[0]))
			result.WriteString(":")
			result.WriteString(strconv.Itoa(sLoc[1]))
			sLoc = re.FindReaderIndex(r)
		}
		r.Reset(f)
		re = regexp.MustCompile("\n")
		//	nLoc := re.FindReaderIndex(r)

		return result.String(), nil
	}
	return "", err
}
