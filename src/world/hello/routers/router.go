package routers

import (
	"hello/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{})
    beego.Router("/register",&controllers.MainController{})
    beego.Router("/login",&controllers.MainController{},"get:Showlogin;post:HandLogin")
	beego.Router("/index",&controllers.MainController{},"get:ShowIndex;")
	beego.Router("/showArticle",&controllers.MainController{},"get:ShowArticleList;post:HandleSelect")
	beego.Router("/addArticle",&controllers.MainController{},"get:ShowAdd;post:HandleAdd")
	beego.Router("/updata",&controllers.MainController{},"get:ShowUpdate;post:HandleUpdate")
	beego.Router("/content",&controllers.MainController{},"get:ShowContent")
	beego.Router("/delete",&controllers.MainController{},"get:HandleDelete")
	beego.Router("/addType",&controllers.MainController{},"get:ShowAddType;post:HandleAddType")
	beego.Router("/logout",&controllers.MainController{},"get:logout")
}
