package elfreader

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

const (
	EI_INDENT = 16
)
const (
	// Ehdr.Ident related constants

	// Ehdr.Ident array indexes

	EI_MAG0    = 0
	EI_MAG1    = 1
	EI_MAG2    = 2
	EI_MAG3    = 3
	EI_CLASS   = 4
	EI_DATA    = 5
	EI_VERSION = 6
	EI_OSABI   = 7
	EI_PAD     = 8

	// Ehdr.Ident array values

	ELFMAG0 = 0x7f
	ELFMAG1 = 'E'
	ELFMAG2 = 'L'
	ELFMAG3 = 'F'

	// Ehdr.Ident[EI_CLASS] values

	ELFCLASSNONE = 0
	ELFCLASS32   = 1
	ELFCLASS64   = 2
	ELFCLASSNUM  = 3

	// Ehdr.eIdent[EI_DATA] values

	ELFDATANONE = 0
	ELFDATA2LSB = 1 // least significant bit, i.e. little endian!
	ELFDATA2MSB = 2 // most significant bit, i.e. big endian!

	// Ehdr.Ident[EI_VERSION] values

	EV_NONE    = 0
	EV_CURRENT = 1
	EV_NUM     = 2

	// [Some of] Ehdr.Ident[EI_OSABI] values

	ELFOSABI_NONE  = 0
	ELFOSABI_LINUX = 3
)

