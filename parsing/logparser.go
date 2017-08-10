package parsing

import (
	"bufio"
	"log"
	"os"
)

// Parse command-line arguments and initialize configuration struct from it.
func ParseLines(filePath string, parse func(string) (string, bool)) error {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	if nil == inputFile {
		log.Fatal("File not found : " + filePath)
		return nil
	}
	defer inputFile.Close()

	scanner := bufio.NewScanner(inputFile)
	for scanner.Scan() {
		parse(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
