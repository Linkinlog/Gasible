# Gasible 
![Gasible logo](https://raw.githubusercontent.com/Linkinlog/Gasible/development/.github/logo.jpeg)

Welcome to **Gasible**, the ultimate solution for automating the setup of your local development environment. Our fast and efficient CLI tool, written in pure Go, makes the process of setting up a new development environment as easy and streamlined as possible. With Gasible, you can have your development environment up and running in minutes, not hours.

## Features
- Installs packages using your chosen package manager (dnf, apt, yum, pacman, etc)
- Sets up a bare git repository for local config management, or uses an existing GitHub repo
- Configures general settings such as hostname, IP, and DNS
- Configures services such as SSH
- Customizable configurations through a config file
- Easy-to-use command-line interface

## Installation
```bash
go install github.com/Linkinlog/gasible
```

## Usage
```bash
gasible [command]
```

## Commands
- `init`: Initializes a new development environment
- `config`: Shows the current configuration
- `update`: Updates packages and configurations

## Configuration
Gasible uses a config file named `config.yml` for customization. You can specify your own package manager, repositories, and configurations in this file. By default, Gasible will look for this file in the home directory, but you can specify a different location by passing the `-c` or `--config` flag.

Below is the default config featuring all the supported options and some explanation
```YAML
---
# Package manager config
pkg-manager-command: dnf
command-args: install -y
packages:
    - python3-pip
    - util-linux-user
    - wget
    - neovim
    - zsh
    - docker
    - gh
# Which processes to run
installer: true
teamviewer: true
ssh: true
git: true
# General config
hostname: development-station
staticIP: 192.168.4.20
mask: 255.255.255.0
TeamViewerCreds:
    user: username
    pass: password

```

## Contribution
We welcome contributions to Gasible. If you find a bug or want to request a new feature, please open an issue. If you want to contribute code, please fork the repository and open a pull request. Our community is always looking for ways to improve and make Gasible even better.

## License
Gasible is licensed under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0).

## Acknowledgements
Gasible was inspired by other [similar](https://github.com/ansible/ansible) projects and we have used their practices as a reference. We are grateful for the contributions of the open-source community, and we hope that Gasible will be a valuable addition to it.
