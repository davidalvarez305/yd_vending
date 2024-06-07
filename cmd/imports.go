package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	// Regular expression to match "https://github.com/davidalvarez305/yd_vending/"
	githubRegex := regexp.MustCompile(`github\.com/davidalvarez305/budgeting/`)

	// Walk through the current directory and its subdirectories
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if the file is a regular file and has a .go extension
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			// Open the file
			file, err := os.Open(path)
			if err != nil {
				fmt.Println("Error opening file:", err)
				return err
			}
			defer file.Close()

			// Create a temporary file
			tmpFile, err := os.Create(path + ".tmp")
			if err != nil {
				fmt.Println("Error creating temporary file:", err)
				return err
			}
			defer tmpFile.Close()

			// Scanner to read file line by line
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()

				// Replace the substring if it exists in the line
				newLine := githubRegex.ReplaceAllString(line, "github.com/davidalvarez305/yd_vending/")

				// Write the modified line to the temporary file
				_, err := tmpFile.WriteString(newLine + "\n")
				if err != nil {
					fmt.Println("Error writing to temporary file:", err)
					return err
				}
			}

			if err := scanner.Err(); err != nil {
				fmt.Println("Error reading file:", err)
				return err
			}

			// Close the original file
			file.Close()

			// Close the temporary file
			tmpFile.Close()

			// Replace the original file with the temporary file
			err = os.Rename(path+".tmp", path)
			if err != nil {
				fmt.Println("Error replacing original file:", err)
				return err
			}
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error walking through directory:", err)
		return
	}

	fmt.Println("Replacement complete!")
}
