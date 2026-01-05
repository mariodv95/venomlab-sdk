package sdk

import (
	"fmt"
	"github.com/google/uuid"
)

// VenomlabNode è l'interfaccia che tutti i plugin devono implementare
type VenomlabNode interface {
	// GetMetadata restituisce i metadata del nodo
	GetMetadata() NodeMetadata

	// Configure configura il nodo con i parametri forniti
	Configure(params map[string]interface{}) error

	// Execute esegue la logica del nodo
	Execute(inputs map[string]Data) (map[string]Data, error)

	// GetID ritorna l'ID univoco dell'istanza del nodo
	GetID() string

	// SetID imposta l'ID del nodo
	SetID(id string)
}

// BaseNode fornisce implementazione base per tutti i nodi
type BaseNode struct {
	ID       string
	Metadata NodeMetadata
	Params   map[string]interface{}
}

func NewBaseNode(metadata NodeMetadata) *BaseNode {
	return &BaseNode{
		ID:       uuid.New().String(),
		Metadata: metadata,
		Params:   make(map[string]interface{}),
	}
}

func (b *BaseNode) GetMetadata() NodeMetadata {
	return b.Metadata
}

func (b *BaseNode) GetID() string {
	return b.ID
}

func (b *BaseNode) SetID(id string) {
	b.ID = id
}

func (b *BaseNode) Configure(params map[string]interface{}) error {
	b.Params = params
	return nil
}

// ValidateInput verifica che l'input sia del tipo accettato
func (b *BaseNode) ValidateInputs(inputs map[string]Data) error {
	for _, port := range b.Metadata.Inputs {
		// Se l'input è opzionale e non presente, skip
		if port.Optional {
			if _, exists := inputs[port.Name]; !exists {
				continue
			}
		}

		// Input richiesto: deve essere presente
		data, exists := inputs[port.Name]
		if !exists {
			return fmt.Errorf("required input '%s' not provided", port.Name)
		}

		// Valida il Data
		if err := data.Validate(); err != nil {
			return fmt.Errorf("input '%s' validation failed: %w", port.Name, err)
		}

		// Verifica compatibilità del tipo
		if data.Type() != port.Type {
			return fmt.Errorf("input '%s': expected type %s, got %s",
				port.Name, port.Type, data.Type())
		}
	}

	return nil
}

// GetInputPort ritorna l'IOPort con il nome specificato, o errore se non esiste
func (b *BaseNode) GetInputPort(name string) (*IOPort, error) {
	for _, port := range b.Metadata.Inputs {
		if port.Name == name {
			return &port, nil
		}
	}
	return nil, fmt.Errorf("input port '%s' not found", name)
}

// GetOutputPort ritorna l'IOPort con il nome specificato, o errore se non esiste
func (b *BaseNode) GetOutputPort(name string) (*IOPort, error) {
	for _, port := range b.Metadata.Outputs {
		if port.Name == name {
			return &port, nil
		}
	}
	return nil, fmt.Errorf("output port '%s' not found", name)
}
