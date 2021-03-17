package models

import (
	"github.com/beego/beego/v2/client/orm"
	"time"
)

type User struct {
	Id         int        `orm:"pk;auto"`
	UserName   string     `json:"user_name" orm:"description(用户名);index;unique"`
	Password   string     `orm:"description(密码)"`
	Cover      string     `orm:"description(头像);default(/static/upload/avatar.jpg)"`
	IsAdmin    int        `orm:"description(1管理员，2普通用户);default(2)"`
	CreateTime time.Time  `orm:"auto_now_add;type(datetime);description(创建时间)"`
	Post       []*Post    `orm:"reverse(many)"` //一对多的反向
	Comment    []*Comment `orm:"reverse(many)"` //一对多的反向
}

//指定表名
func (u *User) TableName() string {
	return "user"
}

func init() {
	//注册模型
	orm.RegisterModel(
		new(User),
		new(Post),
		new(Comment),
	)
}
