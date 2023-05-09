package controllers

import (
	"SunProject/application/models"
	"SunProject/application/service"
	"SunProject/config"
	pay2 "SunProject/libary/pay"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"net/http"
	"strconv"
)

type prePayParams struct {
	PayType string `json:"pay_type"`
	ProductId int `json:"product_id"`
	Num int `json:"num"`
}

func PrePay(c *gin.Context) {
	if !IsLogin(c) {
		return
	}
	payParma := prePayParams{}
	if err := c.Bind(&payParma); err != nil {
		ApiError(c, &Response{Code: -1, Msg: "参数绑定错误"})
		return
	}
	if payParma.PayType != "alipay" && payParma.PayType != "wxpay" {
		ApiError(c, &Response{Code: -1, Msg: "支付方式错误"})
		return
	}
	if payParma.ProductId == 0 || payParma.Num <= 0 {
		ApiError(c, &Response{Code: -1, Msg: "参数错误"})
		return
	}
	userId := c.MustGet("userId").(int)
	product := models.Product{}
	config.DB.Where("id = ?", payParma.ProductId).First(&product)
	if product.Id == 0 {
		ApiError(c, &Response{Code: -1, Msg: "参数错误"})
		return
	}
	order := models.Order{
		UserId: userId,
		OrderSn: service.GenerateOrderNumber(),
		PayType: payParma.PayType,
		Title: product.Title,
		ProductId: payParma.ProductId,
		Price: product.Amount,
		Num: payParma.Num,
	}
	priceF, _ := strconv.ParseFloat(product.Amount, 64)
	order.OrderAmount = math.Round(priceF * float64(payParma.Num) * 100) / 100
	config.DB.Create(&order)
	if order.Id == 0 {
		ApiError(c, &Response{Code: -1, Msg: "生成订单失败"})
		return
	}
	// 调用支付
	pay := pay2.ApiParam{
		Type: order.PayType,
		Name: order.Title,
		OutTradeNo: order.OrderSn,
		Money: fmt.Sprintf("%.2f", order.OrderAmount),
		//Money: "0.10",
		ClientIp: GetReaIp(c),
		Device: "pc",
	}
	payData, err := pay.CreateOrder()
	if err != nil {
		ApiError(c, &Response{Code: -1, Msg: "发起支付失败"})
		return
	}
	ApiResponse(c, &Response{Data: payData})
}

func Notify(c *gin.Context) {
	queryStringParams := c.Request.URL.Query()
	queryData := make(map[string]string)
	for key, values := range queryStringParams {
		if len(values) > 0 {
			queryData[key] = values[0]
		}
	}
	// TODO: md5加密验证
	if pid, ok := queryData["pid"]; ok {
		if pid != "1063" {
			c.String(http.StatusBadRequest, "fail")
			return
		}
	}
	if status, ok := queryData["trade_status"]; ok {
		if status != "TRADE_SUCCESS" {
			c.String(http.StatusBadRequest, "fail")
			return
		}
	}
	// 处理
	res := service.Notify(queryData["out_trade_no"], queryData["trade_no"])
	if !res {
		c.String(http.StatusBadRequest, "fail")
		return
	}
	c.String(http.StatusOK, "success")
	return
}

func OrderStatus(c *gin.Context)  {
	orderSn := c.Query("order_sn")
	if orderSn == "" {
		ApiError(c, &Response{Msg: "参数错误"})
		return
	}
	orderDetail := struct {
		Title string `json:"title"`
		PayType string `json:"pay_type"`
		OrderSn string `json:"order_sn"`
		OrderAmount string `json:"order_amount"`
		Status int `json:"status"`
	}{}
	result := config.DB.Model(&models.Order{}).Where("order_sn = ?", orderSn).First(&orderDetail)
	if result.Error != nil {
		ApiError(c, &Response{Msg: "订单不存在"})
		return
	}
	ApiResponse(c, &Response{Data: orderDetail})
}