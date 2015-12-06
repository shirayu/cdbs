package cdbs

import (
	"testing"
)

func TestDigit(t *testing.T) {
	tests := []struct {
		num   int
		digit int
	}{
		{0, 1},
		{5, 1},
		{12, 2},
		{123, 3},
		{1234, 4},
		{12345, 5},
	}

	for _, test := range tests {
		if GetDigitNum(test.num) != test.digit {
			t.Errorf("Expected %d but got: %d", test.digit, test.num)
		}
	}
}

func TestOutput(t *testing.T) {
	//TODO
	//     tmpdir, err := ioutil.TempDir()
	//     if err != nil {
	//         t.Errorf("Expected no error but got: %v", err)
	//     }

}
