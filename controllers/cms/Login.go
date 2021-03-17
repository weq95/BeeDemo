package cms

import (
	"BeeDemo/controllers"
	"BeeDemo/models"
	"BeeDemo/utils"
	"github.com/beego/beego/v2/client/orm"
	"net/http"
	"time"
)

type LoginController struct {
	controllers.BaseController
}

func (c *LoginController) Get() {
	c.Data["json"] = map[string]interface{}{
		"code":    http.StatusOK,
		"message": "欢饮来到登录界面",
		"data": map[string]interface{}{
			"title":    "SystemManager 管理系统",
			"date":     time.Now().Format("2006-01-02 15:04:05"),
			"username": "admin",
			"password": utils.GetMD5("123456"),
		},
		"session": c.GetSession("cms_user"),
	}

	_ = c.ServeJSON()
}

func (c *LoginController) Post() {
	user := &models.User{
		UserName: utils.FilterStr(c.GetString("username")),
		Password: utils.GetMD5(c.GetString("password")),
	}

	o := orm.NewOrm().QueryTable(user)
	err := o.Filter("user_name", user.UserName).Filter("password", user.Password).One(user, []string{"id", "user_name", "cover", "create_time", "is_admin"}...)

	if err != nil || user.Id <= 0 {
		c.FailMsg("登录失败，用户名或密码错误...")
		return
	}

	//redis 原生存储方式
	//reply, err := utils.GetRedis().Get().Do("SET", "cms_user", string(bytes))
	err = c.SetSession("cms_user", user)
	if err != nil {
		c.FailMsg("登录失败，请稍后重试， session err 500")
		return
	}

	c.Success(map[string]interface{}{
		"user":    user,
		"session": c.GetSession("cms_user"),
	})
}
