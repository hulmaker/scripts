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
| `transcribe` | Record from the microphone and transcribe locally with Whisper. |
| `modified_files` | Lists files modified in the current git branch (useful for piping to linters). |

### Requirements

#### `transcribe`

- [uv](https://docs.astral.sh/uv/getting-started/installation/) - package manager
- [ffmpeg](https://ffmpeg.org/) – records audio from the microphone (macOS `avfoundation`).
- [mlx-whisper](https://pypi.org/project/mlx-whisper/) – local speech-to-text optimised for macOS

```bash
brew install ffmpeg openai-whisper
```

Usage: press Enter or Space to stop recording and transcribe, Ctrl+C to cancel.
Recordings and transcripts are saved to `~/Recordings`.

--- 

**NOTE:**
Some systems have `gcp` setup as an alias for `git cherry-pick`. We recommend removing it with `unalias gcp`
