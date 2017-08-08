package parsing

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// Parse command-line arguments and initialize configuration struct from it.
func ParseLines(filePath string, parse func(string) (string, bool)) ([]string, error) {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	if nil == inputFile {
		log.Fatal("File not found : " + filePath)
		return nil, nil
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	var results []string
	for scanner.Scan() {
		if output, add := parse(scanner.Text()); add {
			fmt.Println(output)
			results = append(results, output)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return results, nil
}
