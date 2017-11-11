/*
Provides a client abstraction around the BloomD interface.

Example:
	client := bloomd.NewClient("127.0.0.1:8673")

	filters, err := client.ListFilters()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", filters)

	filter := client.GetFilter("coolfilter")
	if err := client.CreateFilter(filter); err != nil {
		panic(err)
	}

	if _, err := filter.Set("milkshake"); err != nil {
		panic(err)
	}

	res, err := filter.Multi([]string{"milkshake", "apfelschorle"})
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", res)
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
	if f.Prob > 0 && f.Capacity < 1 {
		return errInvalidCapacity
	}

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
		Name:     name,
		Conn:     c.Conn,
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
