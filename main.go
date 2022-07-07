package main

import (
	"fmt"
	elfreader "github.com/mrata/goelf/pkg"
	"log"
	"os"
)

func main() {
	fmt.Println("Start...")
	inputFile := "myelf"
	fileContent, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Failed to read file: %s", err.Error())
		os.Exit(1)
	}
	elf, err := elfreader.ParseFileContent(fileContent)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)
	}
	fmt.Println(elf.GetFileType())
	elf.PrettyPrint()
	fmt.Println("Done!")
}
