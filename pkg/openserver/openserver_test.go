package openserver

import "testing"

func TestCreateServerValidPort(t *testing.T) {
	_, err := NewServer(1999)
	if err != nil {
		t.Error("Server init threw error")
	}
}

func TestCreateServerInvalidPort(t *testing.T) {
	_, err := NewServer(80)
	if err == nil {
		t.Error("Server init should not be permited on port <=1000")
	}
}
