# go-bloomd [![Build Status](https://drone.io/github.com/geetarista/go-bloomd/status.png)](https://drone.io/github.com/geetarista/go-bloomd/latest) [![GoDoc](https://godoc.org/github.com/geetarista/go-bloomd/bloomd?status.svg)](https://godoc.org/github.com/geetarista/go-bloomd/bloomd)

A [bloomd](https://github.com/armon/bloomd) client powered by [Go](http://golang.org).

## Installation

```bash
go get github.com/geetarista/go-bloomd/bloomd
```

## Testing

I use Vagrant to run the tests against a BloomD server. Use the included [Vagrantfile](Vagrantfile) and make sure you use your VM's IP address in `test_helpers.go`.

## License

MIT. See `LICENSE`.
