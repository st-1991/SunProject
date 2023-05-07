package models

import "SunProject/application/models/custom"

type Order struct {
	Id int `json:"id" gorm:"primaryKey"`
	OrderType int `json:"order_type" gorm:"type:tinyint(1);not null;default:1;comment:订单类型1积分"`
	UserId int `json:"user_id" gorm:"type:int(11);noll null;default:0;"`
	Title string `json:"title" gorm:"type:varchar(126);not null;default:'';comment:订单标题"`
	ProductId int `json:"product_id" gorm:"type:int(11);not null;default:0;comment:产品id"`
	Price string `json:"price" gorm:"type:decimal(8,2);not null;default:0.00;comment:单价"`
	Num int `json:"num" gorm:"type:int(11);not null;default:0;comment:数量"`
	OrderAmount float64 `json:"order_amount" gorm:"type:decimal(8,2);not null;default:0.00;comment:订单金额"`
	PayType string `json:"pay_type" gorm:"type:varchar(16);not null;default:'';comment:支付方式 alipay,wxpay"`
	OrderSn string `json:"order_sn" gorm:"uniqueIndex:uniq_order;type:char(16);not null;default:'';comment:订单号"`
	OutTradeNo string `json:"out_trade_no" gorm:"type:varchar(32);not null;default:'';comment:三方订单号"`
	PayAt custom.JTime `json:"pay_at" gorm:"type:timestamp;not null;default:'1997-01-01 00:00:00';comment:支付时间"`
	Status int `json:"status" gorm:"type:tinyint(1);not null;default:0;comment:订单状态0发起支付1已完成"`
	Date `gorm:"embedded"`
}

func (Order) TableName() string {
	return "keep_orders"
}