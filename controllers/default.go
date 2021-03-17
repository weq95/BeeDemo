package controllers

import (
	beego "github.com/beego/beego/v2/server/web"
	"strconv"
	"time"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "cms/index.tpl"
}

type BaseController struct {
	beego.Controller
}

func (c *BaseController) FailCode(code int) {
	c.ResponseData(code, "failed", nil)
}

func (c *BaseController) FailMsg(msg string) {
	c.ResponseData(400, msg, nil)
}

func (c *BaseController) FailCodeMsg(code int, msg string) {
	c.ResponseData(code, msg, nil)
}

func (c *BaseController) Fail(data interface{}) {
	c.ResponseData(400, "failed", data)
}

func (c *BaseController) SuccessCode(code int) {
	c.ResponseData(code, "success", nil)
}

func (c *BaseController) SuccessMsg(msg string) {
	c.ResponseData(200, msg, nil)
}

func (c *BaseController) SuccessCodeMsg(code int, msg string) {
	c.ResponseData(code, msg, nil)
}

func (c *BaseController) Success(data interface{}) {
	c.ResponseData(200, "success", data)
}

func (c *BaseController) ResponseData(code int, msg string, data interface{}) {
	c.Data["json"] = map[string]interface{}{
		"code":    code,
		"message": msg,
		"data":    data,
	}

	_ = c.ServeJSON()
}

//获取分页参数信息
func (c *BaseController) GetPagePer() (page, per int64) {
	page, _ = c.GetInt64("page", 1)
	if page <= int64(0) {
		page = int64(1)
	}

	per, _ = c.GetInt64("per", 15)
	if per <= int64(0) || per > int64(200) {
		per = int64(15)
	}

	return page, per
}

//上传文件
func (c *BaseController) UploadFile(img string) (string, error) {
	//上传图片
	file, h, err := c.GetFile("cover")
	if err != nil {
		//未上传文件,则使用默认图片
		return img, nil
	}
	defer func() {
		_ = file.Close()
	}()

	date := strconv.FormatInt(time.Now().Unix(), 10)
	path := "/static/upload/" + date + "_" + h.Filename
	err = c.SaveToFile("cover", "."+path)
	if err != nil {
		return "", err
	}

	return path, err
}
