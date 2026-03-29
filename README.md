# scut

`scut` is a small terminal utility for saving and reusing shell commands by working directory. It keeps a local SQLite database of shortcuts, shows the shortcuts for your current folder in a Bubble Tea table UI, and lets you copy a selected command to your clipboard. :contentReference[oaicite:0]{index=0}

## What it does

- Stores shortcuts per working directory
- Shows only the shortcuts that belong to the current directory
- Lets you add a shortcut directly from the command line
- Lets you save your most recent shell command from `.zsh_history` or `.bash_history`
- Opens a terminal UI to browse shortcuts
- Copies the selected shortcut to the clipboard on Enter on macOS using `pbcopy` :contentReference[oaicite:1]{index=1}

## How it works

When `scut` starts, it:

1. Loads configuration from a config file, creating one if needed
2. Opens a local SQLite database
3. Detects the current working directory
4. Loads saved shortcuts from the database
5. Filters them to only the current directory
6. Starts the UI unless a command-line flag is used instead :contentReference[oaicite:2]{index=2}

## Usage

### Installation

Clone the repo, and use the provided scripts to install the app.

```bash
git clone https://github.com/haochend413/scut.git
```

The `build.sh` script will build and run the app, while the `add_to_path.sh` script will add the compiled binary to your system path such that you can run the app using `scat`

### Launch the UI

```bash
scut
```

Add the previous command to shortcut

```bash
scut -l
```

Add your own command

```bash
scut -w "your command"
```