const (
	// ELF64Ehdr.Type related constants

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
	// ELF64Phdr.Type related constants
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

// ELF64Ehdr ELF header
type ELF64Header struct {
	Ident     [EI_INDENT]byte
	Type      uint16
	Machine   uint16
	Version   uint32
	Entry     Elf64_Addr
	PHOff     Elf64_Off // Program Header Table Offset
	SHOff     Elf64_Off // Section Header Table Offset
	Flags     uint32
	EHSize    uint16 // ELF Header's Size
	PHEntSize uint16 // Program Header Entry Size
	PHNum     uint16 // Number of entries in program header table
	SHEntSize uint16 // Section Header Entry Size
	SHNum     uint16 // Number of entries in section header table
	SHStrNdx  uint16 // Section header table index of an entry
}

// ELF64Phdr ELF Program header
type ELF64ProgramHeader struct {
	Type     uint32     `json:"type"`
	Flags    uint32     `json:"flags"`
	Offset   Elf64_Off  `json:"offset"`
	VAddr    Elf64_Addr `json:"VAddr"`
	PAddr    Elf64_Addr `json:"PAddr"`
	FileSize uint64     `json:"fileSize"`
	MemSize  uint64     `json:"memSize"`
	Align    uint64     `json:"align"`
}

// ELF64Shdr ELF Section header
type ELF64SectionHeader struct {
	Name      uint32     `json:"name"`
	Type      uint32     `json:"type"`
	Flags     uint64     `json:"flags"`
	Addr      Elf64_Addr `json:"addr"`
	Offset    Elf64_Off  `json:"offset"`
	Size      uint64     `json:"size"`
	Link      uint32     `json:"link"`
	Info      uint32     `json:"info"`
	AddrAlign uint64     `json:"addrAlign"`
	EntSize   uint64     `json:"entSize"`
}

type ELF64SymbolTable struct {
	Name    uint32     `json:"name"`
	Info    byte       `json:"info"`
	Other   byte       `json:"other"`
	SHIndex uint16     `json:"SHIndex"`
	Value   Elf64_Addr `json:"value"`
	Size    uint64     `json:"size"`
}

type ELF64File struct {
	Header    ELF64Header
	Segments  []ELF64ProgramHeader
	Sections  []ELF64SectionHeader
	byteOrder binary.ByteOrder
	Data      []byte
}

func ParseFileContent(rawBytes []byte) (*ELF64File, error) {

	// File has to have at least the ident.
	if len(rawBytes) < EI_INDENT {
		os.Exit(1)
	}
	var elf ELF64File
	elf.Data = rawBytes

	var header ELF64Header
	data := bytes.NewReader(elf.Data)

	// Read Ident for checking the format and byte order
	var Ident = make([]byte, EI_INDENT)
	err := binary.Read(data, binary.LittleEndian, Ident)
	if err != nil {
		log.Fatalf("Failed to read Ident: %s", err.Error())
		return nil, err
	}

	// Check byte order
	if Ident[EI_DATA] == ELFDATA2LSB {
		elf.byteOrder = binary.LittleEndian
	} else if Ident[EI_DATA] == ELFDATA2MSB {
		elf.byteOrder = binary.BigEndian
	} else {
		log.Fatal("Unknown byte order!")
		return nil, err
	}

	if Ident[EI_CLASS] != ELFCLASS64 {
		log.Fatal("ELF class is not 64 bit!")
		return nil, err
	}

	// read the entire header with the known byte order.
	data = bytes.NewReader(elf.Data)
	err = binary.Read(data, binary.LittleEndian, &header)
	if err != nil {
		log.Fatal("Failed to read the header.")
		return nil, err
	}

	elf.Header = header
	fmt.Println(string(header.Ident[EI_MAG1]), string(header.Ident[EI_MAG2]), string(header.Ident[EI_MAG3]))
	// Parse Program Header

	if elf.Header.PHOff > Elf64_Off(len(elf.Data)) {
		log.Fatalf("Failed to parse program header: Invalid program header offset: %d", elf.Header.PHOff)
		return nil, err
	}
	data = bytes.NewReader(elf.Data[elf.Header.PHOff:])
	segments := make([]ELF64ProgramHeader, elf.Header.PHNum)
	err = binary.Read(data, elf.byteOrder, segments)
	if err != nil {
		log.Fatal("Failed to read segments!")
		return nil, err
	}
	elf.Segments = segments

	// Parse Section Header

	if elf.Header.SHOff > Elf64_Off(len(elf.Data)) {
		log.Fatal("Failed to parse section header: Invalid section header offset.")
		return nil, err
	}
	data = bytes.NewReader(elf.Data[elf.Header.SHOff:])
	sections := make([]ELF64SectionHeader, elf.Header.SHNum)
	err = binary.Read(data, elf.byteOrder, sections)
	if err != nil {
		log.Fatal("Failed to read sections.")
		return nil, err
	}
	elf.Sections = sections

	return &elf, nil
}

func (elf *ELF64File) GetFileType() uint16 {
	return elf.Header.Type
}

func (elf *ELF64File) GetProgramHeader(index uint16) (ELF64ProgramHeader, error) {
	if int(index) > len(elf.Segments) {
		return ELF64ProgramHeader{}, fmt.Errorf("invalid index %d", index)
	}
	return elf.Segments[index], nil
}

func (elf *ELF64File) GetSectionHeader(index uint16) (ELF64SectionHeader, error) {
	if int(index) > len(elf.Sections) {
		return ELF64SectionHeader{}, fmt.Errorf("invalid index %d", index)
	}
	return elf.Sections[index], nil
}

func (elf *ELF64File) GetSectionContent(index uint16) ([]byte, error) {
	if int(index) > len(elf.Sections) {
		return nil, fmt.Errorf("invalid index")
	}
	sectionStart := elf.Sections[index].Offset
	if int(sectionStart) > len(elf.Data) {
		return nil, fmt.Errorf("section start offset is higher than the length of data: %d vs %d", int(sectionStart), len(elf.Data))
	}
	sectionEnd := uint64(sectionStart) + elf.Sections[index].Size
	if int(sectionEnd) > len(elf.Data) {
		return nil, fmt.Errorf("section end offset is higher than the length of data: %d vs %d", int(sectionEnd), len(elf.Data))
	}
	return elf.Data[sectionStart:sectionEnd], nil
}

func (elf *ELF64File) GetSectionName(index uint16) (string, error) {

	content, err := elf.GetSectionContent(elf.Header.SHStrNdx)
	if err != nil {
		return "", err
	}
	end := elf.Sections[index].Name
	if int(end) > len(elf.Data) {
		return "", fmt.Errorf("name start offset is higher than the length of data: %d vs %d", int(end), len(elf.Data))
	}
	for content[end] != 0 {
		end++
	}
	if int(end) > len(elf.Data) {
		return "", fmt.Errorf("name end offset is higher than the length of data: %d vs %d", int(end), len(elf.Data))
	}
	return string(content[elf.Sections[index].Name:end]), nil
}

func (elf *ELF64File) GetSymbolTable(index uint16) ([]ELF64SymbolTable, []string, error) {
	if elf.Sections[index].Type != SHT_SYMTAB && elf.Sections[index].Type != SHT_DYNSYM {
		return nil, nil, fmt.Errorf("invaild section")
	}
	sectionContent, err := elf.GetSectionContent(index)
	if err != nil {
		return nil, nil, err
	}

	nameTable, err := elf.GetSectionContent(uint16(elf.Sections[index].Link))
	if err != nil {
		return nil, nil, err
	}
	numEntries := elf.Sections[index].Size / uint64(binary.Size(&ELF64SymbolTable{}))
	symbols := make([]ELF64SymbolTable, numEntries)
	data := bytes.NewReader(sectionContent)
	err = binary.Read(data, elf.byteOrder, symbols)
	if err != nil {
		return nil, nil, err
	}

	names := make([]string, numEntries)
	var nameOffset uint32
	for symIdx := range symbols {
		nameOffset = symbols[symIdx].Name
		if nameOffset == 0 {
			names[symIdx] = ""
			continue
		}

		end := nameOffset
		if int(end) > len(elf.Data) {
			return nil, nil, fmt.Errorf("name start offset is higher than the length of data: %d vs %d", int(end), len(elf.Data))
		}
		for nameTable[end] != 0 {
			end++
		}
		if int(end) > len(elf.Data) {
			return nil, nil, fmt.Errorf("name end offset is higher than the length of data: %d vs %d", int(end), len(elf.Data))
		}
		names[symIdx] = string(nameTable[nameOffset:end])
	}
	return symbols, names, nil
}

func (elf *ELF64File) PrettyPrint() {

	// Print Segments
	var i uint16 = 0
	for i = 0; i < elf.Header.PHNum; i++ {
		header, err := elf.GetProgramHeader(i)
		if err != nil {
			log.Fatal("Failed to get header!")
		}
		log.Printf("Segment %d: %s\n", i, header)
	}

	// Print Sections
	for i = 1; i < elf.Header.SHNum; i++ {
		name, err := elf.GetSectionName(i)
		if err != nil {
			log.Fatalf("Invalid section name: %s", err.Error())
		}
		header, err := elf.GetSectionHeader(i)
		if err != nil {
			log.Fatal("Invalid section header")
		}
		log.Printf("Section %d: Name: %s,  Content:%s\n", i, name, header)
	}

	// Print Symbols
	for i = 0; i < elf.Header.SHNum; i++ {
		if elf.Sections[i].Type == SHT_SYMTAB || elf.Sections[i].Type == SHT_DYNSYM {
			name, err := elf.GetSectionName(i)
			if err != nil {
				log.Fatalf("Invalid name: %s", err.Error())
			}
			table, names, e := elf.GetSymbolTable(i)
			if e != nil {
				log.Fatalf("Invalid symbol talbe")
			}
			log.Printf("%d symbols in section %s:\n", len(table), name)
			for j := range table {
				log.Printf("  %d. %s: %s\n", j, names[j], table[j])
			}
		}
	}
}
