![header-light-only](https://github.com/user-attachments/assets/800c7b87-55d1-40f7-9cb4-2ea6f1a32be5#gh-light-mode-only)
![header-dark-only](https://github.com/user-attachments/assets/9d121acb-de3f-4bfb-9a1c-b9bf6ce8063a#gh-dark-mode-only)
# Axt

Axt turns structured logs, awesome for applications on servers, to something
beautiful on your local dev machine.

This project is mostly a very thin wrapper for the awesome project [pterm](github.com/pterm)

You can easily configure the project to your needs but you must build it.

## How to use

Axt is dead-simple and you should build the binary yourself.

```
git clone https://github.com/ctholho/axt.git
cd axt
go mod tidy
go build -o axt cmd/cli.go
```

Use the `axt` binary and pipe the output of your server application into it.


```bash
go run test/test.go | ./axt
```

Or install globally, by adding it to your bin

```bash
# linux & macOS
sudo mv axt /usr/local/bin
```

## FAQ
<details>
<summary>Do you guarantee backwards compatibility?</summary>
No, because – come on.
</details>

<details>
<summary>What about version numbers?</summary>
No, because – come on.
</details>

