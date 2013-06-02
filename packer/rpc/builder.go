package rpc

import (
	"github.com/mitchellh/packer/packer"
	"net/rpc"
)

// An implementation of packer.Builder where the builder is actually executed
// over an RPC connection.
type builder struct {
	client *rpc.Client
}

// BuilderServer wraps a packer.Builder implementation and makes it exportable
// as part of a Golang RPC server.
type BuilderServer struct {
	builder packer.Builder
}

type BuilderPrepareArgs struct {
	Config interface{}
}

type BuilderRunArgs struct {
	RPCAddress string
}

func Builder(client *rpc.Client) *builder {
	return &builder{client}
}

func (b *builder) Prepare(config interface{}) (err error) {
	cerr := b.client.Call("Builder.Prepare", &BuilderPrepareArgs{config}, &err)
	if cerr != nil {
		err = cerr
	}

	return
}

func (b *builder) Run(ui packer.Ui, hook packer.Hook) packer.Artifact {
	// Create and start the server for the Build and UI
	// TODO: Error handling
	server := rpc.NewServer()
	RegisterUi(server, ui)
	RegisterHook(server, hook)

	args := &BuilderRunArgs{serveSingleConn(server)}

	var reply string
	if err := b.client.Call("Builder.Run", args, &reply); err != nil {
		panic(err)
	}

	client, err := rpc.Dial("tcp", reply)
	if err != nil {
		panic(err)
	}

	return Artifact(client)
}

func (b *BuilderServer) Prepare(args *BuilderPrepareArgs, reply *error) error {
	err := b.builder.Prepare(args.Config)
	if err != nil {
		*reply = NewBasicError(err)
	}

	return nil
}

func (b *BuilderServer) Run(args *BuilderRunArgs, reply *string) error {
	client, err := rpc.Dial("tcp", args.RPCAddress)
	if err != nil {
		return err
	}

	hook := Hook(client)
	ui := &Ui{client}
	artifact := b.builder.Run(ui, hook)

	// Wrap the artifact
	server := rpc.NewServer()
	RegisterArtifact(server, artifact)

	*reply = serveSingleConn(server)
	return nil
}