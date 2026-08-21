package plugin

import (
	"encoding/json"
	"fmt"
	"net/rpc"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/pkg/languages"
	"github.com/hashicorp/go-plugin"
)

type Language interface {
	ValidateConfig(config map[string]any) error
	Transform(codegenConfig languages.Config, config map[string]any, context languages.Context) (languages.Context, error)
	Generate(codegenConfig languages.Config, config map[string]any, context languages.Context) (codejen.Files, error)
}

var LanguagePluginHandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "COG_LANGUAGE_PLUGIN",
	MagicCookieValue: "joe la frite",
}

var _ Language = (*LanguageRPC)(nil)

// Here is an implementation that talks over RPC
type LanguageRPC struct {
	client *rpc.Client
}

type ValidateConfigArgs struct {
	Config map[string]any
}

func (g *LanguageRPC) ValidateConfig(config map[string]any) error {
	var resp string
	err := g.client.Call("Plugin.ValidateConfig", &ValidateConfigArgs{
		Config: config,
	}, &resp)

	return err
}

type TransformArgs struct {
	CodegenConfig languages.Config
	Config        map[string]any
	Context       []byte
}

func (g *LanguageRPC) Transform(codegenConfig languages.Config, config map[string]any, context languages.Context) (languages.Context, error) {
	transformed := languages.Context{}

	// TODO: marshalling `languages.Context` with `encoding/gob` is a bit of a pain, so we cheat and json-marshal it first.
	ctxPayload, err := json.Marshal(context)
	if err != nil {
		return languages.Context{}, fmt.Errorf("could not marshal context to json: %w", err)
	}

	result := ""
	err = g.client.Call("Plugin.Transform", &TransformArgs{
		CodegenConfig: codegenConfig,
		Config:        config,
		Context:       ctxPayload,
	}, &result)
	if err != nil {
		return languages.Context{}, err
	}

	if err := json.Unmarshal([]byte(result), &transformed); err != nil {
		return languages.Context{}, err
	}

	return transformed, nil
}

type GenerateArgs struct {
	CodegenConfig languages.Config
	Config        map[string]any
	Context       []byte
}

func (g *LanguageRPC) Generate(codegenConfig languages.Config, config map[string]any, context languages.Context) (codejen.Files, error) {
	files := codejen.Files{}

	// TODO: marshalling `languages.Context` with `encoding/gob` is a bit of a pain, so we cheat and json-marshal it first.
	ctxPayload, err := json.Marshal(context)
	if err != nil {
		return nil, fmt.Errorf("could not marshal context to json: %w", err)
	}

	err = g.client.Call("Plugin.Generate", &GenerateArgs{
		CodegenConfig: codegenConfig,
		Config:        config,
		Context:       ctxPayload,
	}, &files)
	if err != nil {
		return nil, err
	}

	return files, err
}

// Here is the RPC server that GreeterRPC talks to, conforming to
// the requirements of net/rpc
type LanguageRPCServer struct {
	// This is the real implementation
	Impl Language
}

func (s *LanguageRPCServer) ValidateConfig(args *ValidateConfigArgs, resp *string) error {
	return s.Impl.ValidateConfig(args.Config)
}

func (s *LanguageRPCServer) Transform(args *TransformArgs, result *string) error {
	var codegenContext languages.Context

	if err := json.Unmarshal(args.Context, &codegenContext); err != nil {
		return fmt.Errorf("could not unmarshal context from json: %w", err)
	}

	res, err := s.Impl.Transform(args.CodegenConfig, args.Config, codegenContext)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(res)
	if err != nil {
		return err
	}

	*result = string(payload)

	return nil
}

func (s *LanguageRPCServer) Generate(args *GenerateArgs, files *codejen.Files) error {
	var codegenContext languages.Context

	if err := json.Unmarshal(args.Context, &codegenContext); err != nil {
		return fmt.Errorf("could not unmarshal context from json: %w", err)
	}

	res, err := s.Impl.Generate(args.CodegenConfig, args.Config, codegenContext)
	if err != nil {
		return err
	}

	*files = res

	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a GreeterRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return GreeterRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type LanguagePlugin struct {
	// Impl Injection
	Impl Language
}

func (p *LanguagePlugin) Server(*plugin.MuxBroker) (any, error) {
	return &LanguageRPCServer{Impl: p.Impl}, nil
}

func (p *LanguagePlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (any, error) {
	return &LanguageRPC{client: c}, nil
}
