package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
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
	configFile, err := ioutil.ReadFile("../regex.yaml")
	if err != nil {
		panic(err)
	}
	// Parse YAML config file
	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}
    
	// Loop over all files in directory
	dir := "."
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// Skip directories and hidden files
		if info.IsDir() || filepath.Base(path)[0] == '.' {
			return nil
		}
        
		// Check file content against each regex pattern
		for _, pattern := range config.RegularExpressions {
			for _, regex := range pattern.Regexes {
				rex:= strings.Replace(regex.Regex, "\n", "", -1)
				re := regexp.MustCompile(rex)
				examples:= regex.Example
				matches := re.FindAllString(examples, -1)
				if len(matches) > 0 {
					fmt.Printf("Found %d matches in %s for regex '%s' (example: %s)\n", len(matches), path, regex.Name, regex.Example)
				}else{
					fmt.Printf("NO MATCHES FOUND with %s: regex (%s)\n  example (%s)\n", regex.Name, rex, regex.Example)
				}
			}
		}
		return nil
	})
}