# Scripts

A collection of personal CLI tools and scripts.

## Usage

### Install (Recommended)
Builds all tools and symlinks them to `~/.local/bin`.

```bash
make install
```

To remove them later:
```bash
make uninstall
```

### Local Build
Builds everything into a local `./bin` directory without modifying your system path.

```bash
make
```

To use these locally, add `./bin` to `PATH`

```bash
export PATH="$(pwd)/bin:$PATH"
```

## Tools

| Tool | Description |
|------|-------------|
| `gcp` | Google Cloud Platform helper (login context, build monitoring, console shortcuts). |
| `modified_files` | Lists files modified in the current git branch (useful for piping to linters). |
