package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"
)

type Regex struct {
	Name           string `yaml:"name"`
	Regex          string `yaml:"regex"`
	Example        string `yaml:"example"`
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
	configPath := ""
	usage := "Usage: -r </path/to/regex.yaml> [-d regex-search <dir> | -g <github-repo> | -gs <github-repo> <github-repo>] [-c]"

	args := os.Args[1:]

	// Check for -r parameter for the regex.yaml path
	for i, arg := range args {
		if arg == "-r" && i+1 < len(args) {
			configPath = args[i+1]
			// Remove the -r and its value from args
			args = append(args[:i], args[i+2:]...)
			break
		}
	}

	// If the configPath hasn't changed from the default, it means -r wasn't provided.
	if configPath == "" {
		fmt.Println("You must provide the '-r' parameter followed by the path to regex.yaml.")
		fmt.Println(usage)
		return
	}

	// Read YAML config file using the determined path
	configFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}

	var config Config
	var isFaslePos bool

	if len(args) == 0 {
		fmt.Println(usage)
		return
	}

	isFaslePos = checkIfFaslePos(args)
	if isFaslePos {
		// Parse YAML config file
		err = yaml.Unmarshal(configFile, &config)
		if err != nil {
			panic(err)
		}

	} else {
		// Parse YAML config file
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
	}
	switch arg := args[0]; arg {
	case "-h":
		fmt.Println(usage)
	case "-d":
		if isFaslePos {
			if len(args) < 3 {
				fmt.Println("Usage: regex-search -d <dir> -c")
				return
			}
		} else {
			if len(args) < 2 {
				fmt.Println("Usage: regex-search -d <dir>")
				return
			}
		}
		dir := args[1]
		searchRegexInDir(dir, config, "")
	case "-g":
		if isFaslePos {
			if len(args) < 3 {
				fmt.Println("Usage: regex-search -g <github-repo> -c")
				return
			}
		} else {
			if len(args) < 2 {
				fmt.Println("Usage: regex-search -g <github-repo>")
				return
			}
		}
		repoUrl := args[1]
		searchRegexInRepoGithub(repoUrl, config)

	case "-gs":
		if isFaslePos {
			if len(args) < 3 {
				fmt.Println("Usage: regex-search -gs <github-repo> <github-repo> -c")
				return
			}
			if len(args) == 3 {
				fmt.Println("Usage: regex-search -gs <github-repo> <github-repo> -c. Add more repos or use -g")
				return
			}
			for i := 2; i < len(os.Args)-1; i++ {
				repoUrl := os.Args[i]
				searchRegexInRepoGithub(repoUrl, config)
			}
		} else {
			if len(args) < 2 {
				fmt.Println("Usage: regex-search -gs <github-repo> <github-repo>")
				return
			}
			if len(args) == 2 {
				fmt.Println("Usage: regex-search -gs <github-repo> <github-repo>. Add more repos or use -r")
				return
			}
			for i := 2; i < len(os.Args); i++ {
				repoUrl := os.Args[i]
				searchRegexInRepoGithub(repoUrl, config)
			}
		}
	default:
		fmt.Println(usage)
	}
}

func checkIfFaslePos(args []string) bool {
	var isFaslePos bool = false
	for i := 0; i < len(args); i++ {
		if args[i] == "-c" {
			isFaslePos := true
			return isFaslePos
		}
	}
	return isFaslePos
}

func searchRegexInDir(dir string, config Config, repoName string) {

	// Loop over all files in the directory
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

		// Check file content against each regex pattern
		for _, pattern := range config.RegularExpressions {
			for _, regex := range pattern.Regexes {
				rex := strings.Replace(regex.Regex, "\n", "", -1)
				re := regexp.MustCompile(rex)
				foundMatches := re.FindAllString(string(content), -1)
				for _, foundMatch := range foundMatches {
					if foundMatch != "" {
						// Truncate match if it's longer than 500 chars
						if len(foundMatch) < 500 {
							match := Match{
								RegexName: regex.Name,
								Regex:     rex,
								Match:     foundMatch,
								File:      filepath.Base(path),
							}
							jsonData, err := json.Marshal(match)
							if err != nil {
								fmt.Println("Error marshaling match to JSON:", err)
								continue
							}

							fmt.Println(string(jsonData))
						}
					}
				}
			}
		}

		return nil
	})
}

func searchRegexInRepoGithub(repoUrl string, config Config) {
	// Download the github repo and split the github log file in chunks of 5MB
	// Then call the filesystem analysis

	const chunkSize = 5 * 1024 * 1024 // 5MB in bytes

	// Create a unique temporary directory
	tempDir, err := ioutil.TempDir(os.TempDir(), "repoClone_")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir) // Ensure removal of temp directory upon function exit

	// Clone the repository
	repoPath := fmt.Sprintf("%s/repo", tempDir)
	cmd := exec.Command("git", "clone", repoUrl, repoPath)
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// Name the gitlog-file
	rmGitUrl := strings.Replace(repoUrl, "https://github.com/", "", 1)
	logUrl := strings.Replace(rmGitUrl, "/", "-", 1)

	// Generate git log and write it directly to logFile
	logCmd := exec.Command("git", "-C", repoPath, "log", "-p")
	output, err := logCmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := logCmd.Start(); err != nil {
		panic(err)
	}

	reader := bufio.NewReader(output)
	chunkCount := 0

	for {
		chunk := make([]byte, chunkSize)
		_, err := io.ReadFull(reader, chunk)
		if err == io.EOF {
			break
		}
		if err != nil && err != io.ErrUnexpectedEOF {
			panic(err)
		}

		logDir := fmt.Sprintf("%s/gitlog-%s-part%d.txt", tempDir, logUrl, chunkCount)
		logFile, err := os.Create(logDir)
		if err != nil {
			panic(err)
		}

		_, err = logFile.Write(chunk)
		if err != nil {
			logFile.Close()
			panic(err)
		}
		logFile.Close()
		chunkCount++
	}

	if err := logCmd.Wait(); err != nil {
		panic(err)
	}

	// Remove the cloned repository as it's no longer needed
	if err := os.RemoveAll(repoPath); err != nil {
		panic(err)
	}

	searchRegexInDir(tempDir, config, logUrl)
}
