package sdk

import (
	"context"
	"fmt"
	proto "github.com/mariodv95/venomlab-sdk/proto/build"

	"github.com/hashicorp/go-plugin"
	"google.golang.org/grpc"
)

// VenomlabPlugin implementa plugin.Plugin
type VenomlabPlugin struct {
	plugin.Plugin
	Impl VenomlabNode
}

func (p *VenomlabPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterVenomlabPluginServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *VenomlabPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewVenomlabPluginClient(c)}, nil
}

// GRPCServer implementa il server gRPC
type GRPCServer struct {
	proto.UnimplementedVenomlabPluginServer
	Impl VenomlabNode
}

func (s *GRPCServer) GetMetadata(ctx context.Context, req *proto.GetMetadataRequest) (*proto.GetMetadataResponse, error) {
	metadata := s.Impl.GetMetadata()

	protoInputs := make([]*proto.IOPort, len(metadata.Inputs))
	for i, input := range metadata.Inputs {
		protoInputs[i] = &proto.IOPort{
			Name:     input.Name,
			Type:     input.Type,
			Optional: input.Optional,
		}
	}

	protoOutputs := make([]*proto.IOPort, len(metadata.Outputs))
	for i, output := range metadata.Outputs {
		protoOutputs[i] = &proto.IOPort{
			Name:     output.Name,
			Type:     output.Type,
			Optional: output.Optional,
		}
	}

	return &proto.GetMetadataResponse{
		Metadata: &proto.NodeMetadata{
			Family:      metadata.Family,
			Name:        metadata.Name,
			DisplayName: metadata.DisplayName,
			Description: metadata.Description,
			Version:     metadata.Version,
			Author:      metadata.Author,
			Inputs:      protoInputs,
			Outputs:     protoOutputs,
		},
	}, nil
}

func (s *GRPCServer) Configure(ctx context.Context, req *proto.ConfigureRequest) (*proto.ConfigureResponse, error) {
	params := make(map[string]interface{})

	for key, param := range req.Params {
		switch v := param.Value.(type) {
		case *proto.ConfigParam_StringValue:
			params[key] = v.StringValue
		case *proto.ConfigParam_IntValue:
			params[key] = int(v.IntValue)
		case *proto.ConfigParam_DoubleValue:
			params[key] = v.DoubleValue
		case *proto.ConfigParam_BoolValue:
			params[key] = v.BoolValue
		case *proto.ConfigParam_BytesValue:
			params[key] = v.BytesValue
		}
	}

	if err := s.Impl.Configure(params); err != nil {
		return &proto.ConfigureResponse{
			Success: false,
			Error:   err.Error(),
		}, nil
	}

	return &proto.ConfigureResponse{Success: true}, nil
}

func (s *GRPCServer) Execute(ctx context.Context, req *proto.ExecuteRequest) (*proto.ExecuteResponse, error) {
	inputs := make(map[string]Data)

	for key, nodeData := range req.Inputs {
		switch data := nodeData.Data.(type) {
		case *proto.NodeData_Shellcode:
			inputs[key] = &ShellcodeData{
				Bytes:             data.Shellcode.Bytes,
				IsEncoded:         data.Shellcode.IsEncoded,
				EncodingTechnique: data.Shellcode.EncodingTechnique,
				Key:               data.Shellcode.Key,
				Metadata:          data.Shellcode.Metadata,
			}
		case *proto.NodeData_SourceCode:
			inputs[key] = &SourceCodeData{
				Files:    data.SourceCode.Files,
				Language: data.SourceCode.Language,
				Metadata: data.SourceCode.Metadata,
			}
		}
	}

	outputs, err := s.Impl.Execute(inputs)
	if err != nil {
		return &proto.ExecuteResponse{Error: err.Error()}, nil
	}

	protoOutputs := make(map[string]*proto.NodeData)
	for key, data := range outputs {
		switch d := data.(type) {
		case *ShellcodeData:
			protoOutputs[key] = &proto.NodeData{
				Data: &proto.NodeData_Shellcode{
					Shellcode: &proto.ShellcodeData{
						Bytes:             d.Bytes,
						IsEncoded:         d.IsEncoded,
						EncodingTechnique: d.EncodingTechnique,
						Key:               d.Key,
						Metadata:          d.Metadata,
					},
				},
			}
		case *SourceCodeData:
			protoOutputs[key] = &proto.NodeData{
				Data: &proto.NodeData_SourceCode{
					SourceCode: &proto.SourceCodeData{
						Files:    d.Files,
						Language: d.Language,
						Metadata: d.Metadata,
					},
				},
			}
		}
	}

	return &proto.ExecuteResponse{Outputs: protoOutputs}, nil
}

