package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCAuthProvider validates a Bearer token and returns the authenticated user ID.
type GRPCAuthProvider interface {
	Verify(token string) (userID string, err error)
}

// GRPCStubProvider accepts any non-empty Bearer token and uses it as the user ID.
// Intended for local development only (AUTH_PROVIDER=stub).
type GRPCStubProvider struct{}

// Verify uses the token directly as the user ID.
func (p *GRPCStubProvider) Verify(token string) (string, error) {
	if token == "" {
		return "", status.Error(codes.Unauthenticated, "missing bearer token")
	}
	return token, nil
}

// UnaryAuthInterceptor returns a gRPC unary interceptor that verifies the Bearer token
// and injects the user ID into the context. HealthService is exempt.
func UnaryAuthInterceptor(provider GRPCAuthProvider) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		if strings.HasSuffix(info.FullMethod, "/HealthService/Check") {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeader := md.Get("authorization")
		if len(authHeader) == 0 || authHeader[0] == "" {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		token := strings.TrimPrefix(authHeader[0], bearerPrefix)

		userID, err := provider.Verify(token)
		if err != nil {
			return nil, err
		}

		ctx = context.WithValue(ctx, userIDKey, userID)
		return handler(ctx, req)
	}
}
