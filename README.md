# scripts
A collection of CLI tools.

## Installation
 1. Use `make` to build and link everything into a new folder `./bin`.
 2. Either add `bin` to your `PATH` manually, or use `source  setup_env.sh`
```

```
todo - installation steps, using setup_env.sh

**Available tools**:
 - gcp:
    - login: switch between gcp SDK configurations.
    - builds
        - list: list recent cloud builds.
        - watch: observe an active cloud build and run a command after it's done.
    - console: open web gcp console with currently set project
 - modified_files: finds all files that have been modified in the current git branch.

