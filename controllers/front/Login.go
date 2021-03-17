package front

import (
	"BeeDemo/controllers"
	"BeeDemo/models"
	"BeeDemo/utils"
	"github.com/beego/beego/v2/client/orm"
)

type LoginController struct {
	controllers.BaseController
}

//用户登录
func (c *LoginController) Post() {
	user := &models.User{
		UserName: c.GetString("username"),
		Password: utils.GetMD5(c.GetString("password")),
	}

	err := orm.NewOrm().
		QueryTable(user).
		Filter("user_name", user.UserName).
		Filter("password", user.Password).
		Filter("is_admin", 2). //普通用户登录
		One(user, []string{"id", "user_name", "cover"}...)

	if err != nil {
		c.FailMsg("用户名或密码错误")
		return
	}

	//保存用户信息到session
	err = c.SetSession("front_user", user)
	if err != nil {
		c.FailMsg("保存用户登录信息失败： session err 500")
		return
	}

	c.Success(user)
}

type FrontController struct {
	controllers.BaseController
}

func (c *FrontController) Get() {

	c.Success(utils.TestGenerateTree())
}
