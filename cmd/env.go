package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`os\.Getenv\("([^"]*)"\)`)
var envRegex = regexp.MustCompile(`([A-Z_]+)=`)

func copyEnvVariables(absEnvFilePath string, outFile *os.File) error {
	envFile, err := os.Open(absEnvFilePath)
	if err != nil {
		fmt.Println("Error opening .env file for appending:", err)
		return err
	}
	defer envFile.Close()

	scanner := bufio.NewScanner(envFile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line != "" {
			_, err := outFile.WriteString(line + "\n")
			if err != nil {
				fmt.Println("Error appending to output file:", err)
				return err
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading .env file for appending:", err)
		return err
	}

	return nil
}

// Function to read lines from a file into a set (map for quick lookup)
func readEnvFile(envFilePath string) (map[string]struct{}, error) {
	envSet := make(map[string]struct{})

	file, err := os.Open(envFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line != "" {
			matches := envRegex.FindAllStringSubmatch(line, -1)
			for _, match := range matches {
				envSet[match[0]] = struct{}{}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return envSet, nil
}

func main() {
	// Read .env file and store its contents in a set
	envPath := os.Args[1]

	absEnvFilePath, err := filepath.Abs(envPath)
	if err != nil {
		fmt.Println("Error getting absolute path for .env file:", err)
		return
	}

	envSet, err := readEnvFile(absEnvFilePath)
	if err != nil {
		fmt.Println("Error reading .env file:", err)
		return
	}

	// Open the output file
	outFile, err := os.Create("env.md")
	if err != nil {
		fmt.Println("Error creating output file:", err)
		return
	}
	defer outFile.Close()

	err = copyEnvVariables(absEnvFilePath, outFile)
	if err != nil {
		fmt.Println("Error copying env variables:", err)
		return
	}

	// Get the current directory
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Walk through the files in the current directory
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file has a .go extension
		if filepath.Ext(info.Name()) == ".go" {
			file, err := os.Open(path)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return err
			}
			defer file.Close()

			// Read the file line by line
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				matches := re.FindAllStringSubmatch(line, -1)
				for _, match := range matches {
					if len(match) > 1 {
						envVariable := match[1] + "="
						// Check if the envVariable already exists in the .env set
						if _, exists := envSet[envVariable]; !exists {
							_, err := outFile.WriteString(envVariable + "\n")
							if err != nil {
								fmt.Println("Error writing to output file:", err)
								return err
							}
						}
					}
				}
			}

			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading file:", err)
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking path:", err)
		return
	}

	fmt.Println("Finished scanning .go files and writing to env.md")
}
