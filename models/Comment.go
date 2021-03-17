package models

import (
	"time"
)

type Comment struct {
	Id         int       `orm:"pk;auto"`
	Author     *User     `orm:"rel(fk);description(评论人)"`
	Pid        int       `orm:"description(父级评论);default(0)"`
	Post       *Post     `orm:"description(文章id);rel(fk)"`
	Content    string    `orm:"size(255);description(评论内容)"`
	CreateTime time.Time `orm:"auto_now_add;type(datetime)"`
}

func (c Comment) TableName() string {
	return "post_comments"
}

//多级评论数据结构体
type CommentTree struct {
	Id         int            `json:"id"`
	Pid        int            `json:"pid"`
	Content    string         `json:"content"`
	CreateTime time.Time      `json:"create_time"`
	Author     *User          `json:"author"`
	Children   []*CommentTree `json:"children"`
}

