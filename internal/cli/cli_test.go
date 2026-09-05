package cli

import "testing"

func TestNewCliValidatesArgumentsBeforeConnecting(t *testing.T) {
	validHost, validPassword, validCharset := "localhost", "", "latin1"
	for _, port := range []int32{0, 65536} {
		if _, err := NewCli(&validHost, &port, &validPassword, &validCharset); err == nil {
			t.Errorf("port %d was accepted", port)
		}
	}
	badCharset := "cp1252"
	port := int32(25575)
	if _, err := NewCli(&validHost, &port, &validPassword, &badCharset); err == nil {
		t.Fatal("invalid charset was accepted")
	}
	if _, err := NewCli(nil, &port, &validPassword, &validCharset); err == nil {
		t.Fatal("nil host was accepted")
	}
}
