package front

import (
	"BeeDemo/controllers"
	"BeeDemo/models"
	"github.com/beego/beego/v2/client/orm"
)

type CommentController struct {
	controllers.BaseController
}

//文章评论
func (c *CommentController) Post() {
	id, err := c.GetInt("id")
	if err != nil {
		c.FailMsg("获取id失败")
		return
	}

	content := c.GetString("content")
	if len(content) <= 0 {
		c.FailMsg("请输入评论内容")
		return
	}
	pid, err := c.GetInt("pid")
	if err != nil {
		//可以不传pid 默认评论顶级
		pid = 0
	}

	user := c.GetSession("front_user")
	if user == nil {
		c.FailMsg("用户未登录，请您先登录")
		return
	}

	userInfo := user.(*models.User)
	if userInfo.Id <= 0 {
		c.FailMsg("用户登录信息异常，请重新登录")
		return
	}

	post := &models.Post{
		Id: id,
	}

	err = orm.NewOrm().
		QueryTable(post).
		Filter("id", post.Id).
		One(post)

	if err != nil {
		c.FailMsg("评论失败，未获取到文章信息")
		return
	}

	comment := &models.Comment{
		Content: content,
		Pid:     pid,
		Post:    post,
		Author:  userInfo,
	}

	_, err = orm.NewOrm().Insert(comment)
	if err != nil {
		c.FailMsg("评论失败：" + err.Error())
		return
	}

	c.Success(comment)
}
