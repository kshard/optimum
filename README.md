<p align="center">
  <img src="./doc/optimum.svg" height="240" />
  <h3 align="center">optimum</h3>
  <p align="center"><strong>data structures management client</strong></p>

  <p align="center">
    <!-- Version -->
    <a href="https://github.com/kshard/optimum/releases">
      <img src="https://img.shields.io/github/v/tag/kshard/optimum?label=version" />
    </a>
    <!-- Documentation -->
    <a href="https://pkg.go.dev/github.com/kshard/optimum">
      <img src="https://pkg.go.dev/badge/github.com/kshard/optimum" />
    </a>
    <!-- Build Status -->
    <a href="https://github.com/kshard/optimum/actions/">
      <img src="https://github.com/kshard/optimum/workflows/build/badge.svg" />
    </a>
    <!-- GitHub -->
    <a href="http://github.com/kshard/optimum">
      <img src="https://img.shields.io/github/last-commit/kshard/optimum.svg" />
    </a>
    <!-- Coverage -->
    <a href="https://coveralls.io/github/kshard/optimum?branch=main">
      <img src="https://coveralls.io/repos/github/kshard/optimum/badge.svg?branch=main" />
    </a>
    <!-- Go Card -->
    <a href="https://goreportcard.com/report/github.com/kshard/optimum">
      <img src="https://goreportcard.com/badge/github.com/kshard/optimum" />
    </a>
  </p>
</p>

--- 

The library is both Golang api and command-line client for managing data
structures.

## What is this about?

> "A data structure is a data organization, and storage format that is usually chosen for efficient access to data" - Wikipedia says.

Data structures are widely utilized across various domains in Computer Science and Software Engineering. Unlike key-value or relational datastores, data structures are an algebraic abstractions that implements a unique properties tailored to meet the specific needs of applications. This library eliminate the extra cost of converting application objects into database entities for each database operation.

This library provides remote access to sophisticated data structures, giving the simplicity of developing application with fewer lines of code to store and access data.

## Getting Started

- [What is this about?](#what-is-this-about)
- [Getting Started](#getting-started)
  - [Getting access](#getting-access)
  - [Typical workflow](#typical-workflow)
    - [List data structures](#list-data-structures)
    - [Create data structure instance](#create-data-structure-instance)
    - [Writing to data structure instance (batch mode)](#writing-to-data-structure-instance-batch-mode)
    - [Reading from data structure instance](#reading-from-data-structure-instance)
    - [Remove data structure instance](#remove-data-structure-instance)
  - [Supported data structures](#supported-data-structures)
- [Using Golang API](#using-golang-api)
  - [Quick Example](#quick-example)
- [How To Contribute](#how-to-contribute)
  - [commit message](#commit-message)
  - [bugs](#bugs)
- [License](#license)


Install the command-line utility from source code. It requires [Golang](https://go.dev) to be installed:

```bash
go install github.com/kshard/optimum/cmd/optimum@latest
```

### Getting access

The library usage requires access to api that provisions and operates data structures for you. Contact your provided for api details.

It is recommended to config environment variables for client usage:

```bash
export HOST=https://example.com
export ROLE=arn:aws:iam::000000000000:role/example-access-role
```


### Typical workflow

Using data structures typically involves a following workflow:
1. List existing data structures.
2. Create a new instance of data structure.
3. Write data.
4. Read data.
5. Remove the data structure instance.

A data structure can be seen as a typed algebraic abstraction that encompasses a collection of data values, the relationships between those values, and the operations or functions that can be applied to manipulate the data. In practical application development, each data structure must be uniquely identifiable to allow efficient access and manipulation. To facilitate this, the application uses a unique reference name called a CURIE (Compact Uniform Resource Identifier). The CURIE combines both the data structure type and a unique identifier, ensuring that the correct data structure is referenced throughout the workflow, enabling smooth interactions within the system.   

See [tutorials](./examples/) for example usage.

#### List data structures

List all data structure instances. It fetches data structure instances of same type. For each provisioned instance it reports NAME, active VERSION, UPDATED AT timestamp, instance STATUS, PENDING version if any, and initialization PARAMS. 

```bash
optimum <type> list -u $HOST

NAME      VERSION          UPDATED AT          | STATUS   PENDING          | PARAMS
example1  NjqOYyOkpMHfg3.6 2024-08-18 10:40:34 | ACTIVE                    | {}
example2                   2024-08-18 10:38:13 | PENDING  NjqOYyOkpMHfg3.6 | {}
```


#### Create data structure instance

Create new instance of data structure. See either documentation of supported
data structure or `optimum help` for details about configuration parameters.

```bash
optimum <type> create -u $HOST -n <name> -j path/to/config.json
```


#### Writing to data structure instance (batch mode)

The batch writing consist of two phases - data upload followed by a commit.
See either documentation of supported data structure or `optimum help` for
details about upload file format.

```bash
# Upload data into server.
optimum <type> upload -u $HOST -n <name> path/to/data.txt

# Commit uploaded data, making it available online.
optimum <type> commit -u $HOST -n <name>
```


#### Reading from data structure instance

Use the REST API for any advanced reading use cases, as the client only supports
basic read operations. See either documentation of supported data structure or
`optimum help` for details about query formats.

```bash
optimum <type> query -u $HOST -n <name> path/to/query.txt
```


#### Remove data structure instance

The command removes data structure instance. The operation is irreversible and
results in the permanent destruction of all data.

```bash
optimum <type> remove -u $HOST -n <name>
```

### Supported data structures

The library supports following data structures:
* `hnsw` [Hierarchical Navigable Small World](./doc/hnsw.md)


Continue with [examples and tutorials](./examples/).

Note: the command line is only support basic operation for data structure manipulation. Use Golang API for any advanced scenario.


## Using Golang API

The latest version of the module is available at `main` branch. All development, including new features and bug fixes, take place on the `main` branch using forking and pull requests as described in contribution guidelines. The stable version is available via Golang modules.

Use `go get` to retrieve the library and add it as dependency to your application.

```bash
go get -u github.com/kshard/optimum
```

### Quick Example

The example below shows usage of client for Hierarchical Navigable Small World.

```go
package main

import (
  "github.com/kshard/optimum"
  "github.com/fogfish/gurl/v2/http"
  "github.com/fogfish/curie"
)

const (
  host = "https://example.com"
  cask = curie.IRI("hnsw:example")
)

func main() {
  // Create client, the library depends on 
  api := optimum.New(http.New(), host)

  // Query the data structure
  neighbors, err := api.Query(context.Background(), cask,
		optimum.Query{Query: []float32{0.1, 0.2, /* ... */ 0.128}},
	)
  
  // Print results
  fmt.Println("Nearest neighbors:", neighbors)
}
```


## How To Contribute

The library is [MIT](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

The build and testing process requires [Go](https://golang.org) version 1.21 or later.


### commit message

The commit message helps us to write a good release note, speed-up review process. The message should address two question what changed and why. The project follows the template defined by chapter [Contributing to a Project](http://git-scm.com/book/ch5-2.html) of Git book.

### bugs

If you experience any issues with the library, please let us know via [GitHub issues](https://github.com/kshard/optimum/issue). We appreciate detailed and accurate reports that help us to identity and replicate the issue. 


## License

[![See LICENSE](https://img.shields.io/github/license/kshard/optimum.svg?style=for-the-badge)](LICENSE)

