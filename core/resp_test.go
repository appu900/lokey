package core

import (
	"testing"
)

func TestErrorMessageDecode(t *testing.T) {
	cases := map[string]string{
		"-error-message\r\n":"error-message",
	}
	for key, value := range cases{
		testIncomingMessage , _ := Decode([]byte(key))
		if testIncomingMessage != value{
			t.Fail()
		}
	}
}
func TestSimpleStringDecode(t *testing.T) {
	cases := map[string]string{
		"+OK\r\n": "OK",
	}
	for key, value := range cases {
		testIncomingValue, _ := Decode([]byte(key))
		if testIncomingValue != value {
			t.Fail()
		}
	}
}
