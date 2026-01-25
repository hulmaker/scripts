# Scripts

A collection of personal CLI tools and scripts.

## Usage

### Install
Installs Go tools to `$GOBIN` (defaults to `~/go/bin`) and symlinks scripts from `bin/` to `~/.local/bin`.

```bash
make install
```

To remove them later:
```bash
make uninstall
```

## Tools

| Tool | Description |
|------|-------------|
| `gcp` | Google Cloud Platform helper (login context, build monitoring, console shortcuts). |
| `modified_files` | Lists files modified in the current git branch (useful for piping to linters). |

--- 

**NOTE:**
Some systems have `gcp` setup as an alias for `git cherry-pick`. We recommend removing it with `unalias gcp`
