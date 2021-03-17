package models

import (
	"time"
)

type Post struct {
	Id         int        `orm:"pk;auto"`
	Title      string     `orm:"description(标题)"`
	Desc       string     `orm:"description(描述)"`
	Content    string     `orm:"description(内容);size(400)"`
	Cover      string     `orm:"description(封面);default(/static/upload/no_cover.png)"`
	ReadNum    int        `orm:"description(阅读数);default(0)"`
	StarNum    int        `orm:"description(点赞数);default(0)"`
	Author     *User      `orm:"description(作者);rel(fk)"` //一对多的正向
	Comments   []*Comment `orm:"description(文章评论);reverse(many)"`
	CreateTime time.Time  `orm:"auto_now_add;type(datetime)"`
}

func (p *Post) TableName() string {
	return "posts"
}
