package sdk

// VenomlabNode è l'interfaccia che tutti i plugin devono implementare
type VenomlabNode interface {
	// GetMetadata restituisce i metadata del nodo
	GetMetadata() NodeMetadata

	// Configure configura il nodo con i parametri forniti
	Configure(params map[string]interface{}) error

	// Execute esegue la logica del nodo
	Execute(inputs map[string]Data) (map[string]Data, error)
}

// NodeMetadata descrive il plugin
type NodeMetadata struct {
	Family      string
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
	Type     string
	Optional bool
}

// Data è un'interfaccia per i dati processati
type Data interface {
	DataType() string
}

// ShellcodeData rappresenta shellcode
type ShellcodeData struct {
	Bytes             []byte
	IsEncoded         bool
	EncodingTechnique string
	Key               []byte
	Metadata          map[string]string
}

func (s *ShellcodeData) DataType() string { return "shellcode" }

// SourceCodeData rappresenta codice sorgente
type SourceCodeData struct {
	Files    map[string]string // filename -> content
	Language string
	Metadata map[string]string
}

func (s *SourceCodeData) DataType() string { return "source_code" }

// TODO: aggiungere gli stessi DataType che vengono aggiunti nel proto
