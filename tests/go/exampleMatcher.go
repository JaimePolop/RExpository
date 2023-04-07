package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Regex struct {
	Name    string `yaml:"name"`
	Regex   string `yaml:"regex"`
	Example string `yaml:"example"`
}

type Pattern struct {
	Name    string  `yaml:"name"`
	Regexes []Regex `yaml:"regexes"`
}

type Config struct {
	RegularExpressions []Pattern `yaml:"regular_expresions"`
}

func main() {
	// Read YAML config file
	configFile, err := ioutil.ReadFile("../../regex.yaml")
	if err != nil {
		panic(err)
	}
	// Parse YAML config file
	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}
    
	// Check file content against each regex pattern
	var any_error bool = false
	for _, pattern := range config.RegularExpressions {
		for _, regex := range pattern.Regexes {
			rex:= strings.Replace(regex.Regex, "\n", "", -1)
			re := regexp.MustCompile(rex)
			examples:= regex.Example
			matches := re.FindAllString(examples, -1)
			if len(matches) == 0 {
				fmt.Printf("NO MATCHES FOUND with %s: regex (%s)\n  example (%s)\n", regex.Name, rex, regex.Example)
				any_error = true
			}
		}
	}
	if (any_error) {
		os.Exit(1)
	}else {
		fmt.Printf("All examples matched their regular expressions\n")
	}
}