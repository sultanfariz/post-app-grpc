package grpc

import (
	"context"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	JWT_METADATA_KEY = "jwt"
	SECRET_KEY       = "thisIs45ecretKey"
)

func JWTInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	token := md.Get(JWT_METADATA_KEY)
	if len(token) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "token is not provided")
	}

	parsedToken, err := jwt.Parse(token[0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}
		return []byte(SECRET_KEY), nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	if !parsedToken.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	
	return handler(ctx, req)
}
