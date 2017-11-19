//go:generate sh api.sh

package api

import (
	"context"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

type faas struct{}

func newFaasService() FaasCliServiceServer {
	return new(faas)
}

func newResponse(code int) *Response {
	return &Response{
		Code:    int32(code),
		Message: http.StatusText(code),
	}
}

func okResponse() *Response {
	return newResponse(http.StatusOK)
}

func badRequestResponse() *Response {
	return newResponse(http.StatusBadRequest)
}

//Run start the gRPC server
func Run() error {
	listen, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}
	server := grpc.NewServer()
	RegisterFaasCliServiceServer(server, newFaasService())
	server.Serve(listen)
	return nil
}

func (f *faas) Build(ctx context.Context, msg *BuildRequest) (*Response, error) {
	return okResponse(), nil
}

func (f *faas) Deploy(ctx context.Context, msg *DeployRequest) (*Response, error) {
	return okResponse(), nil
}

func (f *faas) Invoke(ctx context.Context, msg *InvokeRequest) (*Response, error) {
	return okResponse(), nil
}

func (f *faas) List(ctx context.Context, msg *ListRequest) (*Response, error) {
	return okResponse(), nil
}

func (f *faas) Login(ctx context.Context, msg *LoginRequest) (*Response, error) {
	return okResponse(), nil
}

func (f *faas) Logout(ctx context.Context, msg *LogoutRequest) (*Response, error) {
	return okResponse(), nil
}

func (f *faas) NewFunction(ctx context.Context, msg *NewFunctionRequest) (*Response, error) {
	return okResponse(), nil
}

func (f *faas) Push(ctx context.Context, msg *PushRequest) (*Response, error) {
	return okResponse(), nil
}

func (f *faas) TemplatePull(ctx context.Context, msg *TemplatePullRequest) (*Response, error) {
	return okResponse(), nil
}

func (f *faas) Version(ctx context.Context, msg *VersionRequest) (*Response, error) {
	return okResponse(), nil
}
