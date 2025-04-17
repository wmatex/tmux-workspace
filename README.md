# tmux-workspace

[![Go Report Card](https://goreportcard.com/badge/github.com/wmatex/tmux-workspace)](https://goreportcard.com/report/github.com/wmatex/tmux-workspace)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/wmatex/tmux-workspace)](https://github.com/wmatex/tmux-workspace)

A powerful session manager for Tmux with automatic configuration of each session based on patterns. Easily switch between projects, automatically setup windows and panes, and manage session lifecycle with hooks.

## Table of Contents

- [Motivation](#motivation)
- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [Examples](#examples)
- [Contributing](#contributing)
- [License](#license)

## Motivation

Since switching to my new terminal-with-tmux-only dev workflow, I needed some session manager to easily switch between projects (sessions) and also automatically setup the session with windows and panes. I tried to use [sessionx](https://github.com/omerxx/tmux-sessionx) with [tmuxinator](https://github.com/tmuxinator/tmuxinator), which kind of worked, but I had some issues with it, namely:

- The sessionx popup was really slow to first load for me, which was getting really annoying
- Tmuxinator project doesn't have a workflow to setup realiable commands add the start end end of the session (at least for me and for the way I was using it)
- I have many projects which I switch frequently in and out of and configuring each one was a bit tedious

So I decided to write my own session manager, which would have the following features to satisfy my needs:

### Speed

The interaction must feel _fast_. I should be able to switch active sessions almost instantly and not wait for the UI to load 500ms (sorry sessionx, otherwise you are amazing) before I can do something. This also determined the language as Go

### Native session lifecycle

Tmux has a feature called "hooks", which allow you to react to various triggers e.g. closing window or session, which felt like the exact feature I needed for my workflow. I should be able to define set of commands to run when the session is first started and then a set which will be run _anytime_ the session is closed.

### Project configuration based on patterns

This was probably the biggest reason for me to write new application from scratch. Most of my projects follow some kind of pattern. From PHP projects based on Symfony or full JS projects or even the ones in Go. When I open the project I end up with one main window with neovim and then 1 or 2 panes with the client and/or server running to immediately see the build errors and then an empty pane where I run tests or other commands. This covers 90% of the projects I'm working on. The goal was to have a general configuration which would apply to all these projects and setup the dev workflow in somewhat standardized way. Of course there should be an option to make exceptions for particular projects.

## Features

- **Fast Project Switching**: Quickly switch between projects with a fuzzy finder interface
- **Automatic Session Setup**: Configure windows and panes based on project patterns
- **Session Lifecycle Management**: Define commands to run at session start and end
- **Pattern-Based Configuration**: Apply configurations to projects based on patterns
- **Project-Specific Overrides**: Override default configurations for specific projects

## Requirements

- tmux
- fzf
- Go 1.18+ (for installation from source)

## Installation

### Using Go

```sh
go install github.com/wmatex/tmux-workspace@latest
```

This will download, compile and put the binary to `$GOPATH/bin/tmux-workspace`. You can inspect the value of the variable by running `go env`.

### From Source

```sh
git clone https://github.com/wmatex/tmux-workspace.git
cd tmux-workspace
go build -o tmux-workspace
```

## Usage

The best way to integrate tmux-workspace with tmux is to add a keybinding to your tmux configuration:

```
bind-key W run-shell "$GOPATH/bin/tmux-workspace"
```

Add this line to your `~/.tmux.conf` file and reload the configuration with `tmux source-file ~/.tmux.conf`.

## Configuration

The configuration is managed using a YAML file located in the `$XDG_CONFIG_HOME/tmux-workspace/config.yaml` directory (typically `~/.config/tmux-workspace/config.yaml`).

### Basic Configuration

```yaml
projects:
  # Directories to look for projects
  lookup_dirs:
    - "/path/to/projects"
    - "~/another/path/to/projects"

  # Default window layout
  layout: "main-vertical"
```

### Rules Configuration

You can define rules for different project types in your configuration file:

```yaml
rules:
  - checks:
      dir_exists: .git
    windows:
      editor:
        panes:
          - vim
          - ls

  - checks:
      file_exists: compose-dev.yaml
    hooks:
      start:
        - docker compose -f compose-dev.yaml up -d
      end:
        - docker compose -f compose-dev.yaml down

  - name: turborepo
    checks:
      file_exists: turbo.json
      exec: jq -e '.tasks.dev' turbo.json
    windows:
      editor:
        panes:
          - turbo dev

  - checks:
      file_exists: nest-cli.json
      not_active: turborepo
    windows:
      editor:
        panes:
          - npm run start:debug
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
