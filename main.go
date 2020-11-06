package main

import (
	"github.com/kataras/iris/v12"
	"time"
)

func main() {
	app := iris.New()

	tmpl := iris.Django("./templates", ".html")
	tmpl.Reload(true)                             // reload templates on each request (development mode)
	tmpl.AddFunc("greet", func(s string) string { // {{greet(name)}}
		return "Greetings " + s + "!"
	})

	// tmpl.RegisterFilter("myFilter", myFilter) // {{"simple input for filter"|myFilter}}
	app.RegisterView(tmpl)

	app.Get("/", index)

	// http://localhost:8080
	app.Listen(":8080")
}

var startTime = time.Now()

func index(ctx iris.Context) {
	ctx.View("hi.html", iris.Map{
		"title":           "Hi Page",
		"name":            "iris",
		"serverStartTime": startTime,
	})
}
