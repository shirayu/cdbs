package cdbs

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"testing"
)

func TestOutput(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "cdbs_test")
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}

	lines := `aaa,3
bar,4
`
	prefix := "out"
	linesMap := `aaa 0
bar 1
`

	r := bufio.NewReader(bytes.NewBufferString(lines))
	Output(r, path.Join(tmpdir, prefix), false, ',', false, 1)

	for i := 0; i < 2; i++ {
		filename := path.Join(tmpdir, prefix) + "." + strconv.Itoa(i) + ".cdb"
		if _, err := os.Stat(filename); err != nil {
			t.Errorf("Not exist: %v", err)
		}
	}

	filenameMap := path.Join(tmpdir, prefix) + ".keymap"
	text, err := ioutil.ReadFile(filenameMap)
	if err == nil {
		if string(text) != linesMap {
			t.Errorf("Expected %s but %s", linesMap, text)
		}
	} else {
		t.Errorf("Open error: %v", err)
	}

	if err := os.RemoveAll(tmpdir); err != nil {
		t.Errorf("Remove error: %v", err)
	}
}
