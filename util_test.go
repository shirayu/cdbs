package cdbs

import (
	"io"
	"testing"
)

func TestGet(t *testing.T) {
	cdbs, err := NewCdbs("data/sample")
	if err != nil {
		t.Errorf("Expected no error but got: %v", err)
	}

	var tests = []struct {
		arg   string
		err   error
		value string
	}{
		{
			arg:   "aaaa",
			err:   io.EOF,
			value: "",
		},
		{
			arg:   "Akey0",
			err:   nil,
			value: "value0",
		},
		{
			arg:   "Akey1",
			err:   nil,
			value: "value1",
		},
		{
			arg:   "Akey2",
			err:   nil,
			value: "value2",
		},
		{
			arg:   "Bkey0",
			err:   nil,
			value: "value5",
		},
		{
			arg:   "Bkey1",
			err:   nil,
			value: "value7",
		},
		{
			arg:   "Bkey2",
			err:   nil,
			value: "value9",
		},
		{
			arg:   "ZZZZ",
			err:   io.EOF,
			value: "",
		},
	}

	for _, test := range tests {
		ret, err := cdbs.Get(test.arg)
		if test.err != err {
			t.Errorf("Expected error %v for the key [%v] but got: %v", test.err, test.arg, err)
		}
		if test.value != string(ret) {
			t.Errorf("Expected value [%s] for the key [%v] but got: [%s]", test.value, test.arg, ret)
		}
	}
}
