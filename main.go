package main

import (
	"fmt"
	elf "github.com/mrata/goelf/pkg"
)

func main() {
	fmt.Println("Start...")

	elf.InitELF()

	fmt.Println("Done!")
}
