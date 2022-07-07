package main

import (
	"flag"
	elfreader "github.com/mrata/goelf/pkg"
	"log"
	"os"
)

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "file", "", "Input ELF File")

	flag.Parse()
	if inputFile == "" {
		log.Fatalf("No input ELF file is provided.")
	}

	fileContent, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err.Error())
	}
	elf, err := elfreader.ParseFileContent(fileContent)
	if err != nil {
		log.Fatal(err.Error())
	}
	elf.PrettyPrint()

}
