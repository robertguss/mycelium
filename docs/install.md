# Install mycelium

## Release binary

Download the latest release for your OS/arch into `~/.local/bin`:

```
curl -fsSL https://github.com/robertguss/mycelium/releases/latest/download/mycelium-$(uname -s | tr A-Z a-z)-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') -o ~/.local/bin/mycelium && chmod +x ~/.local/bin/mycelium
```

Prebuilt binaries ship for **linux-amd64** and **darwin-arm64** only. Other
platforms should build from source.

Ensure `~/.local/bin` is on your `PATH`, then run `mycelium version`.

## From source

Requires [Go 1.26](https://go.dev/dl/) and `CGO_ENABLED=0`.

```bash
git clone https://github.com/robertguss/mycelium.git
cd mycelium
CGO_ENABLED=0 go build -o mycelium ./cmd/mycelium
```

Put the resulting `mycelium` binary on your `PATH`.

To install into your Go bin directory:

```bash
CGO_ENABLED=0 go install github.com/robertguss/mycelium/cmd/mycelium@latest
```

`go install` of a tagged release stamps that version. A local `go build`
without `-ldflags` prints `0.1.0-dev`.

## Verify

```bash
mycelium version
mycelium --help
```

Then create a first idea. See [`user-guide.md`](user-guide.md).
