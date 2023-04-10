package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"os/exec"
	"encoding/json"

	"gopkg.in/yaml.v2"
)

type Regex struct {
	Name    string `yaml:"name"`
	Regex   string `yaml:"regex"`
	Example string `yaml:"example"`
	FalsePositives string `yaml:"falsePositives"`
}

type Pattern struct {
	Name    string  `yaml:"name"`
	Regexes []Regex `yaml:"regexes"`
}

type Config struct {
	RegularExpressions []Pattern `yaml:"regular_expresions"`
}


type Match struct {
	Regex     string `json:"regex"`
	Match     string `json:"match"`
	File      string `json:"file"`
	RegexName string `json:"regexName"`
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

	// Remove the regexes with falsePositives set to true
	for i, pattern := range config.RegularExpressions {
		filteredRegexes := make([]Regex, 0, len(pattern.Regexes))
		for _, regex := range pattern.Regexes {
			if !strings.EqualFold(regex.FalsePositives, "true") {
				filteredRegexes = append(filteredRegexes, regex)
			}
		}
		config.RegularExpressions[i].Regexes = filteredRegexes
	}
    
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Println("Usage: -d regex-search <dir> | -r <github-repo> | -f <github-repos-file>")
		return
	}

	for i := 0; i < len(args); i++ {
        if args[i] == "-c" {
            
            configFile = args[i+1]
            break
        }
    }

	switch arg := args[0]; arg {
	case "-h":
		fmt.Println("Usage: -d regex-search <dir> | -r <github-repo> | -f <github-repos-file>")
	case "-d":
		if len(args) < 2 {
			fmt.Println("Usage: regex-search -d <dir>")
			return
		}
		dir := args[1]
		searchRegexInDir(dir, config, "")
	case "-r":
		if len(args) < 2 {
			fmt.Println("Usage: regex-search -r <github-repo>")
			return
		}
		repoUrl := args[1]
		searchRegexInRepoGithub(repoUrl, config)
		
	case "-rs":
		if len(args) < 2 {
			fmt.Println("Usage: regex-search -rs <github-repo> <github-repo>")
			return
		}
		if len(args) == 2 {
			fmt.Println("Usage: regex-search -rs <github-repo> <github-repo>. Add more repos or change to -r")
			return
		}
		for i := 1; i < len(os.Args); i++ {
			repoUrl := os.Args[i];
			searchRegexInRepoGithub(repoUrl,config)
		}
	default:
		fmt.Println("Usage: regex-search <dir> | <github-repo> | <github-repos-file>")
	}
}

func searchRegexInDir(dir string, config Config, repoName string){
	
	// Loop over all files in directory
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
	// Skip directories and hidden files
	if info.IsDir() || filepath.Base(path)[0] == '.' {
		return nil
	}

	// Read file content
	content, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file %s: %v\n", path, err)
		return nil
	}

	if err := os.MkdirAll("matches", 0777); err != nil {
		panic(err)
	}

	matches := make([]Match, 0)
	
	// Check file content against each regex pattern
	for _, pattern := range config.RegularExpressions {
		for _, regex := range pattern.Regexes {
			rex := strings.Replace(regex.Regex, "\n", "", -1)
			re := regexp.MustCompile(rex)
			foundMatches := re.FindAllString(string(content), -1)
			for _, foundMatch := range foundMatches {
				if foundMatch != "" {
					match := Match{
						RegexName: regex.Name,
						Regex:     rex,
						Match:     foundMatch,
						File:      strings.ReplaceAll(path,"\\","\\\\"),
					}
					matches = append(matches, match)
				}
        	}
		}
	}
	jsonData, err := json.MarshalIndent(matches, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling matches to JSON:", err)
		return nil
	}

	fmt.Println(string(jsonData))

	//Getting name for matches.json
	_, direc := filepath.Split(dir)

	//Create the json file
	jsonName := direc + "-" + repoName + ".json"
	jsonDir := "matches/" + jsonName
	jsonFile, err := os.Create(jsonDir)
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	if _, err := jsonFile.Write(jsonData); err != nil {
		panic(err)
	}

	return nil
})
}

func searchRegexInRepoGithub(repoUrl string, config Config) {

	// Create the log directory if it doesn't exist
	if err := os.MkdirAll("tmp", 0777); err != nil {
		panic(err)
	}
	// Clone the repository
	if err := os.RemoveAll("repo"); err != nil {
		panic(err)
	}
	cmd := exec.Command("git", "clone", repoUrl, "repo")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	
	// Change to repository directory
	if err := os.Chdir("repo"); err != nil {
		panic(err)
	}
	
	// Generate git log
	logCmd := exec.Command("git", "log", "-p")
	logOutput, err := logCmd.Output()
	if err != nil {
		panic(err)
	}

	// Change back to original directory
	if err := os.Chdir("../"); err != nil {
		panic(err)
	}

	//Name the gitlog-file
	rmGitUrl := strings.Replace(repoUrl, "https://github.com/", "", 1)
	logUrl := strings.Replace(rmGitUrl, "/", "-", 1)
	logDir := "tmp/gitlog-" + logUrl +".txt"
	logFile, err := os.Create(logDir)
	if err != nil {
		panic(err)
	}

	if _, err := logFile.Write(logOutput); err != nil {
		panic(err)
	}

	searchRegexInDir("tmp", config, logUrl);

	// Remove the cloned repository
	if err := os.RemoveAll("repo"); err != nil {
		panic(err)
	}
	if err := os.RemoveAll("tmp"); err != nil {
		panic(err)
	}
}
