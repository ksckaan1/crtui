```
                       ╭─╮
╭────╮╭────╮╭─╮  ╭─╮╭─╮├─┤
│ ╭──╯│ ╭─┬┴╯ ╰─╮│ ││ ││ │
│ │   │ │ ╰─╮ ╭─╯│ ││ ││ │
│ ╰──╮│ │   │ ╰─╮│ ╰╯ ││ │
╰────╯╰─╯   ╰───╯╰────╯╰─╯
```

A terminal-based TUI tool for managing container registries (`registry:2`). Supports both Docker and Podman credentials.

## Features

### 📦 Registry Management
- **List Registries**: View all configured container registries with their connection status (online/offline/unauthorized)
- **Add Registry**: Create new registry connections with custom URL, username, and password credentials
- **Edit Registry**: Modify existing registry connection details
- **Delete Registry**: Remove registry connections from the configuration
- **Auto-detect Credentials**: Automatically discover and use Docker/Podman credentials from config files

### 📂 Repository & Tag Management
- **Browse Repositories**: Navigate through all repositories in a selected registry
- **Search & Filter**: Instantly filter repositories and tags using the search input
- **View Tags**: Display all tags available in a selected repository
- **Delete Repository**: Remove all tags from a specific repository
- **Delete Tags**: Select and delete individual or multiple tags from a repository

### 🔍 Tag Details & Inspection
- **Multi-Platform Support**: View tags built for multiple architectures and operating systems (linux/amd64, linux/arm64, linux/arm/v7, windows/amd64, etc.)
- **Total Size**: See the complete disk size of the image including all layers
- **Environment Variables**: Inspect all defined environment variables passed to the container
- **Entrypoint**: View the container entrypoint command
- **CMD**: View the default command executed when the container starts
- **Working Directory**: See the configured working directory inside the container
- **User**: View the user/group that runs the container process
- **Layers**: Browse through all container image layers with their individual sizes
- **History**: View the complete build history including creation dates and author information
- **RootFS**: Inspect the root filesystem layers and their digests
- **Labels**: View OCI image labels and annotations

### ⌨️ User Experience
- **Copy Pull Command**: Copy the exact `docker pull` command to clipboard with a single keypress
- **Refresh Data**: Manually refresh registry, repository, or tag lists at any time
- **Status Feedback**: Real-time status messages showing success/error states and operation duration
- **Loading Indicators**: Visual spinners and progress indicators during data fetching operations

## Quick Start

```bash
crtui
```

The application will automatically detect any existing Docker or Podman credentials and display them in the registry list.

[![asciicast](https://asciinema.org/a/ujo5faGGk8PPsBeM.svg)](https://asciinema.org/a/ujo5faGGk8PPsBeM)

## Installation

### Homebrew (macOS/Linux)
```bash
brew install ksckaan1/tap/crtui
```

### Debian (Ubuntu, Linux Mint...)
1. Add repository
	```bash
	curl -1sLf \
	  'https://dl.cloudsmith.io/public/ksckaan1/crtui/setup.deb.sh' \
	  | sudo -E bash
	```

2. Install package
	```bash
	sudo apt update
	sudo apt install crtui
	```

### Fedora/RHEL
1. Add repository
	```bash
	curl -1sLf \
	  'https://dl.cloudsmith.io/public/ksckaan1/crtui/setup.rpm.sh' \
	  | sudo -E bash
	```

2. Install package
	```bash
	sudo dnf update
	sudo dnf install crtui
	```

### Alpine (apk)
1. Add repository
	```bash
	sudo apk add --no-cache bash
	curl -1sLf \
	  'https://dl.cloudsmith.io/public/ksckaan1/crtui/setup.alpine.sh' \
	  | sudo -E bash
	```

2. Install package
```bash
sudo apk add crtui
```

### AUR (Arch Linux, CachyOS, Manjaro)
```bash
yay -S crtui-bin
```

### Go Install
```bash
go install github.com/ksckaan1/crtui/cmd/crtui@latest
```

Make sure `$HOME/go/bin` (or `$GOPATH/bin`) is in your PATH.

### Direct Download Prebuilt Binary
Download from [Releases](https://github.com/ksckaan1/crtui/releases):
- **Linux**: `.tar.gz` (amd64, arm64)
- **macOS**: `.tar.gz` (amd64, arm64)
- **Windows**: `.zip` (amd64, arm64)

Extract and run:
```bash
tar -xzf crtui_*.tar.gz
./crtui
```

## Build From Source

### Prerequisites
- Go 1.26+

### Build

```bash
git clone https://github.com/ksckaan1/crtui.git
cd crtui
go build -o crtui ./cmd/crtui
```

## License

MIT License
