package bloomd

import (
	"testing"
)

var TestHash = "61736466da39a3ee5e6b4b0d3255bfef95601890afd80709"

func TestGetKeyHash(t *testing.T) {
	dummyFilter.HashKeys = true
	k := dummyFilter.getKey("asdf")
	if k != TestHash {
		t.Fail()
	}
}

func TestGetKeyString(t *testing.T) {
	dummyFilter.HashKeys = false
	k := dummyFilter.getKey("asdf")
	if k != "asdf" {
		t.Fail()
	}
}

func TestSet(t *testing.T) {
	ok, err := validFilter.Set("derp")
	failIfError(t, err)
	if ok != true {
		t.Error(ok)
	}
}

func TestBulk(t *testing.T) {
	responses, err := validFilter.Bulk([]string{"herp", "derpina"})
	failIfError(t, err)
	for _, response := range responses {
		if response != true {
			t.Error("Bulk fail")
		}
	}
}

func TestMulti(t *testing.T) {
	responses, err := validFilter.Multi([]string{"derp", "herp", "derpina"})
	failIfError(t, err)
	for _, response := range responses {
		if response != true {
			t.Error("Multi fail")
		}
	}
}

func TestClose(t *testing.T) {
	err := validFilter.Close()
	failIfError(t, err)
	err = anotherFilter.Close()
	failIfError(t, err)
}

func TestClear(t *testing.T) {
	err := anotherFilter.Clear()
	failIfError(t, err)
}

func TestFilterFlush(t *testing.T) {
	err := validFilter.Flush()
	failIfError(t, err)
}

func TestInfo(t *testing.T) {
	info, err := validFilter.Info()
	failIfError(t, err)
	for _, field := range infoFields {
		if info[field] == "" {
			t.Error(field + " not in info")
		}
	}
}
