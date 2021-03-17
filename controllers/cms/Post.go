package cms

import (
	"BeeDemo/controllers"
	"BeeDemo/models"
	"BeeDemo/utils"
	"errors"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

type PostController struct {
	controllers.BaseController
}

//文章列表
func (c *PostController) Get() {
	posts := make([]*models.Post, 0)
	page, per := c.GetPagePer()

	qs := orm.NewOrm().QueryTable(new(models.Post))

	_, err := qs.RelatedSel().All(&posts,
		[]string{
			"id", "title", "desc", "cover",
			"read_num", "star_num", "author_id",
			"create_time", "content"}...)
	if err != nil {
		c.ResponseData(400, "获取文章失败", map[string]interface{}{
			"posts":      posts,
			"pagination": utils.CreatePaging(page, per, int64(0)),
		})
		return
	}

	count, _ := qs.Count()
	c.Success(map[string]interface{}{
		"posts":      posts,
		"pagination": utils.CreatePaging(page, per, count),
	})
}

//文章详情
func (c *PostController) Detail() {
	id, err := c.GetInt("id")
	if err != nil {
		c.FailMsg("获取文章id失败")
		return
	}

	post := &models.Post{
		Id: id,
	}

	err = orm.NewOrm().QueryTable(post).
		Filter("id", post.Id).One(post, []string{
		"id", "title", "desc", "cover", "read_num",
		"star_num", "author_id", "create_time", "content",
	}...)
	if err != nil {
		c.FailMsg("获取文章详情失败，请稍后重试...")
		return
	}

	c.Success(post)
}

//文章添加
func (c *PostController) Post() {
	post, err := c.GetPostParam()
	if err != nil {
		c.FailMsg(err.Error())
		return
	}

	id, err := orm.NewOrm().Insert(post)
	if err != nil {
		c.FailMsg("文章添加失败：" + err.Error())
		return
	}

	post.Id = int(id)
	c.Success(post)
}

//文章编辑
func (c *PostController) Patch() {
	id, err := c.GetInt("id")
	if err != nil || id <= 0 {
		fmt.Println("id：", id)
		c.FailMsg("获取文章id失败")
		return
	}

	post, err := c.GetPostParam()
	if err != nil {
		c.FailMsg(err.Error())
		return
	}

	//开始执行更新
	num, err := orm.NewOrm().
		QueryTable(new(models.Post)).
		Filter("id", id).
		Update(map[string]interface{}{
			"title": post.Title,
			"desc":  post.Desc,
			"cover": post.Cover,
		})
	if err != nil {
		c.FailMsg("文章更新失败：" + err.Error())
		return
	}

	if num <= 0 {
		c.SuccessMsg("文章更新成功，影响行数 0 条")
		return
	}
	c.Success(post)
}

//添加,编辑文章 获取参数
func (c *PostController) GetPostParam() (*models.Post, error) {
	imgPath, err := c.UploadFile("/static/upload/no_cover.png")

	if err != nil {
		c.FailMsg(err.Error())
		return nil, err
	}

	//获取session中用户信息
	userInfo := c.GetSession("cms_user")
	if userInfo == nil {
		return nil, errors.New("用户未登录")
	}

	//进行类型断言
	user := userInfo.(*models.User)
	if user.Id <= 0 {
		return nil, errors.New("用户登录信息异常，请重新登录")
	}

	post := &models.Post{
		Title:   c.GetString("title"),
		Desc:    c.GetString("desc"),
		Content: c.GetString("content"),
		Cover:   imgPath,
		Author:  user,
	}

	if len(post.Title) <= 0 ||
		len(post.Desc) <= 0 ||
		len(post.Content) <= 0 {

		return nil, errors.New("添加文章失败，缺少必传参数")
	}

	return post, nil
}

//文章删除
func (c *PostController) Delete() {
	id, err := c.GetInt("id")
	if err != nil {
		c.FailMsg("获取文章id失败")
		return
	}

	post := &models.Post{
		Id: id,
	}

	num, err := orm.NewOrm().
		QueryTable(post).
		Filter("id", post.Id).
		Delete()

	if err != nil {
		c.FailMsg("删除失败：" + err.Error())
		return
	}

	if num <= 0 {
		c.SuccessMsg("删除操作成功，影响行数0行")
		return
	}

	c.SuccessMsg("删除成功， succeed")
}
