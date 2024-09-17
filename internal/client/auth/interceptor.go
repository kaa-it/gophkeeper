package auth

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Interceptor struct {
	authClient  *Client
	authMethods map[string]bool
	accessToken string
}

func NewAuthInterceptor(
	authClient *Client,
	authMethods map[string]bool,
	refreshDuration time.Duration,
) (*Interceptor, error) {
	interceptor := &Interceptor{
		authClient:  authClient,
		authMethods: authMethods,
	}

	err := interceptor.scheduleRefreshToken(refreshDuration)
	if err != nil {
		return nil, err
	}

	return interceptor, nil
}

func (i *Interceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply any,
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if i.authMethods[method] {
			return invoker(i.attachToken(ctx), method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (i *Interceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		if i.authMethods[method] {
			return streamer(i.attachToken(ctx), desc, cc, method, opts...)
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (i *Interceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", i.accessToken)
}

func (i *Interceptor) scheduleRefreshToken(refreshDuration time.Duration) error {
	err := i.refreshToken()
	if err != nil {
		return err
	}

	go func() {
		wait := refreshDuration
		for {
			<-time.After(wait)
			err := i.refreshToken()
			if err != nil {
				wait = time.Second
			} else {
				wait = refreshDuration
			}
		}
	}()

	return nil
}

func (i *Interceptor) refreshToken() error {
	accessToken, err := i.authClient.Login()
	if err != nil {
		return err
	}

	i.accessToken = accessToken

	return nil
}
