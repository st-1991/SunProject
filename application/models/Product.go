package models

type Product struct {
	Id int `json:"id" gorm:"primaryKey"`
	Title string `json:"title" gorm:"type:varchar(64);not null;default:'';comment:产品名称"`
	Desc string `json:"desc" gorm:"type:varchar(64);not null;default:'';comment:产品描述"`
	Flag string `json:"flag" gorm:"type:varchar(20);not null;default:'';comment:标识"`
	Amount string `json:"amount" gorm:"type:decimal(8,2);not null;default:'0.00';comment:金额"`
	Integral int `json:"integral" gorm:"type:int(11);unsigned;not null;default:0;comment:积分"`
	Date `gorm:"embedded"`
}

func (Product) TableName() string {
	return "keep_product"
}