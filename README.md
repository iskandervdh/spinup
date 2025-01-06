<p align="center"><img src="https://raw.githubusercontent.com/iskandervdh/spinup/refs/heads/main/images/logo.png" width="367" height="150" alt="Spinup Logo"></p>

# Spinup

Quickly spin up your multi command projects.

## Installation

### Pre-requisites

To run spinup you need to have go installed on your system. You can download it from the [official website](https://golang.org/dl/).

### Installation

Download an installer from the releases and run it.

Afterwards you can initialize spinup using the following command:

```bash
spinup init
```

This will create a `spinup` folder in your `.config` folder with the configuration files needed to run spinup.

## Configuration

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

`port` and `domain` are reserved variables that are used to define the port and domain of the project. These are required to be when adding a project.

More information on variables can be found in the [Variables](#variables) section.

### Projects

#### Adding a project

To add a project you can use the following command:

```bash
spinup project add <name> <domain> <port> [commands...]
```

This will create a configuration for the project in the sqlite database for spinup located in your `.config/spinup` folder.

Example:

```bash
spinup project add example example.local 8001 example1 example2
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

## Variables

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

## Usage

To run a project you can use the following command:

```bash
spinup run <project>
```

This will run the commands defined in the configuration for the project.
