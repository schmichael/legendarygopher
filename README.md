# Legendary Gopher

Simple web viewer for Dwarf Fortress 42 Legends XML files.

## Usage

The easiest way to get started is to download the latest release binary for
your platform here:

https://github.com/schmichael/legendarygopher/releases

If you're on Windows you'll have to add a `.exe` extension. See issue #15.

Then run it and pass in your legends xml:

```sh
legendarygopher some-legends-dump.xml
```

Once the xml is parsed open http://localhost:6565/ in a browser.

## Features

* gzipped (`.xml.gz`) and bzipped (`.xml.bz2`) files
* Code Page 437 encoding handling
* JSON HTTP API (for example `/api/world`)
* JSON support (save `/api/world` and pass it in instead of xml)
* Text dump mode (with `-http=""`)

## Development

If you change templates you must install go-bindata and run go generate:

```sh
go get -u github.com/jteeuwen/go-bindata/...
go generate
go build
./legendarygopher -http=:6565 some-legends-dump.xml
```

### Contributing

Pull requests welcome!

Check out https://github.com/schmichael/legendarygopher/issues for the roadmap.

### WARNING

Everything, including the package/repo name/location is subject to change
without warning. **If you want to use this code your best bet is probably just
to copy and paste it into your own projet.**

*Go nuts!*
