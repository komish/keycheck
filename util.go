package main

import (
	"fmt"
	"github.com/spf13/pflag"
	"io/ioutil"
	"os"
)

// getInputFileContent returns a byte slice of a file at the specified path.
func getInputFileContent(path string) ([]byte, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// printUsage overides the default usage statement that pflag provides.
func printUsage() {
	fmt.Println(`Check YAML or JSON files for specific paths, and print a message if 
those paths are present/missing. Useful for checking for missing 
keys or targets in your YAML and JSON documents.
	`)
	fmt.Printf("Usage:\n %s -s specfile /path/to/yaml-or-json ...\n\n", os.Args[0])
	fmt.Println("Flags:")
	pflag.PrintDefaults()
}
