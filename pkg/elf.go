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
	// ELF64Ehdr.eType related constants

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

const (
	// ELF64Phdr.pType related constants
	PT_NULL      = 0
	PT_LOAD      = 1
	PT_DYNAMIC   = 2
	PT_INTERP    = 3
	PT_NOTE      = 4
	PT_SHLIB     = 5
	PT_PHDR      = 6
	PT_LOPROC    = 0x70000000
	PT_HIPROC    = 0x7fffffff
	PT_LOOS      = 0x60000000
	PT_GNU_STACK = PT_LOOS + 0x474e551
)

const (
	// ELF64Phdr.pFlags related constans
	PF_X = 0x1
	PF_W = 0x2
	PF_R = 0x4
)

const (
	// ELF64Shdr.shType related constans
	SHT_NULL     = 0
	SHT_PROGBITS = 1 // e.g. for .comment, .data, .debug, etc sections
	SHT_SYMTAB   = 2
	SHT_STRTAB   = 3
	SHT_RELA     = 4
	SHT_HASH     = 5
	SHT_DYNAMIC  = 6
	SHT_NOTE     = 7
	SHT_NOBITS   = 8 // e.g. for .bss section
	SHT_REL      = 9
	SHT_SHLIB    = 10
	SHT_DYNSYM   = 11
	SHT_LOUSER   = 0x80000000
	SHT_HIUSER   = 0xffffffff
)

const (
	// ELF64Shdr.shFlags related constants
	SHF_WRITE     = 0x1
	SHF_ALLOC     = 0x2
	SHF_EXECINSTR = 0x3
	SHF_MASKPROC  = 0xf0000000
)

const (
	// ELF64SymbolTable.stInfo related constants
	STT_NOTYPE  = 0
	STT_OBJECT  = 1
	STT_FUNC    = 2
	STT_SECTION = 3
	STT_FILE    = 4
	STT_COMMON  = 5
	STB_LOCAL   = 0
	STB_GLOBAL  = 1
	STB_WEAK    = 2
)

func InitELF() {
	fmt.Println("Hello From ELF!")
}

// ELF64Ehdr ELF header
type ELF64Header struct {
	eIdent     [EI_INDENT]byte
	eType      uint16
	eMachine   uint16
	eVersion   uint32
	eEntry     Elf64_Addr
	ePHOff     Elf64_Off // Program Header Table Offset
	eSHOff     Elf64_Off // Section Header Table Offset
	eFlags     uint32
	eEHSize    uint16 // ELF Header's Size
	ePHEntSize uint16 // Program Header Entry Size
	ePHNum     uint16 // Number of entries in program header table
	eSHEntSize uint16 // Section Header Entry Size
	eSHNum     uint16 // Number of entries in section header table
	eSHStrNdx  uint16 // Section header table index of an entry
}

// ELF64Phdr ELF Program header
type ELF64ProgramHeader struct {
	pType     uint32
	pFlags    uint32
	pOffset   Elf64_Off
	pVAddr    Elf64_Addr
	pPAddr    Elf64_Addr
	pFileSize uint64
	pMemSize  uint64
	pAlign    uint64
}

// ELF64Shdr ELF Section header
type ELF64SectionHeader struct {
	shName      uint32
	shType      uint32
	shFlags     uint32
	shAdrr      Elf64_Addr
	shOffset    Elf64_Off
	shSize      uint64
	shLink      uint32
	shInfo      uint32
	shAddrAlign uint64
	shEntSize   uint64
}

type ELF64SymbolTable struct {
	stName    uint32
	stInfo    byte
	stOther   byte
	stSHIndex uint16
	stValue   Elf64_Addr
	stSize    uint64
}
