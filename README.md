# Axt

Axt turns structured logs, awesome for applications on servers, to something
beautiful on your local dev machine.

This project is mostly a very thin wrapper for the awesome project [pterm](github.com/pterm)

You can easily configure the project to your needs.

## How to use

Axt is dead-simple and you should build the binary yourself.

```
git clone axt
cd axt
go mod tidy
go build axt.go
```

Use the `axt` binary and pipe the output of your server application into it.


```bash
go run test/test.go | ./axt
```

Or install globally, by adding it to your bin

```bash
# linux & macOS
mv axt /usr/local/bin
```
