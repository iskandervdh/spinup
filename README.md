<p align="center"><img src="https://raw.githubusercontent.com/iskandervdh/spinup/refs/heads/main/images/logo.png" width="367" height="150" alt="Spinup Logo"></p>

# Spinup

Quickly spin up your multi command projects.

## Installation

### Requirements

- Nginx
- Dnsmasq (only for unix based systems)

### Debian based systems

To install the required packages on a debian based system you can use the following command:

```bash
sudo apt install libgtk-3-0 libwebkit2gtk-4.0-dev nginx dnsmasq
```

> [!NOTE]
> For Ubuntu 24.04 and up, you should install `libwebkit2gtk-4.1-dev` instead of `libwebkit2gtk-4.0-dev`.

Download the `.deb` package from the releases. There is a separate version for Ubuntu 24.04 due to the different `libwebkit2gtk` package version being used.

To install the package run the following command where `{{version}}` is the version of the package:

```bash
sudo dpkg -i spinup-{{version}}.deb
```

### RPM based systems

> [!WARNING]
> This has not been tested

To install the required packages on a rpm based system you can use the following command:

```bash
sudo dnf install nginx dnsmasq
```

Download the `.rpm` package from the releases.

To install the package run the following command where `{{version}}` is the version of the package:

```bash
sudo rpm -i spinup-{{version}}.rpm
```

### MacOS

> [!WARNING]
> This has not been tested

Install the required packages:

```bash
brew install nginx dnsmasq
```

Download the `spinup-{{version}}-macos.zip` archive from the releases and unzip the archive:

```bash
unzip spinup-{{version}}-macos.zip
```

Run the installation script:

```bash
sudo ./spinup-{{version}}-macos/install.sh
```

This will create the required directories and files.

## Running the app

To run the app you can use the following command:

```bash
spinup
```

This will start the app where you can add command (templates) and projects.

## CLI

### Commands

#### Adding a command

To add a command template you can use the following command:

```bash
spinup command add <name> <command>
```

**Example:**

```bash
spinup command add example "npm run dev"
```

#### Removing a command

To remove a command template you can use the following command:

```bash
spinup command remove|rm <name>
```

**Example:**

```bash
spinup command remove example
```

#### Listing commands

To list all the command templates you can use the following command:

```bash
spinup command list|ls
```

#### Custom Variables

Commands are templates, so we can use variables that are then defined in the project configuration.

**Example:**

```bash
spinup command add example "npm run dev -- --port {{port}}"
```

`port` and `domain` are reserved variables that are used to define the port and domain (based on the name) of the project. These are required to be when adding a project.

More information on variables can be found in the [Variables](#variables) section.

### Projects

#### Adding a project

To add a project you can use the following command:

```bash
spinup project add <name> <port> [commands...]
```

This will create a configuration for the project in the sqlite database for spinup located in your `.config/spinup` folder.

Example:

```bash
spinup project add example 8001 example1 example2
```

#### Removing a project

To remove a project you can use the following command:

```bash
spinup project remove|rm <name>
```

**Example:**

```bash
spinup project remove example
```

#### Listing projects

To list all the projects you can use the following command:

```bash
spinup project list|ls
```

### Variables

You can add custom variables to the project configuration file. These variables can be used in the command templates.

```bash
spinup variable add <project> <name> <value>
```

**Example:**

If we define a command template like this:

```bash
spinup command add "example" "npx vite -- --loglevel {{loglevel}}"
```

We can add a variable to the project configuration like this:

```bash
spinup variable add example loglevel silent
```

```bash
spinup run example
```

### Running a project

To run a project you can use the following command:

```bash
spinup run <project>
```

This will run the commands defined in the configuration for the project.
