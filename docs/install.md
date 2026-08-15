# Install mycelium

Download the latest release binary for your OS/arch into `~/.local/bin`:

```
curl -fsSL https://github.com/robertguss/mycelium/releases/latest/download/mycelium-$(uname -s | tr A-Z a-z)-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/') -o ~/.local/bin/mycelium && chmod +x ~/.local/bin/mycelium
```

Ensure `~/.local/bin` is on your `PATH`, then run `mycelium version`.

This phase ships `linux-amd64` and `darwin-arm64` only. Cutting the GitHub
Release tag and verifying install on a clean VM are human evidence, not the
hermetic merge gate.
