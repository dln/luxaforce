package server

import (
	"context"
	"errors"
	"github.com/dln/luxaforce/api"
)

type luxaforceAPI struct {
	state string
}

func NewAPI() api.LuxaforceServer {
	return luxaforceAPI{
		state: "foo",
	}
}

func (d luxaforceAPI) CreateClient(ctx context.Context, req *api.CreateClientReq) (*api.CreateClientResp, error) {
	if req.Client == nil {
		return nil, errors.New("no client supplied")
	}
	return nil, nil
}
func (d luxaforceAPI) DeleteClient(ctx context.Context, req *api.DeleteClientReq) (*api.DeleteClientResp, error) {
	return nil, nil
}

func (d luxaforceAPI) GetVersion(ctx context.Context, req *api.VersionReq) (*api.VersionResp, error) {
	return &api.VersionResp{
		Server: "0.0.1",
		Api:    1,
	}, nil
}
