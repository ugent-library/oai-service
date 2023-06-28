package grpcserver

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/bufbuild/connect-go"
)

var (
	ErrNoToken      = errors.New("no token provided")
	ErrInvalidToken = errors.New("invalid token")
)

type AuthConfig struct {
	Token string
}

// see https://connect.build/docs/go/interceptors and
// https://connect.build/docs/go/streaming
type authInterceptor struct {
	AuthConfig
}

func NewAuthInterceptor(config AuthConfig) *authInterceptor {
	return &authInterceptor{
		AuthConfig: config,
	}
}

func (a *authInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		if err := a.checkAuth(req.Header()); err != nil {
			return nil, err
		}
		return next(ctx, req)
	})

}

func (ai *authInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return connect.StreamingClientFunc(func(ctx context.Context, spec connect.Spec) connect.StreamingClientConn {
		return next(ctx, spec)
	})
}

func (a *authInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return connect.StreamingHandlerFunc(func(ctx context.Context, conn connect.StreamingHandlerConn) error {
		if err := a.checkAuth(conn.RequestHeader()); err != nil {
			return err
		}
		return next(ctx, conn)
	})
}

func (a *authInterceptor) checkAuth(headers http.Header) *connect.Error {
	token := strings.TrimPrefix(headers.Get("Authorization"), "Bearer ")
	if token == "" {
		return connect.NewError(connect.CodeUnauthenticated, ErrNoToken)
	}
	if token != a.Token {
		return connect.NewError(connect.CodeUnauthenticated, ErrInvalidToken)
	}
	return nil
}
