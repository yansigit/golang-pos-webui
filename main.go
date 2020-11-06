package main

import "github.com/kataras/iris/v12"

func main() {
	app := iris.New()

	app.Get()
}

func index(ctx iris.Context) {
	ctx.Text()
}
