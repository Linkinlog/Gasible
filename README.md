# Gasible 
![Gasible logo](https://raw.githubusercontent.com/Linkinlog/Gasible/development/.github/logo.jpeg)

Welcome to **Gasible**, the ultimate solution for automating the setup of your local development environment. Our fast and efficient CLI tool, written in pure Go, makes the process of setting up a new development environment as easy and streamlined as possible. With Gasible, you can have your development environment up and running in minutes, not hours.

## Features
- Customizable configurations through a config file
- Modular add-ons for endless expandability
- Easy-to-use command-line interface
- Package management module for ease of use, using your favorite manager

## Installation
- Through Go:
```bash
$ go install github.com/Linkinlog/gasible@latest
```
- Through cloning and building:
```bash
$ git clone https://github.com/Linkinlog/Gasible
$ cd gasible
$ make
$ ./gasible
```
- Through the binary in our [releases](https://github.com/Linkinlog/Gasible/releases)

## Usage
```bash
gasible [command]
```

## Commands

- `setup`: Runs the setup method on all modules
- `update`: Runs the update method on all modules
- `teardown`: Runs the teardown method on all modules
- `generate`: Generates a new config, overwriting the old

For more detail on what each Module does, please check out our Wiki: (TODO)
## Usage

Gasible uses a config file named `config.yml` for customization.
You can specify your own package manager, packages, and all the modules' configuration in this file.
By default, Gasible will look for this file in `$HOME/.gas/`

Below is the default config featuring all the supported options and some explanation

```YAML
---
# Package manager config
GenericPackageManager:
  enabled: true # dictates if this module gets ran
  settings:
    manager: "apt-get" # which package manager to use for package/dependency management
    packages: ["cowsay", "lolcat"] # (optional) array of packages to install when `Setup()` is ran
GitHub:
  enabled: true # dictates if this module gets ran
  settings:
    token-env-key: "GASIBLE_GH" # (optional) the environment variable containing your Github personal access token
SysCall:
  enabled: true # dictates if this module gets ran, this should always be true
  settings: {}
```

## Contribution
We welcome contributions to Gasible. If you find a bug or want to request a new feature, please open an issue. If you want to contribute code, please fork the repository and open a pull request. Our community is always looking for ways to improve and make Gasible even better.
Also check out the CONTRIBUTING.md for extra info

## License
Gasible is licensed under the [Apache 2.0 License](https://www.apache.org/licenses/LICENSE-2.0).

## Acknowledgements
Gasible was inspired by other [similar](https://github.com/ansible/ansible) projects, and we have used their practices as a reference. We are grateful for the contributions of the open-source community, and we hope that Gasible will be a valuable addition to it.
