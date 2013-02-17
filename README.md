# go-bloomd ![Build Status](https://travis-ci.org/geetarista/go-bloomd.png)

A [bloomd](https://github.com/armon/bloomd) client powered by [Go](http://golang.org).

## Installation

```bash
go get github.com/geetarista/go-bloomd/bloomd
```

## Documentation

[Read it online](http://godoc.org/github.com/geetarista/go-bloomd)

Or read it locally:

```bash
go doc github.com/geetarista/go-bloomd
```

## Testing

I use Vagrant to run the tests against a BloomD server. Use the included [Vagrantfile](Vagrantfile) and make sure you use your VM's IP address in `test_helpers.go`.

## License

MIT. See `LICENSE`.
