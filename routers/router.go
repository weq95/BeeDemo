package routers

import (
	"BeeDemo/controllers"
	"BeeDemo/controllers/cms"
	"BeeDemo/controllers/front"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	fmt.Println("router 注册已开始开始...")
	beego.Router("/", &controllers.MainController{})

	beego.Router("/cms/login", &cms.LoginController{}, "get:Get;post:Post")
	beego.Router("/cms/system", &cms.LoginController{})
	beego.Router("/cms/system/login", &cms.LoginController{})

	beego.Router("/cms/posts", &cms.PostController{})
	beego.Router("/cms/posts/detail", &cms.PostController{}, "get:Detail")

	beego.Router("/front/user/register", &front.RegisterController{}) //用户注册
	beego.Router("/front/user/login", &front.LoginController{})       //用户登录
	beego.Router("/front/comment", &front.CommentController{})
	beego.Router("/front/comment/post", &front.IndexController{})

	beego.Router("/trees", &front.FrontController{})

}
