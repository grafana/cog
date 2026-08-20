package plugin

import (
	"net/rpc"

	"github.com/grafana/codejen"
	"github.com/hashicorp/go-plugin"
)

type Language interface {
	ValidateConfig(config map[string]any) error
	Generate() (codejen.Files, error)
}

var LanguagePluginHandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "COG_LANGUAGE_PLUGIN",
	MagicCookieValue: "joe la frite",
}

// Here is an implementation that talks over RPC
type LanguageRPC struct{ client *rpc.Client }

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

type GenerateArgs struct {
}

func (g *LanguageRPC) Generate() (codejen.Files, error) {
	files := codejen.Files{}

	err := g.client.Call("Plugin.Generate", &GenerateArgs{}, &files)
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

func (s *LanguageRPCServer) Generate(_ *GenerateArgs, files *codejen.Files) error {
	res, err := s.Impl.Generate()
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
