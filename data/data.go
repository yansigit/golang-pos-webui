package data

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"strconv"
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
	Id       uint   `gorm:"primaryKey;unique;autoIncrement"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	MenuId   int
}

type Menu struct {
	gorm.Model
	Id         uint     `gorm:"primaryKey;unique;autoIncrement"`
	Name       string   `json:"name"`
	Options    []Option `json:"options"`
	Price      int      `json:"price"`
	TotalPrice int      `json:"totalPrice"`
	IsTakeOut  bool     `json:"isTakeOut"`
	IsTumbler  bool     `json:"isTumbler"`
	Temp       string   `json:"temp"`
	OrderId    int
}

type Order struct {
	gorm.Model
	Id          uint   `gorm:"primaryKey;unique;autoIncrement"`
	IsConfirmed bool   `gorm:"default:false"`
	Menus       []Menu `json:"menus"`
	TotalPrice  int    `json:"totalPrice"`
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Println("Db 연결에 실패하였습니다.")
		log.Fatal(err.Error())
	}

	// 테이블 자동 생성
	err = db.AutoMigrate(&Order{}, &Menu{}, &Option{})
	if err != nil {
		panic("DB 초기화에 실패했습니다")
	}
}

func InsertOrderList(orderList *Order) *gorm.DB {
	result := db.Create(orderList)
	// print(result)
	return result
}

func FindOrderList(id uint) (orderList Order) {
	db.First(&orderList, id)
	return orderList
}

func Paginate(_page string, _pageSize string) func(db *gorm.DB) *gorm.DB {
	page, _ := strconv.Atoi(_page)
	pageSize, _ := strconv.Atoi(_pageSize)
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

//db.Scopes(Paginate(r)).Find(&users)
//db.Scopes(Paginate(r)).Find(&articles)

// 컨펌 상태만 업데이트
func UpdateOrderListConfirmation(orderNumber uint) {
	orderList := FindOrderList(orderNumber)
	db.Model(&orderList).Update("IsConfirmed", !orderList.IsConfirmed)
}

func DeleteOrderList(orderNumber uint) {
	orderList := FindOrderList(orderNumber)
	db.Delete(&orderList)
}
