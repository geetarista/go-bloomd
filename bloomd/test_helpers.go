package bloomd

import (
	"testing"
)

var (
	serverHost    = "127.0.0.1"
	serverPort    = "8673"
	serverAddress = serverHost + ":" + serverPort
	dummyFilter   = Filter{Name: "asdf"}
	validFilter   = Filter{
		Name:     "thing",
		InMemory: true,
		Conn:     &Connection{Server: serverAddress},
	}
	anotherFilter = Filter{
		Name: "another",
		Conn: &Connection{Server: serverAddress},
	}
	infoFields = []string{"storage", "check_hits", "in_memory", "page_outs",
		"page_ins", "size", "check_misses", "capacity", "sets", "checks",
		"set_misses", "set_hits", "probability"}
)

func failIfError(t *testing.T, err error) {
	if err != nil {
		t.Error(err.Error())
	}
}
