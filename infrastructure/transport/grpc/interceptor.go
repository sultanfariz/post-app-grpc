package grpc

import (
	"context"
	"fmt"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	JWT_METADATA_KEY = "authorization"
)

var jwtMethodsAuth = map[string]bool{
	"/user.AuthService/Register":    false,
	"/user.AuthService/Login":       false,
	"/post.PostService/GetPosts":    false,
	"/post.PostService/GetPostById": false,
	"/post.PostService/CreatePost":  true,
}

func JWTInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// check if the method is required to be authenticated
	checkJWT, ok := jwtMethodsAuth[info.FullMethod]
	if !ok || !checkJWT {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	token := md.Get(JWT_METADATA_KEY)
	if len(token) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "token is not provided")
	}

	// parse and validate token
	jwtToken := strings.Split(token[0], " ")
	parsedToken, err := jwt.Parse(jwtToken[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.Errorf(codes.Unauthenticated, "invalid token")
		}
		return []byte(viper.GetString("JWT_SECRET_KEY")), nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	// get email from token
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	fmt.Printf("Authenticated user with email: %s\n", claims)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	email, ok := claims["sub"].(string)
	fmt.Printf("Authenticated user with email: %s\n", email)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	// add email to context
	ctx = context.WithValue(ctx, "email", email)

	return handler(ctx, req)
}
