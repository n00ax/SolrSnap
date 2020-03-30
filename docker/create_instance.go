package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
)

func CreateInstance() (*client.Client, context.Context, error) {
	fmt.Printf("== Docker Startup ==\n")
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, nil, err
	}
	return cli, ctx, nil
}
