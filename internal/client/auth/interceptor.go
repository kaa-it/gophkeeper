package auth

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// Interceptor provides gRPC interceptors for unary and stream methods to attach authentication tokens.
type Interceptor struct {
	authClient  *Client
	authMethods map[string]bool
}

// NewAuthInterceptor creates a new Interceptor with the given authClient, authMethods, and token refresh duration.
func NewAuthInterceptor(
	authClient *Client,
	authMethods map[string]bool,
	refreshDuration time.Duration,
) *Interceptor {
	interceptor := &Interceptor{
		authClient:  authClient,
		authMethods: authMethods,
	}

	interceptor.scheduleRefreshToken(refreshDuration)

	return interceptor
}

// Unary returns a grpc.UnaryClientInterceptor that attaches an auth token to
// outgoing gRPC requests if the method requires it.
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

// Stream returns a grpc.StreamClientInterceptor. It attaches an authentication token
// to outgoing gRPC streams if the method requires it.
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
	token := i.authClient.AccessToken()
	return metadata.AppendToOutgoingContext(ctx, "authorization", token)
}

func (i *Interceptor) scheduleRefreshToken(refreshDuration time.Duration) {
	go func() {
		<-i.authClient.notifyChannel

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
}

func (i *Interceptor) refreshToken() error {
	accessToken, err := i.authClient.refreshToken()
	if err != nil {
		return err
	}

	log.Printf("received token: %s", accessToken)

	i.authClient.accessToken = accessToken

	return nil
}
