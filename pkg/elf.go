package elf

import "fmt"

const (
	EI_INDENT = 16
)
const (
	// Ehdr.eIdent related constants

	// Ehdr.eIdent array indexes

	EI_MAG0    = 0
	EI_MAG1    = 1
	EI_MAG2    = 2
	EI_MAG3    = 3
	EI_CLASS   = 4
	EI_DATA    = 5
	EI_VERSION = 6
	EI_OSABI   = 7
	EI_PAD     = 8

	// Ehdr.eIdent array values

	ELFMAG0 = 0x7f
	ELFMAG1 = 'E'
	ELFMAG2 = 'L'
	ELFMAG3 = 'F'

	// Ehdr.eIdent[EI_CLASS] values

	ELFCLASSNONE = 0
	ELFCLASS32   = 1
	ELFCLASS64   = 2
	ELFCLASSNUM  = 3

	// Ehdr.eIdent[EI_DATA] values

	ELFDATANONE = 0
	ELFDATA2LSB = 1 // least significant bit, i.e. little endian!
	ELFDATA2MSB = 2 // most significant bit, i.e. big endian!

	// Ehdr.eIdent[EI_VERSION] values

	EV_NONE    = 0
	EV_CURRENT = 1
	EV_NUM     = 2

	// [Some of] Ehdr.eIdent[EI_OSABI] values

	ELFOSABI_NONE  = 0
	ELFOSABI_LINUX = 3
)

const (
	// Ehdr.eType related constants

	ET_NONE = 0 // unknown elf type
	ET_REL  = 1 // relocatable elf file
	ET_EXEC = 2 // executable elf file
	ET_DYN  = 3 // dynamic shared library elf file
	ET_CORE = 4 // core elf file type
)

const (
	X86 = 0x03
	// TODO: Support more machine type
)

func InitELF() {
	fmt.Println("Hello From ELF!")
}

// Ehdr ELF header
type Ehdr struct {
	eIdent     [EI_INDENT]byte
	eType      uint16
	eMachine   uint16
	eVersion   uint32
	eEntry     Elf64_Addr
	ePHOff     Elf64_Off // Program Header Table Offset
	eSHOff     Elf64_Off // Section Header Table Offset
	eFlags     uint32
	eEhSize    uint16
	ePhEntSize uint16
	ePhNum     uint16
	eShentSize uint16
	eShNum     uint16
	eShStrndx  uint16
}
