package front

import (
	"BeeDemo/controllers"
	"BeeDemo/models"
	"fmt"
	"github.com/beego/beego/v2/client/orm"
)

type IndexController struct {
	controllers.BaseController
}

//评论获取文章详情，获取登录用户信息
func (c *IndexController) Get() {
	id, _ := c.GetInt("id")
	post := &models.Post{
		Id: id,
	}

	err := orm.NewOrm().
		QueryTable(post).
		Filter("id", id).
		RelatedSel().
		One(post)

	if err != nil {
		c.FailMsg("获取文章详情失败：" + err.Error())
		return
	}
	fmt.Println("获取文章详情成功了哦 ===============================")

	comments := make([]models.Comment, 0)
	_, err = orm.NewOrm().
		QueryTable(new(models.Comment)).
		Filter("post_id", post.Id).
		Filter("pid", 0).
		All(&comments)
	if err != nil {
		c.FailMsg("获取文章评论列表信息失败")
		return
	}

	//文章阅读数量+1
	go func() {
		_, _ = orm.NewOrm().
			QueryTable(post).
			Filter("id", post.Id).
			Update(map[string]interface{}{
				"read_num": post.ReadNum + 1,
			})
	}()

	//递归分类
	commentTree := make([]*models.CommentTree, 0)
	for _, comment := range comments {

		result := &models.CommentTree{
			Id:         comment.Id,
			Pid:        comment.Pid,
			Content:    comment.Content,
			CreateTime: comment.CreateTime,
			Children:   []*models.CommentTree{},
		}

		commentChildren(comment.Id, result)
		commentTree = append(commentTree, result)

	}

	//获取用户登录信息
	user := c.GetSession("front_user")
	c.Success(map[string]interface{}{
		"post":        post,
		"user":        user,
		"commentTree": commentTree,
	})
}

//递归获取顶级评论的下级评论内容
func commentChildren(pid int, tree *models.CommentTree) {
	comments := []models.Comment{}
	_, err := orm.NewOrm().
		QueryTable(new(models.Comment)).
		Filter("pid", pid).
		RelatedSel().
		All(&comments)

	//查询出错，或没有下级评级，直接返回
	if err != nil {
		return
	}

	//多级评论
	for _, val := range comments {
		fmt.Println(val)

		child := &models.CommentTree{
			Id:         val.Id,
			Pid:        val.Pid,
			Content:    val.Content,
			CreateTime: val.CreateTime,
			Children:   []*models.CommentTree{},
		}

		tree.Children = append(tree.Children, child)
		commentChildren(val.Id, child)
	}

	return
}
