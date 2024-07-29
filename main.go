package main

import (
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/devmizumizurice/go-jwt-graphql/graph"
	"github.com/devmizumizurice/go-jwt-graphql/initializers"
	"github.com/devmizumizurice/go-jwt-graphql/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.SetUpDB()
	initializers.SyncDB()
}

func graphqlHandler() gin.HandlerFunc {
	services := InitializeService()

	h := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers:  &graph.Resolver{Srv: services},
				Directives: graph.Directive,
			},
		),
	)
	h.Use(extension.FixedComplexityLimit(10))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	gin.SetMode(gin.DebugMode)

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil)

	r.Use(
		cors.New(
			cors.Config{
				AllowOrigins: []string{
					"https://localhost:8080/",
					"https://localhost:3000/",
				},
				AllowMethods: []string{
					"POST",
					"GET",
					"OPTIONS",
				},
				AllowHeaders: []string{
					"Access-Control-Allow-Credentials",
					"Access-Control-Allow-Headers",
					"Content-Type",
					"Content-Length",
					"Accept-Encoding",
					"Authorization",
				},
				AllowCredentials: true,
				MaxAge:           12 * time.Hour,
			},
		),
	)
	r.Use(utils.GinContextToContextMiddleware())

	r.POST("/query", graphqlHandler())
	if gin.Mode() != gin.ReleaseMode {
		r.GET("/", playgroundHandler())
	}

	r.Run("localhost:3000")
}
