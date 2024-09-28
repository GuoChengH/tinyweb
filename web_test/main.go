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

	g.Use(func(next tinyweb.HandleFunc) tinyweb.HandleFunc {
		return func(ctx *tinyweb.Context) {
			fmt.Println("before middleware")
			next(ctx)
			fmt.Println("after middleware")
		}
	})
	g.Get("/users/:id", func(ctx *tinyweb.Context) {
		fmt.Fprintf(ctx.W, "OK user/:id matched\n")
	})
	engine.Run()
}
