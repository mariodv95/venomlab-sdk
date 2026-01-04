package example_plugin

import (
	"fmt"
	sdk "github.com/mariodv95/venomlab-sdk"
)

type MyAwesomeNode struct {
	config string
}

func NewMyAwesomeNode() *MyAwesomeNode {
	return &MyAwesomeNode{}
}

func (m *MyAwesomeNode) GetMetadata() sdk.NodeMetadata {
	return sdk.NodeMetadata{
		Family:      "encoding",
		Name:        "my_awesome_encoder",
		DisplayName: "My Awesome Encoder",
		Description: "Does something cool to shellcode",
		Version:     "1.0.0",
		Author:      "Developer Name",
		Inputs: []sdk.IOPort{
			{Name: "shellcode", Type: "shellcode", Optional: false},
		},
		Outputs: []sdk.IOPort{
			{Name: "shellcode", Type: "shellcode", Optional: false},
		},
	}
}

func (m *MyAwesomeNode) Configure(params map[string]interface{}) error {
	if cfg, ok := params["config"]; ok {
		m.config = cfg.(string)
	}
	return nil
}

func (m *MyAwesomeNode) Execute(inputs map[string]sdk.Data) (map[string]sdk.Data, error) {
	shellcode, ok := inputs["shellcode"].(*sdk.ShellcodeData)
	if !ok {
		return nil, fmt.Errorf("invalid input")
	}

	// Logica del plugin
	encoded := make([]byte, len(shellcode.Bytes))
	for i, b := range shellcode.Bytes {
		encoded[i] = b ^ 0xAA // XOR semplice come esempio
	}

	result := &sdk.ShellcodeData{
		Bytes:             encoded,
		IsEncoded:         true,
		EncodingTechnique: "XOR-0xAA",
		Metadata:          map[string]string{"key": "0xAA"},
	}

	return map[string]sdk.Data{"shellcode": result}, nil
}
