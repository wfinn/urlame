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
			t.Logf("urls didn't match\nGot:\n%s\nExpected:\n%s", outputBuf.String(), string(expected))
			t.Fail()
		}
	}
}

func staticTest(testname, input, expected string, t *testing.T) {
	if expected[0] == '\n' {
		expected = expected[1:]
	}
	outputBuf := bytes.NewBufferString("")
	runurlame(bytes.NewBufferString(input), outputBuf, false)
	output := outputBuf.String()
	if expected != output {
		t.Logf("'%s' failed\nGot:\n%s\nExpected:\n%s", testname, output, expected)
		t.Fail()
	}
}

func TestLameThings(t *testing.T) {
	urls := `
http://example.com/notlame
http://example.com/lame.png?x=y
http://example.com/docs/boring
https://example.org/user/rick
https://example.org/user/morty
https://example.org/10-things-you-already-knew
`
	expected := `
http://example.com/notlame
`
	staticTest("images, user profiles, blogs etc. are ignored", urls, expected, t)
}

func TestProtocolsAreIgnored(t *testing.T) {
	urls := `
http://localhost/foo
https://localhost/foo
htTPS://localhost/foo
`
	expected := `
http://localhost/foo
`
	staticTest("protocols are ignored", urls, expected, t)
}

func TestTrailingJunk(t *testing.T) {
	urls := `
http://localhost/foo
http://localhost/foo/
http://localhost/foo/?
http://localhost/foo/#
http://localhost/foo/?#
`
	expected := `
http://localhost/foo
`
	staticTest("trailing / ? # are ignored", urls, expected, t)
}

func TestQueryValues(t *testing.T) {
	urls := `
http://localhost/foo?param=1
http://localhost/foo?param=something-else
http://localhost/foo?otherparam=1
`
	expected := `
http://localhost/foo?param=1
http://localhost/foo?otherparam=1
`
	staticTest("query parameter values are ignored", urls, expected, t)
}

func TestLameFilenames(t *testing.T) {
	urls := `
http://localhost/
http://localhost/robots.txt
http://localhost/robots.txt?newparam
`
	expected := `
http://localhost/
http://localhost/robots.txt?newparam
`
	staticTest("common lame files are ignored", urls, expected, t)
}

func TestLameParams(t *testing.T) {
	urls := `
http://localhost/foobar
http://localhost/foobar?utm_source=twitter
`
	expected := `
http://localhost/foobar
`
	staticTest("common lame params are ignored", urls, expected, t)
}
