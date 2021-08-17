## go-config-parser [![Go Reference][1]][2]

`go-config-parser` package provides a simple parser for reading/writing `.ini` configuration files.

## Features
This package supports reading/writing the `.ini` files as well as:
- Reading/writing duplicate section names (ex: MySQL NDB engine's `config.ini`)
- Options without values (ex: can be used to group a set of hostnames)
- Options without a named section (ex: a simple `option=value` file)
- Find sections with regexp pattern matching on section names, ex: `dc1.east.webservers` where regexp is `.webservers`
- `#` or `;` as comment delimiter
- `=` or `:` as value delimiter

## Example
```go
package example

import (
	"fmt"
	"log"
	"testing"

	"github.com/teocci/go-config-parser/core/parser"
)

const (
	configFileName      = "config.ini"
	hostnameOption      = "HostName"
	sectionNameWB       = "dc1.webservers"
	sectionNameMYSQLD   = "MYSQLD DEFAULT"
	sectionNameToDelete = "NDB_MGMD DEFAULT"
	sectionNameRegexp   = ".webservers$"
)

// TestExample Example Read and modify a configuration file
func TestExample(t *testing.T) {
	parser.SetDelimiter(parser.EqualDelimiter) // default delimiter

	config, err := parser.Read(configFileName)
	if err != nil {
		log.Fatal(err)
	}
	// Print the full configuration
	fmt.Println(config)

	// get a section
	section, err := config.Section(sectionNameMYSQLD)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("TotalSendBufferMemory=%s\n", section.ValueOf("TotalSendBufferMemory"))

		// set new value
		var oldValue = section.SetValueFor("TotalSendBufferMemory", "256M")
		fmt.Printf("TotalSendBufferMemory=%s, old value=%s\n", section.ValueOf("TotalSendBufferMemory"), oldValue)

		// delete option
		oldValue = section.Delete("DefaultOperationRedoProblemAction")
		fmt.Println("Deleted DefaultOperationRedoProblemAction: " + oldValue)

		// add new options
		section.Add("innodb_buffer_pool_size", "64G")
		section.Add("innodb_buffer_pool_instances", "8")
	}

	// add a new section and options
	section = config.NewSection("NDBD MGM")
	section.Add("NodeId", "2")
	section.Add("HostName", "10.10.10.10")
	section.Add("PortNumber", "1186")
	section.Add("ArbitrationRank", "1")

	// find all sections ending with .webservers
	sections, err := config.Find(sectionNameRegexp)
	if err != nil {
		log.Fatal(err)
	}
	for _, section := range sections {
		fmt.Print(section)
	}
	// or
	config.PrintSection(sectionNameWB)

	sections, err = config.Delete(sectionNameToDelete)
	if err != nil {
		log.Fatal(err)
	}
	// deleted sections
	for _, section := range sections {
		fmt.Print(section)
	}

	options := section.Options()
	fmt.Println(options[hostnameOption])

	// save the new config. the original will be renamed to /etc/config.ini
	err = parser.Save(config, configFileName)
	if err != nil {
		log.Fatal(err)
	}
}
```

[1]: https://pkg.go.dev/badge/github.com/teocci/go-config-parser.svg
[2]: https://pkg.go.dev/github.com/teocci/go-config-parser
