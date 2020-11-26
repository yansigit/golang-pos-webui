package main

import (
	"crypto/md5"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"koisk-noti-desktop/data"
)

func main() {
	initWeb()
}

func initWeb() {
	app := iris.New()

	tmpl := iris.Django("./templates", ".html")
	tmpl.Reload(true)                             // reload templates on each request (development mode)
	tmpl.AddFunc("greet", func(s string) string { // {{greet(name)}}
		return "Greetings " + s + "!"
	})

	// tmpl.RegisterFilter("myFilter", myFilter) // {{"simple input for filter"|myFilter}}
	app.RegisterView(tmpl)
	app.HandleDir("/assets", "./assets")

	app.Get("/", index)
	app.Get("/queue", queue)
	app.Post("/", storeOrderList)
	app.Connect("/", func(context *context.Context) {
		_, _ = context.ResponseWriter().Write([]byte("Hello World"))
	})

	// http://localhost:8080
	_ = app.Listen(":8080")
}

func index(ctx iris.Context) {
	var orders []data.Order
	data.Paging(1, &orders)
	_ = ctx.View("index.html", iris.Map{
		"order_list": orders,
	})
}

func queue(ctx iris.Context) {
	_ = ctx.View("queue.html", iris.Map{
		"confirmedOrders": []int{1123, 1124},
		"waitingOrders":   []int{1125, 1126, 1127, 1128},
	})
}

func storeOrderList(ctx iris.Context) {
	test, _ := ctx.GetBody()
	fmt.Printf("%x\n", md5.Sum(test))
	id := data.InsertOrderList(test)

	response := iris.Map{"state": "OK", "orderNumber": id}
	ctx.JSON(response)
	println(response)
}
