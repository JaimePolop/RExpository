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
	usage := "Usage: Rex [-d <dir> | -g <github-repo>] -r </path/to/regex.yaml> [-t gh_token] [-c]"

	params := map[string]string{}
	for i := 0; i < len(os.Args)-1; i++ {
		switch os.Args[i] {
		case "-r", "-t":
			params[os.Args[i]] = os.Args[i+1]
			i++
		}
	}

	configPath, hasConfigPath := params["-r"]
	if !hasConfigPath {
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
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

	
	if checkIfFaslePos(os.Args) {
		// If not looking for false positives, filter regexes
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

	switch arg := os.Args[1]; arg {
	case "-h":
		fmt.Println(usage)

	case "-d", "-g":
		if len(os.Args) < 3 {
			fmt.Printf("Usage: Rex %s <value>\n", arg)
			return
		}
		value := os.Args[2]
		if arg == "-d" {
			searchRegexInDir(value, config, "")
		} else {
			searchRegexInRepoGithub(value, config, params["-t"])
		}

	default:
		fmt.Println("The first argument must be '-g' or '-d'", arg)
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

func searchRegexInRepoGithub(repoUrl string, config Config, githubToken string) {
	const chunkSize = 5 * 1024 * 1024 // 5MB in bytes

	// If githubToken is not an empty string, then modify the repoUrl to include it
	if githubToken != "" {
		repoUrl = fmt.Sprintf("https://%s@%s", githubToken, strings.TrimPrefix(repoUrl, "https://"))
	}

	// Create a unique temporary directory
	tempDir, err := ioutil.TempDir(os.TempDir(), "repoClone_")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tempDir)

	// Clone the repository
	repoPath := fmt.Sprintf("%s/repo", tempDir)
	cmd := exec.Command("git", "clone", repoUrl, repoPath)
	if err := cmd.Run(); err != nil {
		panic(err)
	}

	// Name the gitlog-file
	sanitizedUrl := strings.ReplaceAll(repoUrl, "https://github.com/", "")
	sanitizedUrl = strings.ReplaceAll(sanitizedUrl, "/", "-")
	sanitizedUrl = strings.ReplaceAll(sanitizedUrl, ":", "-")
	sanitizedUrl = strings.ReplaceAll(sanitizedUrl, "@", "-")
	logUrl := strings.ReplaceAll(sanitizedUrl, "/", "-")

	// Generate git log and write it directly to a large temporary logFile
	logCmd := exec.Command("git", "-C", repoPath, "log", "-p", "-n", "1000", "--since", "5.years.ago")
	tempLogFile := fmt.Sprintf("%s/gitlog-%s-large.txt", tempDir, logUrl)
	file, err := os.Create(tempLogFile)
	if err != nil {
		panic(err)
	}

	logCmd.Stdout = file
	if err := logCmd.Run(); err != nil {
		file.Close()
		panic(err)
	}
	file.Close()

	// Now, read from the large logFile and split it into 5MB chunks
	inputFile, err := os.Open(tempLogFile)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	reader := bufio.NewReader(inputFile)
	chunkCount := 0

	for {
		chunk := make([]byte, chunkSize)
		bytesRead, err := reader.Read(chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		logDir := fmt.Sprintf("%s/gitlog-%s-part%d.txt", tempDir, logUrl, chunkCount)
		logFile, err := os.Create(logDir)
		if err != nil {
			panic(err)
		}

		_, err = logFile.Write(chunk[:bytesRead])
		if err != nil {
			logFile.Close()
			panic(err)
		}
		logFile.Close()
		chunkCount++
	}

	// Remove the cloned repository and large temp log file as they're no longer needed
	if err := os.RemoveAll(repoPath); err != nil {
		panic(err)
	}
	os.Remove(tempLogFile)

	searchRegexInDir(tempDir, config, logUrl)
}
