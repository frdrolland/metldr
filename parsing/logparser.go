package parsing

import (
	"bufio"
	"compress/gzip"
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

	bReader := bufio.NewReader(inputFile)
	testBytes, err := bReader.Peek(2) //read 2 bytes
	if nil != err {
		log.Fatal("FATAL: ", err)
		return err
	}

	var scanner *bufio.Scanner
	if testBytes[0] == 31 && testBytes[1] == 139 {
		// gzipped file
		gzipReader, err := gzip.NewReader(bReader)
		if err != nil {
			log.Fatal("FATAL: ", err)
			return err
		}
		scanner = bufio.NewScanner(gzipReader)
	} else {
		//Not gzipped, just make a scanner based on the reader
		scanner = bufio.NewScanner(bReader)
	}

	//	scanner = bufio.NewScanner(inputFile)
	for scanner.Scan() {
		parse(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
