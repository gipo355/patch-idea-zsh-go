# patch-jetbrains-ide

⚠ requires using jetbrains-toolbox

Small CLI utility to patch JetBrains desktop files to use a shell to launch the IDE and inherit the environment variables (paths to runtimes, etc.)

- It will find all the desktop files in the local data directory (usually ~/.local/share/applications)
- list all the JetBrains IDEs installed (e.g. `jetbrains-{name}-{hash}.desktop`)
- patch them to use the shell you choose (sh/bash/zsh) and the path to the shell executable. (replaces the `Exec` line by appending the shell executable to the path)

The shell executable is chosen by the user, and the path to the shell executable is determined by the operating system.

## flags

- `-h` or `--help`: Show help
- `-d` or `--dry-run`: Dry run
- `-a` or `--all-ides`: Select all IDEs
- `-y` or `--all-files`: Select all files
- `-r` or `--repatch`: Repatch all files
- `-c` or `--current-shell`: Use current shell from $SHELL

## example

```bash
patch-jetbrains-ide
```

patch everything without prompting:

```bash
patch-jetbrains-ide -acyr
```

This will patch all the jetbrains desktop files to use the current shell from $SHELL.

## Features

- Find all the desktop files in the local data directory (usually ~/.local/share/applications)
- Patch the desktop files to use the shell you choose (sh/bash/zsh) and the path to the shell executable
- Choose the shell you want to use (sh/bash/zsh)
- Choose the JetBrains IDEs you want to patch (comma-separated numbers, default is all)
- Choose the files you want to patch (default is all)
- Patch the files you choose
- Show the patching results

## Installation

### Install from source

```bash
go install github.com/gipo355/patch-jetbrains-ide@latest
```

## License

MIT
