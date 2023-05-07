package service

import (
	"SunProject/application/models"
	"SunProject/config"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

// GenerateOrderNumber 生成订单号
func GenerateOrderNumber() string {
	rand.Seed(time.Now().UnixNano())

	orderNum := ""

	for i := 0; i < 16; i++ {
		digit := rand.Intn(10)
		orderNum += fmt.Sprintf("%d", digit)
	}
	return orderNum
}

func Notify(orderSn, outTradeNo string) bool {
	orderDetail := struct {
		Id int `json:"id"`
		UserId int `json:"user_id"`
		ProductId int `json:"product_id"`
		OrderSn string `json:"order_sn"`
		OrderAmount string `json:"order_amount"`
		PayType string `json:"pay_type"`
	}{}
	config.DB.Model(&models.Order{}).Where("order_sn = ?", orderSn).Where("status = 0").First(&orderDetail)
	if orderDetail.Id == 0 {
		config.Logger().Error(fmt.Sprintf("订单不存在: %s", orderSn))
		return false
	}
	var product struct {
		Id       int    `json:"id"`
		Title    string `json:"title"`
		Amount   string `json:"amount"`
		Integral int    `json:"integral"`
	}
	config.DB.Model(&models.Product{}).
		Select("id", "title", "amount", "integral").
		Where("id = ?", orderDetail.ProductId).
		First(&product)
	if product.Id == 0 {
		config.Logger().Error(fmt.Sprintf("产品不存在:productId - %d", orderDetail.ProductId))
		return false
	}
	tx := config.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			config.Logger().Error("发放积分失败", r)
			tx.Rollback()
		}
	}()
	if err := tx.Error; err != nil {
		return false
	}
	// 修改订单状态
	orderRes := tx.Model(&models.Order{}).Where("id = ?", orderDetail.Id).Updates(map[string]interface{}{
		"status": 1,
		"pay_at": time.Now().Format("2006-01-02 15-04-05"),
		"out_trade_no": outTradeNo,
	})
	if orderRes.Error != nil {
		tx.Rollback()
		config.Logger().Error(fmt.Sprintf("更新用户积分失败: %s", orderRes.Error))
		return false
	}
	// 增加积分余额
	userRes := tx.Model(&models.User{}).
		Where("id = ?", orderDetail.UserId).
		Update("integral", gorm.Expr("integral + ?", product.Integral))
	if userRes.Error != nil {
		tx.Rollback()
		config.Logger().Error(fmt.Sprintf("更新用户积分失败: %s", userRes.Error))
		return false
	}
	tx.Commit()
	go func(userId, integral int) {
		result := config.DB.Create(&models.IntegralLog{
			Title: "购买积分",
			Integral: integral,
			UserId: userId,
			Type: 1,
		})
		if result.Error != nil {
			config.Logger().Error(fmt.Sprintf("创建用户使用记录失败: %s", result.Error))
		}
	}(orderDetail.UserId, product.Integral)
	return true

}