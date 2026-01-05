package sdk

import (
	"fmt"
)

// NodeFamily rappresenta la famiglia/categoria di un nodo
type NodeFamily string

const (
	FamilyInput       NodeFamily = "input"
	FamilyEncoding    NodeFamily = "encoding"
	FamilyInjection   NodeFamily = "injection"
	FamilyObfuscation NodeFamily = "obfuscation"
)

func (nf NodeFamily) String() string {
	return string(nf)
}

func NewNodeFamily(s string) NodeFamily {
	return NodeFamily(s)
}

// DataType rappresenta i tipi di dati che possono fluire tra i nodi
type DataType string

const (
	DataTypeShellcode  DataType = "shellcode"   // Raw shellcode bytes
	DataTypeSourceCode DataType = "source_code" // C/C++ source files
	DataTypeBinary     DataType = "binary"      // Compiled executable
)

func (dt DataType) String() string {
	return string(dt)
}

func NewDataType(s string) DataType {
	return DataType(s)
}

// NodeMetadata descrive il plugin
type NodeMetadata struct {
	Family      NodeFamily
	Name        string
	DisplayName string
	Description string
	Version     string
	Author      string
	Inputs      []IOPort
	Outputs     []IOPort
}

// IOPort descrive una porta di input/output
type IOPort struct {
	Name     string
	Type     DataType
	Optional bool
}

// Data è un'interfaccia per i dati processati
type Data interface {
	Type() DataType
	Validate() error
}

// ShellcodeData rappresenta shellcode raw
type ShellcodeData struct {
	Bytes             []byte            // Raw shellcode bytes
	IsEncoded         bool              // Se true, lo shellcode è codificato
	EncodingTechnique string            // Encoding technique used (se applicabile)
	Key               []byte            // Encoding key (se applicabile)
	Metadata          map[string]string // Metadata opzionale (arch, format, etc)
}

func (s *ShellcodeData) Type() DataType { return DataTypeShellcode }

func (s *ShellcodeData) Validate() error {
	if len(s.Bytes) == 0 {
		return fmt.Errorf("shellcode is empty")
	}
	return nil
}

// SourceCodeData rappresenta codice sorgente generato
type SourceCodeData struct {
	Files    map[string]string // filename -> content
	Language string            // "c", "cpp" // potenzialmente potrebbe supportare qualsiasi linguaggio, al netto di template e moduli
	Metadata map[string]string
}

func (s *SourceCodeData) Type() DataType { return DataTypeSourceCode }

func (s *SourceCodeData) Validate() error {
	if len(s.Files) == 0 {
		return fmt.Errorf("no source files provided or generated")
	}
	return nil
}

// BinaryData rappresenta un eseguibile compilato
type BinaryData struct {
	Data     []byte
	Platform string // "windows", "linux"
	Arch     string // "x64", "x86"
	Metadata map[string]string
}

func (b *BinaryData) Type() DataType { return DataTypeBinary }

func (b *BinaryData) Validate() error {
	if len(b.Data) == 0 {
		return fmt.Errorf("binary is empty")
	}
	return nil
}
