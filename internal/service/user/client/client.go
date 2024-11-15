package client

import (
	"github.com/dobyte/due-chat/internal/service/user/pb"
	"github.com/dobyte/due/v2/transport"
	"google.golang.org/grpc"
)

const target = "discovery://user"

func NewClient(fn transport.NewMeshClient) (pb.UserClient, error) {
	client, err := fn(target)
	if err != nil {
		return nil, err
	}

	return pb.NewUserClient(client.Client().(grpc.ClientConnInterface)), nil
}
