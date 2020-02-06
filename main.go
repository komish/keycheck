package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/tidwall/gjson"
	goyamlv2 "gopkg.in/yaml.v2"
	"sigs.k8s.io/yaml"
)

func main() {
	os.Exit(Run())
}

// Target represents a path in a given JSON or YAML document to be checked. These are found
// in the provided specfile. Path is JSON dot notation (this.that.2.theother) as defined in github.com/tidwall/gjson
// which is used for parsing.
type Target struct {
	Path     string `yaml:"path" json:"path"`
	Message  string `yaml:"msg" json:"msg"`
	Required bool   `yaml:"required" json:"required"`
}

const vers string = "0.0.1"

var (
	specFile string // specFile flag
	help     bool   // help flag. prevents pflag default behavior
	version  bool   // version flag
)

func init() {
	pflag.StringVarP(&specFile, "specfile", "s", "", "(REQUIRED) The specification file containing a list of\npaths to check in your json/yaml document. This can\nbe YAML or JSON.")
	pflag.BoolVarP(&help, "help", "h", false, "Prints this help output.")
	pflag.BoolVarP(&version, "version", "v", false, "Prints the version.")
	// override default usage statement
	pflag.Usage = printUsage
}

// Run is the primary entrypoint for keycheck. This function returns an int which
// is passed to os.Exit() as our exit code.
func Run() int {
	pflag.Parse()

	// User asked for help. Stop here if so.
	if help {
		pflag.Usage()
		return 0
	}

	// User asked for version. Stop here if so.
	if version {
		fmt.Printf("%s %s\n", os.Args[0], vers)
		return 0
	}

	// Get the positional args
	args := pflag.Args()

	// User should provide files to parse, otherwise we stop.
	if len(args) < 1 {
		fmt.Printf("ERROR: positional arguments are required\n\n")
		pflag.Usage()
		return 3
	}

	// Get the raw data from the files.
	specFileData, err := getInputFileContent(specFile)
	if err != nil {
		fmt.Printf("ERROR: Unable to get spec file at path %s.\n%s\n", specFile, err)
	}

	// Unmarshal the specfile into our type. Note that we take
	// advantage of the default value of Required (Type: bool). If the
	// user did not define it in our spec file, we assume the path item
	// should be missing in the input file(s).
	var targets []Target
	err = goyamlv2.Unmarshal(specFileData, &targets)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return 4
	}

	// Check for targets in each provided input file.
	for _, v := range args {
		fmt.Printf("============> Results for file: %s\n", v)
		targetData, err := getInputFileContent(v)
		if err != nil {
			// if get this, we don't want to terminate so we'll just skip the file.
			fmt.Printf("ERROR: Unable to get target file at path %s.\n%s\n", v, err)
			continue
		}
		// gjson allows the use of dotnotation so we convert to json here
		targetDataAsJSON, _ := yaml.YAMLToJSON([]byte(targetData))
		// Check for each deprecation
		for _, v := range targets {
			// We have two cases to match so we'll toggle this
			// and print at the end
			printWarning := false

			// Get the result once and parse.
			result := gjson.GetBytes(targetDataAsJSON, v.Path)

			// Default behavior is to print a warning message if we find something
			// and expect it not to be there (i.e. "deprecation")
			if result.Exists() && !v.Required {
				printWarning = true
			}

			// alternatively catch the case if we don't have a value and it's required.
			if !result.Exists() && v.Required {
				printWarning = true
			}

			// Based on what we found, print the message.
			if printWarning {
				fmt.Println()
				fmt.Printf("   Item  %s\nMessage  %s\n", v.Path, v.Message)
			}
		}
		fmt.Println()
	}

	return 0
}
