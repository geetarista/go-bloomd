/*
Provides a client abstraction around the BloomD interface.

Example:
	client := bloomd.Client{Server: "10.0.0.30:8673"}
	filter := bloomd.Filter{Name: "coolfilter"}
	if err := bloomd.CreateFilter(filter); err != nil {
		// handle error
	}
	filters, _ := bloomd.ListFilters()
	fmt.Printf("%+v", filters["coolfilter"])
*/
package bloomd

import (
	"strconv"
	"strings"
)

// If using multiple BloomD servers, it is recommended to use a BloomD Ring
// and only use the proxy as the Server field for your client.
type Client struct {
	Server     string
	Timeout    int
	Conn       *Connection
	ServerInfo string
	InfoTime   int
	HashKeys   bool
}

func NewClient(address string) Client {
	return Client{Server: address, Conn: &Connection{Server: address}}
}

func (c *Client) CreateFilter(f *Filter) error {
	if f.Prob > 0 && f.Capacity < 1 { return errInvalidCapacity }

	cmd := "create " + f.Name
	if f.Capacity > 0 {
		cmd = cmd + " capacity=" + strconv.Itoa(f.Capacity)
	}
	if f.Prob > 0 {
		cmd = cmd + " prob=" + strconv.FormatFloat(f.Prob, 'f', -1, 64)
	}
	if f.InMemory {
		cmd = cmd + " in_memory=1"
	}

	err := c.Conn.Send(cmd)
	if err != nil {
		return err
	}
	resp, err := c.Conn.Read()
	if err != nil {
		return err
	}
	if resp != "Done" && resp != "Exists" {
		return errInvalidResponse(resp)
	}
	f.Conn = c.Conn
	f.HashKeys = c.HashKeys
	return nil
}

func (c *Client) GetFilter(name string) *Filter {
	return &Filter{
		Name: name,
		Conn: c.Conn,
		HashKeys: c.HashKeys,
	}
}

// Lists all the available filters
func (c *Client) ListFilters() (responses map[string]string, err error) {
	err = c.Conn.Send("list")
	if err != nil {
		return
	}

	responses = make(map[string]string)
	resp, err := c.Conn.ReadBlock()
	if err != nil {
		return
	}
	for _, line := range resp {
		split := strings.SplitN(line, " ", 2)
		responses[split[0]] = split[1]
	}
	return responses, nil
}

// Instructs server to flush to disk
func (c *Client) Flush() error {
	err := c.Conn.Send("flush")
	if err != nil {
		return err
	}
	resp, err := c.Conn.Read()
	if err != nil {
		return err
	}
	if resp != "DONE" {
		return err
	}
	return nil
}
