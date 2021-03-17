package front

import (
	"BeeDemo/controllers"
	"BeeDemo/models"
	"BeeDemo/utils"
	"github.com/beego/beego/v2/client/orm"
)

type RegisterController struct {
	controllers.BaseController
}

//用户注册
func (c *RegisterController) Post() {
	pwd2 := utils.GetMD5(c.GetString("password_2"))
	user := &models.User{
		UserName: c.GetString("username"),
		Password: utils.GetMD5(c.GetString("password")),
		Cover:    "/static/upload/avatar.jpg",
		IsAdmin:  2, //普通用户
	}

	if pwd2 != user.Password {
		c.FailMsg("注册失败，两次输入密码不一致")
		return
	}

	_, err := orm.NewOrm().Insert(user)
	if err != nil {
		c.FailMsg("注册失败：用户名或已被使用" + err.Error())
		return
	}

	c.Success(user)
}
