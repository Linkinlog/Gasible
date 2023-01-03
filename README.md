# Gasible
A lightweight tool written in Go to automate aspects of setting up a local development environment.
Written by Log

---
---

Currently, this tool aims to automate the following processes:
- Installing packages
- Using an existing bare Git repo to unload and sync all config files
  - Or creating one
- Configuring certain system services
  - SSH, SSH Config, SSH Keys
  - TeamViewer
  - Hostname
  - Static IP

---
---

## Usage

Gasible can be ran via the command line as such
```bash
gasible
```

It will look for a `gas.yml` file in the directory that the executable is in and then go from there

---
---

## Configuring via gas.yml

Below is the default config featuring all the supported options and some explanation
```YAML
---
# Package Configuration
pkg-manager: "dnf"
packages: # Packages to install
  - "python3-pip"
  - "util-linux-user"
  - "wget"
  - "neovim"
  - "zsh"
  - "docker"
  - "gh"
# Services we will set up
services:
  - teamviewer: true # By default we will install and enable teamviewer
  - ssh: true # By default we will set up ssh to be enabled on boot, we will also create a ssh key and config file
# General Configuration
config:
  - hostname: "development-station" # Set the hostname to be "development-station
  - staticIP: # Sets the static IP, required IP and mas
    - IP: "192.168.4.20" # Must be a valid IP/Subnet mask
    - mask: "255.255.255.0" 
  # TODO make this encrypted or store it better
  - teamViewerCreds: # We can sign into teamviewer and enable easy logon if we have creds
    - user: "username"
    - pass: "password"
```

