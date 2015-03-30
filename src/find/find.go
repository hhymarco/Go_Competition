package bench

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"runtime"
	"sort"
	"strconv"
	"sync"
)

type Line struct {
	LineNo  int
	LineStr []byte
}

func Find(path, s string) (string, error) {
	if s == "" {
		return "", errors.New("null word")
	}
	f, err := os.Open(path)
	defer f.Close()
	if err == nil {
		bs := []byte(s)
		var wg sync.WaitGroup
		resultPos := make([][]int, runtime.NumCPU())
		chLine := make(chan *Line, runtime.NumCPU())
		for i := 0; i < runtime.NumCPU(); i++ {
			resultPos[i] = make([]int, 0, 200)
			go func(item []int, pos int) {
				for {
					line := <-chLine
					if line == nil {
						continue
					}
					item = SearchStr(line.LineStr, bs, line.LineNo, item)
					resultPos[pos] = item
					wg.Done()
				}
			}(resultPos[i], i)
		}

		rowLens := make([]int, 0, 5000)
		rowLens = rowLens[:0]
		sumLen := 0
		scanner := bufio.NewScanner(f)
		loop := 0
		for scanner.Scan() {
			loop++
			wg.Add(1)
			lineBytes := scanner.Bytes()
			length := len(lineBytes)
			rowLens = append(rowLens, length)
			chLine <- &Line{LineStr: lineBytes, LineNo: sumLen}
			sumLen = sumLen + length
		}
		wg.Wait()
		for j := loop + 1; j <= runtime.NumCPU(); j++ {
			chLine <- nil
		}
		posSum := make([]int, 0, 10000)
		for _, value := range resultPos {
			posSum = append(posSum, value...)
		}
		sort.Ints(posSum)
		//return strconv.Itoa(loop), nil
		posBuf := &bytes.Buffer{}
		sumLen = 0
		for row, rowLen := range rowLens {
			nextPost := 0
			for i := nextPost; i < len(posSum); i++ {
				if posSum[i] >= sumLen && posSum[i] < (sumLen+rowLen) {
					if posBuf.Len() > 0 {
						posBuf.WriteByte(',')
					}
					posBuf.WriteString(strconv.Itoa(row + 1))
					posBuf.WriteByte(':')
					posBuf.WriteString(strconv.Itoa(posSum[i] - sumLen))
					nextPost = i + 1
				}
			}
			sumLen = sumLen + rowLen
		}
		return posBuf.String(), nil
	}
	return "", err
}

func SearchStr(orgin []byte, s []byte, sumLen int, item []int) []int {
	index := bytes.Index(orgin, s)
	increment := 0
	for index > -1 {
		item = append(item, sumLen+index+increment)
		increment = index + 1 + increment
		orgin = orgin[index+1:]
		index = bytes.Index(orgin, s)
	}
	return item
}
