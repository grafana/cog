package plugins

import (
	"encoding/json"
	"fmt"
	"net/rpc"

	"github.com/grafana/codejen"
	"github.com/grafana/cog/pkg/ir"
	"github.com/grafana/cog/pkg/languages"
	"github.com/hashicorp/go-plugin"
)

// Language describes the interface that must be implemented by language plugins.
// Note: the methods are listed in the order in which they will be called.
type Language interface {
	// ValidateConfig receives the configuration for the current language
	// plugin and validates it.
	// Invalid configurations will result in an error while this method will
	// return nil for valid ones.
	ValidateConfig(config map[string]any) error

	// NullableConfig describes some properties of nullable types for a given language.
	// See also: [languages.NullableConfig]
	NullableConfig(config map[string]any) (languages.NullableConfig, error)

	// TransformSchemas modifies the input schemas to make them suitable to
	// the current language specifically.
	//
	// These transformations should only alter the schemas to make them
	// "compatible" with the target language, NOT to add missing elements or
	// fix incorrect ones.
	//
	// Examples of possible transformations:
	//   * [transforms.AnonymousStructsToNamed] can be used for languages that
	//     don't support anonymous structures
	//   * [transforms.DisjunctionOfConstantsToEnum] for languages that don't
	//     support disjunctions (or: union/sum types). See also other
	//     Disjunction* transformations.
	//   * …
	//
	// Returns the transformed schemas, or an error.
	TransformSchemas(config map[string]any, schemas ir.Schemas) (ir.Schemas, error)

	// Generate performs the code generation.
	//
	// In addition to the codegen context, it receives both the "global"
	// configuration (what should be generated) and the language-specific
	// configuration.
	//
	// The context contains intermediate representation (IR) of the schemas and
	// builders.
	//
	// Note: the builders IR is in the context only if builders generation is
	// enabled (ie: codegenConfig.Builders == true)
	Generate(codegenConfig languages.Config, config map[string]any, context languages.Context) (codejen.Files, error)
}

//nolint:gochecknoglobals
var LanguagePluginHandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "COG_LANGUAGE_PLUGIN",
	MagicCookieValue: "joe la frite",
}

var _ Language = (*LanguageRPC)(nil)

// LanguageRPC is an implementation of the Language interface that talks over RPC.
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

type NullableConfigArgs struct {
	Config map[string]any
}

func (g *LanguageRPC) NullableConfig(config map[string]any) (languages.NullableConfig, error) {
	var resp languages.NullableConfig
	err := g.client.Call("Plugin.NullableConfig", &NullableConfigArgs{
		Config: config,
	}, &resp)

	return resp, err
}

type TransformArgs struct {
	Config  map[string]any
	Schemas []byte
}

func (g *LanguageRPC) TransformSchemas(config map[string]any, schemas ir.Schemas) (ir.Schemas, error) {
	// TODO: marshalling `ir.Schemas` with `encoding/gob` is a bit of a pain, so we cheat and json-marshal it first.
	payload, err := json.Marshal(schemas)
	if err != nil {
		return ir.Schemas{}, fmt.Errorf("could not marshal schemas to json: %w", err)
	}

	result := ""
	err = g.client.Call("Plugin.TransformSchemas", &TransformArgs{
		Config:  config,
		Schemas: payload,
	}, &result)
	if err != nil {
		return ir.Schemas{}, err
	}

	transformed := ir.Schemas{}
	if err := json.Unmarshal([]byte(result), &transformed); err != nil {
		return ir.Schemas{}, err
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

type LanguageRPCServer struct {
	// This is the real implementation
	Impl Language
}

func (s *LanguageRPCServer) ValidateConfig(args *ValidateConfigArgs, _ *string) error {
	return s.Impl.ValidateConfig(args.Config)
}

func (s *LanguageRPCServer) NullableConfig(args *NullableConfigArgs, resp *languages.NullableConfig) error {
	config, err := s.Impl.NullableConfig(args.Config)
	if err != nil {
		return err
	}

	*resp = config

	return nil
}

func (s *LanguageRPCServer) TransformSchemas(args *TransformArgs, result *string) error {
	var schemas ir.Schemas

	if err := json.Unmarshal(args.Schemas, &schemas); err != nil {
		return fmt.Errorf("could not unmarshal context from json: %w", err)
	}

	res, err := s.Impl.TransformSchemas(args.Config, schemas)
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

type LanguagePluginRunner struct {
	// Impl Injection
	Impl Language
}

func (p *LanguagePluginRunner) Server(*plugin.MuxBroker) (any, error) {
	return &LanguageRPCServer{Impl: p.Impl}, nil
}

func (p *LanguagePluginRunner) Client(b *plugin.MuxBroker, c *rpc.Client) (any, error) {
	return &LanguageRPC{client: c}, nil
}
