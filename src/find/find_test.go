package bench

import "testing"

// path is a text file path.
const path = "/Users/ghost/codesource/go_benchmark_competiton/src/find/data1.txt"

// word is a word to find.
const word = "aa"

// want is an expected number.
const want = "1:0,1:4,3:0,3:1"

/*
func TestSearch(t *testing.T) {
	item := make([]int, 0, 10)
	item = SearchStr([]byte("aaabbaa"), []byte("aa"), 0, item)
	item = SearchStr([]byte("ccacaa"), []byte("aa"), 7, item)
	for _, value := range item {
		t.Error(value)
	}
}
*/

func TestFind(t *testing.T) {
	got, err := Find(path, word)
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Errorf("Find(%q, %q) => %q, want %q", path, word, got, want)
	}
}

func TestFind_emptyWord(t *testing.T) {
	if _, err := Find(path, ""); err == nil {
		t.Error("some kind of error should be returned")
	}
}

func TestFind_emptyResult(t *testing.T) {
	got, err := Find(path, "not_exist_word")
	if err != nil {
		t.Fatal(err)
	}
	if got != "" {
		t.Error("empty value should be returned")
	}
}

func BenchmarkFind(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Find(path, word)
	}
}
