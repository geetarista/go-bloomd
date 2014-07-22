package bloomd

import (
	"bufio"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type Connection struct {
	Server   string
	Timeout  time.Duration
	Socket   *net.TCPConn
	File     *os.File
	Attempts int
	Reader   *bufio.Reader
}

// Create a TCP socket for the connection
func (c *Connection) createSocket() (err error) {
	addr, err := net.ResolveTCPAddr("tcp", c.Server)
	if err != nil {
		return err
	}
	c.Socket, err = net.DialTCP("tcp", nil, addr)
	if err != nil {
		return err
	}
	c.Reader = bufio.NewReader(c.Socket)
	if c.Attempts == 0 {
		c.Attempts = 3
	}
	return nil
}

// Sends a command to the server
func (c *Connection) Send(cmd string) error {
	if c.Socket == nil || c.Socket.LocalAddr() == nil {
		err := c.createSocket()
		if err != nil {
			return &BloomdError{ErrorString: err.Error()}
		}
	}
	for i := 0; i < c.Attempts; i++ {
		_, err := c.Socket.Write([]byte(cmd + "\n"))

		if err != nil {
			c.createSocket()
			break
		}
		return nil
	}

	return errSendFailed(cmd, strconv.Itoa(c.Attempts))
}

// Returns a single line from the socket file
func (c *Connection) Read() (line string, err error) {
	if c.Socket == nil || c.Socket.LocalAddr() == nil {
		err := c.createSocket()
		if err != nil {
			return "", &BloomdError{ErrorString: err.Error()}
		}
	}

	l, rerr := c.Reader.ReadString('\n')
	if rerr != nil && rerr != io.EOF {
		return l, &BloomdError{ErrorString: rerr.Error()}
	}
	return strings.TrimRight(l, "\r\n"), nil
}

// Reads a response block from the server. The servers responses are between
// `start` and `end` which can be optionally provided. Returns an array of
// the lines within the block.
func (c *Connection) ReadBlock() (lines []string, err error) {
	first, err := c.Read()
	if err != nil {
		return lines, err
	}
	if first != "START" {
		return lines, &BloomdError{ErrorString: "Did not get block start START! Got '" + string(first) + "'!"}
	}

	for {
		line, err := c.Read()
		if err != nil {
			return lines, err
		}
		if line == "END" || line == "" {
			break
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}

// Convenience wrapper around `send` and `read`. Sends a command,
// and reads the response, performing a retry if necessary.
func (c *Connection) SendAndReceive(cmd string) (string, error) {
	err := c.Send(cmd)
	if err != nil {
		return "", err
	}
	return c.Read()
}

func (c *Connection) responseBlockToMap() (map[string]string, error) {
	lines, err := c.ReadBlock()
	if err != nil {
		return nil, err
	}
	theMap := make(map[string]string)
	for _, line := range lines {
		split := strings.SplitN(line, " ", 2)
		theMap[split[0]] = split[1]
	}
	return theMap, nil
}
