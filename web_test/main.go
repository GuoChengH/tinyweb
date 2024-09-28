package main

import (
	"fmt"

	"github.com/GuoChengH/tinyweb"
)

func main() {

	engine := tinyweb.New()

	g := engine.Group("/api")

	// g.Get("/users", func(ctx *tinyweb.Context) {
	// 	fmt.Fprintf(ctx.W, "users\n")
	// })
	g.Get("/users/:id", func(ctx *tinyweb.Context) {
		fmt.Fprintf(ctx.W, "OK user/:id matched\n")
	})

	engine.Run()
}
