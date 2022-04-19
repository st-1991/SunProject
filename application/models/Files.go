package models

import "SunProject/config"

type Files struct {
	Id int `json:"id" gorm:"primaryKey"`
	FileName string `json:"file_name" gorm:"type:varchar(255);not null;default:'';comment:文件名"`
	FileMd5 string `json:"file_md5" gorm:"index:idx_md5;type:varchar(255);not null;default:'';comment:文件md5"`
	FileUrl string `json:"file_url" gorm:"type:varchar(255);not null;default:'';comment:文件url"`
	Date `gorm:"embedded"`
}

func (f Files) TableName() string {
	return "keep_files"
}

func (f Files) CreateFile() bool {
	return config.DB.Create(&f).Error == nil
}

func (f Files) GetFileByMd5() Files {
	config.DB.Where("file_md5 = ?", f.FileMd5).Select("id", "file_url").First(&f)
	return f
}