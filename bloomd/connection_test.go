package bloomd

import (
	"testing"
)

func TestCreateSocket(t *testing.T) {
	conn := Connection{Server: serverAddress}
	conn.createSocket()
	if conn.Socket == nil {
		t.Fail()
	}
}

func TestSend(t *testing.T) {
	conn := Connection{Server: serverAddress}
	err := conn.Send("derp")
	failIfError(t, err)
}

func TestReadDerp(t *testing.T) {
	conn := Connection{Server: serverAddress}
	err := conn.Send("list")
	failIfError(t, err)
	resp, err := conn.Read()
	failIfError(t, err)
	if resp != "START" {
		t.Error("Got: " + resp)
	}
}
