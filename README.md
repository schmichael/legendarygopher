# Legendary Gopher

Just playing around with Dwarf Fortress's legends export in Go.

## Usage

The easiest way to get started is to download the latest release binary for
your platform here:

https://github.com/schmichael/legendarygopher/releases

If you have Go installed you can use it with:

```sh
go get github.com/schmichael/legendarygopher
```

Regardless of how you install it, simply run it and pass it your legends xml:

```sh
legendarygopher some-legends-dump.xml
```

Or to get a web interface:

```sh
legendarygopher -http=:6060 some-legends-dump.xml
```

Once the xml is parsed open http://localhost:6060/ in a browser.

## WARNING

Everything, including the package/repo name/location is subject to change
without warning. **If you want to use this code your best bet is probably just
to copy and paste it into your own projet.**

*Go nuts!*

## Development

If you change templates you must install go-bindata and run go generate:

```sh
go get -u github.com/jteeuwen/go-bindata/...
go generate
go build
./legendarygopher -http=:6060 some-legends-dump.xml
```
