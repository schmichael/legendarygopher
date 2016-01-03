# Legendary Gopher

Just playing around with Dwarf Fortress's legends export in Go.

## Usage

The easiest way to get started is to download the latest release binary for
your platform here:

https://github.com/schmichael/legendarygopher/releases

Then run it and pass in your legends xml:

```sh
legendarygopher some-legends-dump.xml

# gzipped dumps are also supported
legendarygopher some-legends-dump.xml.gz
```

**Once the xml is parsed open http://localhost:6060/ in a browser.**

*Need an XML file?* [Download a gzipped sample](https://gist.github.com/schmichael/3c8e0a9f0f36cbda9089/raw/32a25692941025f9d2d9688942a45160d5f5a494/region1-00255-10-18-legends.xml.gz)

Turning off the web interface dumps raw text:

```sh
legendarygopher -http="" some-legends-dump.xml
```

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

### Contributing

Pull requests welcome!

Check out https://github.com/schmichael/legendarygopher/issues for the roadmap.
