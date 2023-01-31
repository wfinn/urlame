package main

import (
	"bytes"
	"os"
	"testing"
)

func TestBasics(t *testing.T) {
	outputBuf := bytes.NewBufferString("")
	if urlfile, err := os.Open("test_urls/basic.txt"); err == nil {
		runurlame(urlfile, outputBuf, false)
	} else {
		t.Fail()
	}
	if expected, err := os.ReadFile("test_urls/basic_expected.txt"); err != nil {
		t.Fail()
	} else {
		if !bytes.Equal(expected, outputBuf.Bytes()) {
			t.Logf("urls didn't match\n%s\nExpected:\n%s", outputBuf.String(), string(expected))
			t.Fail()
		}
	}
}
