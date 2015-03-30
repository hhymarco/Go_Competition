package bench1

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type Line struct {
	LineNo  int
	LineStr string
}
type Set struct {
	m map[int]string
	sync.RWMutex
}

func NewSet() *Set {
	return &Set{
		m: map[int]string{},
	}
}
func (s *Set) Add(key int, value string) {
	s.Lock()
	defer s.Unlock()
	s.m[key] = value
}
func (s *Set) Has(key int) bool {
	s.RLock()
	defer s.RUnlock()
	_, ok := s.m[key]
	return ok
}
func (s *Set) Value(key int) (string, bool) {
	s.RLock()
	defer s.RUnlock()
	value, ok := s.m[key]
	return value, ok
}
func (s *Set) Size() int {
	return len(s.m)
}
func (s *Set) KeySet() []int {
	s.RLock()
	defer s.RUnlock()
	list := []int{}
	for key, _ := range s.m {
		list = append(list, key)
	}
	return list
}

func Find(path, s string) (string, error) {
	if s == "" {
		return "", errors.New("")
	}
	f, err := os.Open(path)
	defer f.Close()
	if err == nil {
		r := bufio.NewReader(f)
		var wg sync.WaitGroup
		resultMap := NewSet()
		chLine := make(chan *Line, 10)
		for i := 0; i < 10; i++ {
			go func() {
				for {
					line := <-chLine
					searchStr(line.LineStr, s, line.LineNo, resultMap)
					wg.Done()
				}
			}()
		}

		lineN := 0
		for {
			line, err := r.ReadString('\n')
			if err != nil || err == io.EOF {
				break
			}
			wg.Add(1)
			lineN++
			chLine <- &Line{LineStr: line, LineNo: lineN}
		}
		wg.Wait()
		var posBuf bytes.Buffer
		keys := resultMap.KeySet()
		sort.Ints(keys)
		for _, key := range keys {
			v, _ := resultMap.Value(key)
			posBuf.WriteString(v)
		}
		posStr := posBuf.String()
		if strings.HasSuffix(posStr, ",") {
			posStr = posStr[:len(posStr)-1]
		}
		return posStr, nil
	}
	return "", err
}

func searchStr(orgin string, s string, lineNo int, m *Set) {
	result := ""
	index := strings.Index(orgin, s)
	increment := 0
	for index > -1 {
		result = strings.Join([]string{result, strconv.Itoa(lineNo), ":", strconv.Itoa(index + increment), ","}, "")
		//	result = result + strconv.Itoa(lineNo) + ":" + strconv.Itoa(index+increment) + ","
		increment = index + 1 + increment
		//	orgin = substr(orgin, index+1, len(orgin)-index-1)
		orgin = orgin[index+1:]
		index = strings.Index(orgin, s)
	}
	if result != "" {
		m.Add(lineNo, result)
	}
}
