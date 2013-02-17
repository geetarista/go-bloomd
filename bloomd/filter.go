// Provides an interface to a single Bloomd filter

package bloomd

import (
	"crypto/sha1"
	"fmt"
	"strings"
)

type Filter struct {
	Name     string
	Conn     *Connection
	HashKeys bool
	// Optional
	Capacity int     // The initial capacity of the filter
	Prob     float64 // The inital probability of false positives
	InMemory bool    // If True, specified that the filter should be created
}

// Returns the key we should send to the server
func (f *Filter) getKey(key string) string {
	if f.HashKeys {
		h := sha1.New()
		s := h.Sum([]byte(key))
		return fmt.Sprintf("%x", s)
	}
	return key
}

// Adds a new key to the filter. Returns True/False if the key was added
func (f *Filter) Set(key string) (bool, error) {
	cmd := "s " + f.Name + " " + f.getKey(key)
	resp, err := f.Conn.SendAndReceive(cmd)
	if err != nil {
		return false, err
	}
	if resp == "Yes" || resp == "No" {
		return resp == "Yes", nil
	}
	return false, errInvalidResponse(resp)
}

func (f *Filter) groupCommand(kind string, keys []string) (rs []bool, e error) {
	cmd := kind + " " + f.Name
	for _, key := range keys {
		cmd = cmd + " " + f.getKey(key)
	}
	resp, e := f.Conn.SendAndReceive(cmd)
	if e != nil {
		return rs, &BloomdError{ErrorString: e.Error()}
	}
	if strings.HasPrefix(resp, "Yes") || strings.HasPrefix(resp, "No") {
		split := strings.Split(resp, " ")
		for _, res := range split {
			rs = append(rs, res == "Yes")
		}
	}
	return rs, nil
}

// Performs a bulk set command, adds multiple keys in the filter
func (f *Filter) Bulk(keys []string) (responses []bool, err error) {
	return f.groupCommand("b", keys)
}

// Performs a multi command, checks for multiple keys in the filter
func (f *Filter) Multi(keys []string) (responses []bool, err error) {
	return f.groupCommand("m", keys)
}

func (f *Filter) sendCommand(cmd string) error {
	resp, err := f.Conn.SendAndReceive(cmd + " " + f.Name)
	if err != nil {
		return err
	}
	if resp != "Done" {
		return errInvalidResponse(resp)
	}
	return nil
}

// Deletes the filter permanently from the server
func (f *Filter) Drop() error {
	return f.sendCommand("drop")
}

// Closes the filter on the server
func (f *Filter) Close() error {
	return f.sendCommand("close")
}

// Clears the filter on the server
func (f *Filter) Clear() error {
	return f.sendCommand("clear")
}

// Forces the filter to flush to disk
func (f *Filter) Flush() error {
	return f.sendCommand("flush")
}

// Returns the info dictionary about the filter
func (f *Filter) Info() (map[string]string, error) {
	if err := f.Conn.Send("info " + f.Name); err != nil {
		return nil, err
	}
	info, err := f.Conn.responseBlockToMap()
	if err != nil {
		return nil, err
	}
	return info, nil
}
