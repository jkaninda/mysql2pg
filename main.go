package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

var inputFile string
var outputFile string

func init() {
	flag.StringVar(&inputFile, "file", "", "Input file")
	flag.StringVar(&outputFile, "output", "", "Output file")
	flag.Parse()

}
func main() {

	if inputFile == "" || outputFile == "" {
		flag.Usage()
		os.Exit(1)
	}
	// Open MySQL backup file for reading backup.sql
	inputFile, err := os.Open(inputFile)
	if err != nil {
		fmt.Println("Error opening MySQL file: ", err)
		return
	}
	defer inputFile.Close()

	// Create output PostgreSQL file
	outputFile, err := os.Create(outputFile)
	if err != nil {
		fmt.Println("Error creating PostgreSQL file:", err)
		return
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)

	// Regex patterns for common MySQL to PostgreSQL changes
	autoIncrementRegex := regexp.MustCompile(`AUTO_INCREMENT=\d+`)
	engineRegex := regexp.MustCompile(`ENGINE=\w+`)
	backtickRegex := regexp.MustCompile("`")
	datetimeRegex := regexp.MustCompile(`DATETIME`)

	for scanner.Scan() {
		line := scanner.Text()

		// Convert MySQL AUTO_INCREMENT to PostgreSQL SERIAL
		line = autoIncrementRegex.ReplaceAllString(line, "")
		line = strings.ReplaceAll(line, "AUTO_INCREMENT", "SERIAL")

		// Remove MySQL-specific ENGINE clause
		line = engineRegex.ReplaceAllString(line, "")

		// Replace MySQL backticks with PostgreSQL double quotes
		line = backtickRegex.ReplaceAllString(line, `"`)

		// Adjust DATETIME to TIMESTAMP if needed
		line = datetimeRegex.ReplaceAllString(line, "TIMESTAMP")

		//TODO: Additional conversions

		// Write the modified line to the PostgreSQL file
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			fmt.Println("Error writing to PostgreSQL file:", err)
			return
		}
	}

	// Flush the buffered writer
	if err := writer.Flush(); err != nil {
		fmt.Println("Error flushing buffer:", err)
		return
	}

	fmt.Println("Conversion complete. The PostgreSQL .sql file is ready.")
}
