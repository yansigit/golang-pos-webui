package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"koisk-noti-desktop/data"
	"strconv"
	"time"
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
	app.Get("/orders", refreshedOrderList)
	app.Get("/waits", refreshWaitNumbers)
	app.Get("/out-display", outDisplay)
	app.Get("/jungsan", jungSan)

	app.Post("/", storeOrderList)
	app.Post("/action", action)

	app.Connect("/", func(context *context.Context) {
		_, _ = context.ResponseWriter().Write([]byte("Hello World"))
	})

	// http://localhost:8080
	_ = app.Listen(":8080")
}

func index(ctx iris.Context) {
	_ = ctx.View("index.html")
}

func outDisplay(ctx iris.Context) {
	var confirmedOrders []data.Order
	var unconfirmedOrders []data.Order

	confirmedOrders = data.FindOrderListWithStatus(1, 5)
	unconfirmedOrders = data.FindOrderListWithStatus(0, 10)

	var unconfirmedOrdersArray [][]data.Order
	if len(unconfirmedOrders) > 5 {
		unconfirmedOrdersArray = append(unconfirmedOrdersArray, unconfirmedOrders[:5], unconfirmedOrders[5:])
	} else {
		unconfirmedOrdersArray = append(unconfirmedOrdersArray, unconfirmedOrders)
	}

	_ = ctx.View("displayWaitingNumbers.html", iris.Map{
		"confirmed_orders":         confirmedOrders,
		"unconfirmed_orders_array": unconfirmedOrdersArray,
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

func refreshWaitNumbers(ctx iris.Context) {
	var confirmedOrders []data.Order
	var unconfirmedOrders []data.Order

	confirmedOrders = data.FindOrderListWithStatus(1, 40)
	unconfirmedOrders = data.FindOrderListWithStatus(0, 40)

	var confirmedOrdersMashed [][]data.Order
	length := len(confirmedOrders)
	for i := 0; i < length; i += 3 { // 0,1,2,3,4,5 length:6
		if i >= length {
			break
		}
		if length-i < 3 {
			if length-i == 1 {
				confirmedOrdersMashed = append(confirmedOrdersMashed, []data.Order{confirmedOrders[i]})
			} else if length-i == 2 {
				confirmedOrdersMashed = append(confirmedOrdersMashed, []data.Order{confirmedOrders[i], confirmedOrders[i+1]})
			}
			break
		}
		confirmedOrdersMashed = append(confirmedOrdersMashed, []data.Order{confirmedOrders[i], confirmedOrders[i+1], confirmedOrders[i+2]})
	}

	var unconfirmedOrdersMashed [][]data.Order
	length = len(unconfirmedOrders)
	for i := 0; i < length; i += 3 { // 0,1,2,3,4,5 length:6
		if i >= length {
			break
		}
		if length-i < 3 {
			if length-i == 1 {
				unconfirmedOrdersMashed = append(unconfirmedOrdersMashed, []data.Order{unconfirmedOrders[i]})
			} else if length-i == 2 {
				unconfirmedOrdersMashed = append(unconfirmedOrdersMashed, []data.Order{unconfirmedOrders[i], unconfirmedOrders[i+1]})
			}
			break
		}
		unconfirmedOrdersMashed = append(unconfirmedOrdersMashed, []data.Order{unconfirmedOrders[i], unconfirmedOrders[i+1], unconfirmedOrders[i+2]})
	}

	args := map[string]interface{}{
		"unconfirmed_orders": unconfirmedOrdersMashed,
		"confirmed_orders":   confirmedOrdersMashed,
	}

	buf := new(bytes.Buffer)
	ctx.Application().View(buf, "waitingNumberTable.html", "refreshedOrderList.html", args)
	ctx.WriteString(buf.String())
}

func action(ctx iris.Context) {
	action := ctx.PostValue("action")
	if action == "confirm" {
		orderNumber, _ := strconv.Atoi(ctx.PostValue("orderNumber"))
		data.UpdateOrderListConfirmation(uint(orderNumber))
	}
	if action == "cancel" {
		orderNumber, _ := strconv.Atoi(ctx.PostValue("orderNumber"))
		data.CancelOrderList(uint(orderNumber))
	}
	if action == "insertbogus" {
		data.InsertBogusOrderList()
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

func jungSan(ctx iris.Context) {
	var orders []data.Order
	date := time.Now() //time.Date(2021, time.January, 3, 0, 0, 0, 0, time.UTC)
	orders = data.FindOrderListWithDate(date)

	totalPrice := 0
	canceledCnt, canceledPrice := 0, 0
	discountCnt, discountPrice := 0, 0

	// 메뉴별 수량, 금액 카운트용
	menuTable := make(map[string][]int)

	for i := 0; i < len(orders); i++ {
		data.Db.Model(&orders[i]).Association("Menus").Find(&orders[i].Menus)
		if orders[i].IsConfirmed == 3 {
			canceledCnt, canceledPrice = canceledCnt+1, canceledPrice+orders[i].TotalPrice
		}
		for _, menu := range orders[i].Menus {
			if menu.IsTakeOut {
				discountCnt, discountPrice = discountCnt+1, discountPrice+200
			}
			if menuTable[menu.Name] == nil {
				menuTable[menu.Name] = []int{0, 0}
			}
			menuTable[menu.Name][0] += 1
			menuTable[menu.Name][1] += menu.TotalPrice
		}
		totalPrice += orders[i].TotalPrice
	}

	args := map[string]interface{}{
		"order_list":    orders,
		"date":          date.Format("2006-01-02"),
		"canceledCnt":   canceledCnt,
		"canceledPrice": canceledPrice,
		"discountCnt":   discountCnt,
		"discountPrice": discountPrice,
		"totalPrice":    totalPrice,
		"menuTable":     menuTable,
	}

	_ = ctx.View("jungsan.html", args)
}
