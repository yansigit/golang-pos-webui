package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"koisk-noti-desktop/data"
	"strconv"
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
	app.Get("/test", refreshedOrderList)

	app.Post("/", storeOrderList)
	app.Post("/action", action)

	app.Connect("/", func(context *context.Context) {
		_, _ = context.ResponseWriter().Write([]byte("Hello World"))
	})

	// http://localhost:8080
	_ = app.Listen(":8080")
}

func index(ctx iris.Context) {
	var orders []data.Order
	data.Paging(1, &orders)

	for i, _ := range orders {
		var menus []data.Menu
		data.GetMenusFromOrder(orders[i], &menus)
		orders[i].Menus = make([]data.Menu, len(menus))
		copy(orders[i].Menus, menus)

		for j, _ := range orders[i].Menus {
			var options []data.Option
			data.GetOptionsFromMenu(orders[i].Menus[j], &options)
			orders[i].Menus[j].Options = make([]data.Option, len(options))
			copy(orders[i].Menus[j].Options, options)
		}
	}

	_ = ctx.View("index.html", iris.Map{
		"order_list": orders,
	})
}

func refreshedOrderList(ctx iris.Context) {
	var orders []data.Order
	data.Paging(1, &orders)

	for i, _ := range orders {
		var menus []data.Menu
		data.GetMenusFromOrder(orders[i], &menus)
		orders[i].Menus = make([]data.Menu, len(menus))
		copy(orders[i].Menus, menus)

		for j, _ := range orders[i].Menus {
			var options []data.Option
			data.GetOptionsFromMenu(orders[i].Menus[j], &options)
			orders[i].Menus[j].Options = make([]data.Option, len(options))
			copy(orders[i].Menus[j].Options, options)
		}
	}

	args := map[string]interface{}{
		"order_list": orders,
	}

	buf := new(bytes.Buffer)
	ctx.Application().View(buf, "refreshedOrderList.html", "refreshedOrderList.html", args)
	ctx.WriteString(buf.String())
}

func action(ctx iris.Context) {
	action := ctx.PostValue("action")
	if action == "confirm" {
		orderNumber, _ := strconv.Atoi(ctx.PostValue("orderNumber"))
		data.UpdateOrderListConfirmation(uint(orderNumber))
	}
}

var newOrderAvailable bool = false

func storeOrderList(ctx iris.Context) {
	test, _ := ctx.GetBody()
	fmt.Printf("%x\n", md5.Sum(test))
	id := data.InsertOrderList(test)

	newOrderAvailable = true

	response := iris.Map{"state": "OK", "orderNumber": id}
	ctx.JSON(response)
	println(response)
}

func queue(ctx iris.Context) {
	ctx.JSON(iris.Map{"new": newOrderAvailable})
	newOrderAvailable = false
}
