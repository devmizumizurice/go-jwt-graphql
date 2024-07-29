package graph

import (
	"context"
	"errors"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/devmizumizurice/go-jwt-graphql/utils"
	"github.com/golang-jwt/jwt/v5"
)

var Directive DirectiveRoot = DirectiveRoot{
	IsAuthenticated:        IsAuthenticated,
	IsRefreshAuthenticated: IsRefreshAuthenticated,
}

func IsAuthenticated(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	gc, err := utils.GinContextFromContext(ctx)

	if err != nil {
		return nil, errors.New("CONTEXT_ERROR")
	}

	token, err := gc.Cookie("access_token")

	if err != nil {
		return nil, errors.New("MISSING_TOKEN")
	}

	parsedToken, err := utils.VerifyToken(token, false)
	if err != nil {
		return nil, errors.New("INVALID_TOKEN")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if ok && parsedToken.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, errors.New("EXPIRED_TOKEN")
		}
	}

	ctx = context.WithValue(ctx, utils.SubKey, claims["sub"])

	return next(ctx)
}

func IsRefreshAuthenticated(ctx context.Context, obj interface{}, next graphql.Resolver) (res interface{}, err error) {
	gc, err := utils.GinContextFromContext(ctx)

	if err != nil {
		return nil, errors.New("CONTEXT_ERROR")
	}

	token, err := gc.Cookie("refresh_token")

	if err != nil {
		return nil, errors.New("MISSING_TOKEN")
	}

	parsedToken, err := utils.VerifyToken(token, true)
	if err != nil {
		return nil, errors.New("INVALID_TOKEN")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)

	if ok && parsedToken.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, errors.New("EXPIRED_TOKEN")
		}
	}

	ctx = context.WithValue(ctx, utils.SubKey, claims["sub"])

	return next(ctx)
}