// GRPCClient implementa il client gRPC
type GRPCClient struct {
	client proto.VenomlabPluginClient
}

func (c *GRPCClient) GetMetadata() NodeMetadata {
	resp, err := c.client.GetMetadata(context.Background(), &proto.GetMetadataRequest{})
	if err != nil {
		return NodeMetadata{}
	}

	inputs := make([]IOPort, len(resp.Metadata.Inputs))
	for i, input := range resp.Metadata.Inputs {
		inputs[i] = IOPort{
			Name:     input.Name,
			Type:     input.Type,
			Optional: input.Optional,
		}
	}

	outputs := make([]IOPort, len(resp.Metadata.Outputs))
	for i, output := range resp.Metadata.Outputs {
		outputs[i] = IOPort{
			Name:     output.Name,
			Type:     output.Type,
			Optional: output.Optional,
		}
	}

	return NodeMetadata{
		Family:      resp.Metadata.Family,
		Name:        resp.Metadata.Name,
		DisplayName: resp.Metadata.DisplayName,
		Description: resp.Metadata.Description,
		Version:     resp.Metadata.Version,
		Author:      resp.Metadata.Author,
		Inputs:      inputs,
		Outputs:     outputs,
	}
}

func (c *GRPCClient) Configure(params map[string]interface{}) error {
	protoParams := make(map[string]*proto.ConfigParam)

	for key, value := range params {
		switch v := value.(type) {
		case string:
			protoParams[key] = &proto.ConfigParam{Value: &proto.ConfigParam_StringValue{StringValue: v}}
		case int:
			protoParams[key] = &proto.ConfigParam{Value: &proto.ConfigParam_IntValue{IntValue: int64(v)}}
		case float64:
			protoParams[key] = &proto.ConfigParam{Value: &proto.ConfigParam_DoubleValue{DoubleValue: v}}
		case bool:
			protoParams[key] = &proto.ConfigParam{Value: &proto.ConfigParam_BoolValue{BoolValue: v}}
		case []byte:
			protoParams[key] = &proto.ConfigParam{Value: &proto.ConfigParam_BytesValue{BytesValue: v}}
		}
	}

	resp, err := c.client.Configure(context.Background(), &proto.ConfigureRequest{Params: protoParams})
	if err != nil {
		return err
	}

	if !resp.Success {
		return fmt.Errorf(resp.Error)
	}

	return nil
}

func (c *GRPCClient) Execute(inputs map[string]Data) (map[string]Data, error) {
	protoInputs := make(map[string]*proto.NodeData)

	for key, data := range inputs {
		switch d := data.(type) {
		case *ShellcodeData:
			protoInputs[key] = &proto.NodeData{
				Data: &proto.NodeData_Shellcode{
					Shellcode: &proto.ShellcodeData{
						Bytes:             d.Bytes,
						IsEncoded:         d.IsEncoded,
						EncodingTechnique: d.EncodingTechnique,
						Key:               d.Key,
						Metadata:          d.Metadata,
					},
				},
			}
		case *SourceCodeData:
			protoInputs[key] = &proto.NodeData{
				Data: &proto.NodeData_SourceCode{
					SourceCode: &proto.SourceCodeData{
						Files:    d.Files,
						Language: d.Language,
						Metadata: d.Metadata,
					},
				},
			}
		}
	}

	resp, err := c.client.Execute(context.Background(), &proto.ExecuteRequest{Inputs: protoInputs})
	if err != nil {
		return nil, err
	}

	if resp.Error != "" {
		return nil, fmt.Errorf(resp.Error)
	}

	outputs := make(map[string]Data)
	for key, nodeData := range resp.Outputs {
		switch data := nodeData.Data.(type) {
		case *proto.NodeData_Shellcode:
			outputs[key] = &ShellcodeData{
				Bytes:             data.Shellcode.Bytes,
				IsEncoded:         data.Shellcode.IsEncoded,
				EncodingTechnique: data.Shellcode.EncodingTechnique,
				Key:               data.Shellcode.Key,
				Metadata:          data.Shellcode.Metadata,
			}
		case *proto.NodeData_SourceCode:
			outputs[key] = &SourceCodeData{
				Files:    data.SourceCode.Files,
				Language: data.SourceCode.Language,
				Metadata: data.SourceCode.Metadata,
			}
		}
	}

	return outputs, nil
}
