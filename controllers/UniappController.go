package controllers

import (
	"SunProject/libary/uniapp"
	"fmt"
	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["files[]"]
	//filenames := make([]string, len(files))
	//for _, file := range files {
	//	//filename := file.Filename
	//	//if filename != "" {
	//	filenames = append(filenames, file.Filename)
	//	log.Println(file.Filename)
	//
	//	dst := "./" + file.Filename
	//	// 上传文件至指定的完整文件路径
	//	c.SaveUploadedFile(file, dst)
	//	//}
	//}
	var errs []error
	var fileUrls []string
	for _, file := range files {
		UniApp := uniapp.UniApp{}
		AliYun, err := UniApp.InitConfig().PreUploadFile(file.Filename)
		if err != nil {
			errs = append(errs, fmt.Errorf("预加载失败：%s", err))
		}
		_, err = AliYun.Upload(file)
		if err != nil {
			errs = append(errs, fmt.Errorf("上传文件失败：%s", err))
		}
		d := UniApp.CompleteUploadFile(AliYun.ID)
		if d == false {
			errs = append(errs, fmt.Errorf("上传文件失败：%s", err))
		}
		fileUrls = append(fileUrls, "https://" + AliYun.Host + "/" + AliYun.OssPath)
	}

	ApiResponse(c, &Response{Code: 200, Msg: "上传成功", Data: map[string]interface{}{
		"file_urls": fileUrls,
	}})
}
