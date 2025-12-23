// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"
	"fungo/articles/api/internal/config"
	"fungo/articles/api/internal/handler"
	"fungo/articles/api/internal/svc"
	"fungo/common/jwts"
	"fungo/common/middleware"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/article.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	c.MaxBytes = 10485760

	server := rest.MustNewServer(c.RestConf, rest.WithUnauthorizedCallback(jwts.JwtUnauthorizedResult))
	server.Use(middleware.OptionalJWT(c.Auth.AccessSecret))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
