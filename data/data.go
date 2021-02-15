package data

import (
	"encoding/json"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"time"
)

//type OrderList struct {
//	Id 			uint 		`gorm:"primaryKey;unique"`
//	TotalPrice  uint 		`gorm:"not null"`
//	Menus       string	 	`gorm:"not null"`
//	Options     string 		/* 예시 : "샷 추가: 1 (500원), 헤이즐넛 시럽 추가: 1 (500원), 포장: Y (0원) 1" */
//	IsConfirmed bool  		 `gorm:"default:false"`
//}

type Option struct {
	gorm.Model
	//Id       uint   `gorm:"primaryKey;unique;autoIncrement"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	MenuId   uint
}

type Menu struct {
	gorm.Model
	//Id         uint     `gorm:"primaryKey;unique;autoIncrement"`
	Name       string   `json:"name"`
	Options    []Option `gorm:"foreignKey:MenuId;references:ID" json:"options"`
	Price      int      `json:"price"`
	TotalPrice int      `json:"totalPrice"`
	IsTakeOut  bool     `json:"isTakeOut"`
	IsTumbler  bool     `json:"isTumbler"`
	Temp       string   `json:"temp"`
	OrderId    uint
}

type Order struct {
	gorm.Model
	// Id          uint   `gorm:"primaryKey"`
	IsConfirmed    int    `gorm:"default:0"`
	Menus          []Menu `gorm:"foreignKey:OrderId;references:ID" json:"menus"`
	TotalPrice     int    `json:"totalPrice"`
	ApprovalDate   string `gorm:"default:'EMPTY'"`
	ApprovalNumber string `gorm:"default:'EMPTY'"`
}

var Db *gorm.DB

func init() {
	var err error
	Db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Println("Db 연결에 실패하였습니다.")
		log.Fatal(err.Error())
	}

	// 테이블 자동 생성
	err = Db.AutoMigrate(&Order{}, &Menu{}, &Option{})
	// Db.Raw("UPDATE SQLITE_SEQUENCE SET seq = 999 WHERE name = 'orders'")
	if err != nil {
		panic("DB 초기화에 실패했습니다")
	}
}

func GetMenusFromOrder(order Order, menu *[]Menu) {
	_ = Db.Model(order).Association("Menus").Find(menu)
}

func GetOptionsFromMenu(menu Menu, option *[]Option) {
	_ = Db.Model(menu).Association("Options").Find(option)
}

func InsertOrderList(test []byte) uint {
	var order Order
	err := json.Unmarshal(test, &order)
	if err != nil {
		log.Fatal(err)
	}

	Db.Create(&order)
	return order.ID
}

func InsertBogusOrderList() {
	order := Order{
		IsConfirmed: 0,
		Menus:       nil,
		TotalPrice:  0,
	}
	Db.Create(&order)
}

func _InsertOrderList(test []byte) uint {
	var order Order
	err := json.Unmarshal(test, &order)
	if err != nil {
		log.Fatal(err)
	}

	var tempOrder Order
	tempOrder = Order{
		IsConfirmed: 0,
		Menus:       nil,
		TotalPrice:  order.TotalPrice,
	}
	Db.Create(&tempOrder)

	var tempMenus = make([]Menu, len(order.Menus))
	copy(tempMenus, order.Menus)
	for i, _ := range tempMenus {
		tempMenus[i].OrderId = tempOrder.ID
		tempMenus[i].Options = nil
	}
	Db.Create(&tempMenus)

	for i, _ := range tempMenus {
		var tempOptions = make([]Option, len(order.Menus[i].Options))
		copy(tempOptions, order.Menus[i].Options)
		for j, _ := range tempOptions {
			tempOptions[j].MenuId = tempMenus[i].ID
			if tempOptions[j].Quantity > 0 {
				Db.Create(&tempOptions[j])
			}
		}
	}

	tempOrder.Menus = tempMenus
	Db.Save(&tempOrder)

	return tempOrder.ID
}

func FindOrderList(id uint) (orderList Order) {
	Db.First(&orderList, id)
	return orderList
}

func FindOrderListWithStatus(status int, limit int) (orders []Order) {
	Db.Where("is_confirmed = ?", status).Limit(limit).Find(&orders)
	return orders
}

func FindOrderListWithDate(date time.Time) (orders []Order) {
	Db.Raw("SELECT * FROM orders WHERE strftime('%s', updated_at) BETWEEN strftime('%s', ?) AND strftime('%s', ?)", date.Format("2006-01-02"), date.AddDate(0, 0, 1).Format("2006-01-02")).Scan(&orders)
	return orders
}

func Paging(page int, this interface{}) {
	Db.Scopes(paginate(page, 10)).Order("id desc").Find(this)
}

func paginate(page int, pageSize int) func(db *gorm.DB) *gorm.DB {
	if pageSize < 0 {
		pageSize = 0
	}

	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

//Db.Scopes(Paginate(r)).Find(&users)
//Db.Scopes(Paginate(r)).Find(&articles)

// 컨펌 상태만 업데이트
func UpdateOrderListConfirmation(orderNumber uint) {
	orderList := FindOrderList(orderNumber)
	if orderList.IsConfirmed < 2 {
		orderList.IsConfirmed += 1
	}
	Db.Model(&orderList).Update("IsConfirmed", orderList.IsConfirmed)
}

// 주문 취소
func CancelOrderList(orderNumber uint) {
	orderList := FindOrderList(orderNumber)
	orderList.IsConfirmed = 3
	Db.Model(&orderList).Update("IsConfirmed", orderList.IsConfirmed)
}

func DeleteOrderList(orderNumber uint) {
	orderList := FindOrderList(orderNumber)
	Db.Delete(&orderList)
}
