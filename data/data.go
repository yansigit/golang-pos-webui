package data

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
)

type OrderList struct {
	OrderNumber uint `gorm:"primaryKey;unique"`
	TotalPrice  uint `gorm:"not null"`
	Menus       string/* 예시 : "아메리카노/에스프레소/카푸치노" */ `gorm:"not null"`
	Options     string /* 예시 : "샷 추가: 1,사이즈: L, 온도: Hot/없음/우유 추가: 1" */
	IsConfirmed bool   `gorm:"default:false"`
}

var db *gorm.DB

func init() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Println("Db 연결에 실패하였습니다.")
		log.Fatal(err.Error())
	}

	// 테이블 자동 생성
	err = db.AutoMigrate(&OrderList{})
	if err != nil {
		panic("DB 초기화에 실패했습니다")
	}
}

func InsertOrderList(orderList *OrderList) {
	// orderList = &OrderList{OrderNumber: 1524, Price: 1600, Detail: "아이스아메리카노,휘핑크림:X"}
	// 생성

	db.Create(orderList)
}

func FindOrderList(orderNumber uint) (orderList OrderList) {
	db.First(&orderList, orderNumber)
	return orderList
}

// 컨펌 상태만 업데이트
func UpdateOrderListConfirmation(orderNumber uint) {
	orderList := FindOrderList(orderNumber)
	db.Model(&orderList).Update("IsConfirmed", !orderList.IsConfirmed)
}

func DeleteOrderList(orderNumber uint) {
	orderList := FindOrderList(orderNumber)
	db.Delete(&orderList)
}
