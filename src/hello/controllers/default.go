package controllers

import (
	"hello/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"math"
	"path"
	"strconv"
	//"strconv"
	"time"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	/*c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"*/
	c.TplName = "register.html"
}
func (c *MainController) Post() {
	userName := c.GetString("userName")
	pwd := c.GetString("pwd")
	if userName == "" || pwd == ""{
		beego.Info("数据不能为空")
		c.Redirect("/",302)
		return
	}
	o := orm.NewOrm()
	user := models.User{}
	user.Name =userName
	user.Pwd=pwd
	_,err:=o.Insert(&user)
	if err!=nil{
		beego.Info("插入失败用户")
		c.Redirect("/regist",303)
		return
	}
	c.Redirect("/login",302)
}
func (c *MainController) ShowLogin()  {
	name :=c.Ctx.GetCookie("userName")
	c.Data["UserName"] =name

	c.TplName ="login.html"
}
/*func (c *MainController) Showlogin() {
	c.TplName="login.html"
}*/
func (c *MainController) HandLogin()  {
	userName :=c.GetString("userName")
	pwd :=c.GetString("pwd")
	if userName ==""||pwd=="" {
		beego.Info("输入有误")
	}
	var user models.User
	user.Name =userName
	o := orm.NewOrm()
	err := o.Read(&user,"Name")
	if err != nil {
		beego.Info("用户名或密码不准确")
		return
			}
	if user.Pwd != pwd {
		beego.Info("密码错误")
		return
	}
	check :=c.GetString("remember")
	if check == "on"{
		c.Ctx.SetCookie("userName",userName,time.Second*3600)
	}else {
		c.Ctx.SetCookie("userName",userName,-1)
	}
	c.SetSession("userName",userName)

	c.Redirect("/index",302)
		}
func (c *MainController) ShowIndex()  {
	userName :=c.GetSession("userName")
	if userName == nil{
		c.Redirect("/login.html",302)
		return
	}
	o := orm.NewOrm()
	id,_ := c.GetInt("select")

	var articles []models.Article
	_,err := o.QueryTable("Article").All(&articles)
	if err != nil{
		beego.Info("查询所有文章信息出错")
		return
	}
	pageIndex := c.GetString("pageIndex")
	/*id,_ := c.GetInt("select")*/
	pageIndex1,err := strconv.Atoi(pageIndex)
	if err != nil {
		pageIndex1 = 1
	}
	var articals []models.Article

	//_,err := qs.All(&articals)
	pageSize := 2

	start := pageSize*(pageIndex1-1)
	o.QueryTable("Article").Limit(pageSize,start).RelatedSel("ArticleType").Filter("ArticleType__Id",id).All(&articals)

	count,err := o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__Id",id).Count()
	//_,err :=o.QueryTable("Article").All(&articals)
	if err != nil {
		beego.Info("查询失败")
	}
	pageCount :=float64(count)/float64(pageSize)
	pageCount1 :=math.Ceil(pageCount)

	var artiTypes []models.ArticleType
	_,err = o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil{
		beego.Info("获取类型错误")
		return
	}

	c.Data["articleType"] = artiTypes
	c.Data["typeid"] = id
	c.Data["pageCount"] = pageCount1
	c.Data["count"] = count
	c.Data["articles"] = articals
	c.Data["pageIndex"] = pageIndex1
	c.TplName = "index.html"

}
func (c *MainController)ShowAdd()  {
	o := orm.NewOrm()
	var artiTypes []models.ArticleType
	_,err := o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil{
		beego.Info("获取类型错误")
		return
	}

	c.Data["articleType"] = artiTypes

	c.TplName ="add.html"
}

//func (c *MainController) HandLogin()  {
//	userName :=c.GetString("userName")
//	pwd :=c.GetString("pwd")
//	if userName == ""|| pwd ==""{
//		beego.Info("输入数据不合法")
//		c.TplName = "login.html"
//		return
//	}
//	o := orm.NewOrm()
//	user := models.User{}
//	user.Name = userName
//	user.Pwd = pwd
//
//	err := o.Read(&user,"Name")
//	if err != nil {
//		beego.Info("用户名或密码不准确")
//		return
//	}
//	c.Redirect("/index",302)
//}
func (c *MainController) HandleAdd() {
	articleName :=c.GetString("articleName")
	articleContent := c.GetString("content")
	id,err :=c.GetInt("select")
	if err !=nil{
		beego.Info("获取类型错误")
		return
	}
	//beego.Info(articleContent,articleName)
	f,h,err:=c.GetFile("uploadname")
	if articleName==""||articleContent=="" {
		beego.Info("输入有误")
	}
	fileext := path.Ext(h.Filename)
	if fileext != "jpg" {
		beego.Info("上传文件格式错误")
	}
	filename := time.Now().Format("2006-01-02 15:04:05")
	defer f.Close()
	if err != nil {
		beego.Info("上传文件失败")
		return
	}else {
		c.SaveToFile("uploadname","./static/img"+filename)
	}

	o := orm.NewOrm()
	var arti models.Article
	arti.Acontent = articleContent
	arti.ArtiName = articleName
	arti.Aimg = "/static/img"+h.Filename
	_,err = o.Insert(&arti,)
	artiType := models.ArticleType{Id:id}
	o.Read(&artiType)
	arti.ArticleType = &artiType
	if err != nil {
		beego.Info("添加文件失败")
		return
	}
	c.Redirect("/index",302)
}
func (c *MainController)ShowContent()  {
	id ,err :=c.GetInt("Id")
	if err != nil {
		beego.Info("获取文章ID错误",err)
		return
	}
	o := orm.NewOrm()
	var arti models.Article
	arti.Id =id

	err =o.Read(&arti)
	if err != nil{
		beego.Info("查询错误",err)
		return
	}
	c.Data["article"] = arti
	c.TplName = "content.html"
}
func (c *MainController)ShowUpdate()  {
	id ,err :=c.GetInt("Id")
	if err != nil {
		beego.Info("获取文章ID错误",err)
		return
	}
	o := orm.NewOrm()
	var arti models.Article
	arti.Id =id

	err =o.Read(&arti)
	if err != nil{
		beego.Info("查询错误",err)
		return
	}
	c.Data["article"] = arti
	c.TplName = "update.html"
}

func (c *MainController)HandleUpdate()  {
	id,_ := c.GetInt("id")
	artiName := c.GetString("articleName")
	content := c.GetString("content")
	f,h,err:=c.GetFile("uploadname")
	if err != nil{
		beego.Info("上传文件失败")
		return
	}else {
		defer f.Close()


		//1.要限定格式
		fileext := path.Ext(h.Filename)
		if fileext != ".jpg" && fileext != "png"{
			beego.Info("上传文件格式错误")
			return
		}
		filename := time.Now().Format("2006-01-02 15:04:05")
		defer f.Close()
		if err != nil {
			beego.Info("上传文件失败")
			return
		}else {
			c.SaveToFile("uploadname","./static/img"+filename)
		}
		if artiName == "" || content ==""{
			beego.Info("更新数据获取失败")
			return
		}

		//3.更新操作
		o := orm.NewOrm()
		arti := models.Article{Id:id}
		err = o.Read(&arti)
		if err != nil{
			beego.Info("查询数据错误")
			return
		}
		arti.ArtiName = artiName
		arti.Acontent = content
		arti.Aimg = "./static/img/"+filename


		_,err = o.Update(&arti,"ArtiName","Acontent","Aimg")
		if err != nil{
			beego.Info("更新数据显示错误")
			return
		}
		//4.返回列表页面
		c.Redirect("/index",302)
	}

}
func (c *MainController)HandleDelete(){
	Id,_ := c.GetInt("Id")
	var arti models.Article
	arti.Id =Id
	o := orm.NewOrm()
	err :=o.Read(&arti)
	if err != nil{
		beego.Info("查询错误")
		return
	}
	o.Delete(&arti)

	//3.返回列表页面
	c.Redirect("/index",302)

}
func (c *MainController)ShowAddType() {
	o := orm.NewOrm()
	var artiTypes []models.ArticleType
	_,err:=o.QueryTable("ArticleType").All(&artiTypes)
	if err!=nil{
		beego.Info("没有获取到类型数据")
	}

	c.Data["articleType"] = artiTypes
	c.TplName = "addType.html"
}
func (c*MainController)HandleAddType(){
	//1.获取内容
	typeName := c.GetString("typeName")
	//2.判断数据是否合法
	if typeName == ""{
		beego.Info("获取天津爱类型信息错误")
		return
	}
	//3.写入数据
	o := orm.NewOrm()
	artiType := models.ArticleType{}
	artiType.Tname = typeName
	_,err := o.Insert(&artiType)
	if err != nil{
		beego.Info("插入类型错误")
		return
	}
	//4.返回界面
	c.Redirect("/addType",302)
}
func (c *MainController)logout()  {
	c.DelSession("userName")
	c.Redirect("/ShowLogin",302)
}
/*func (c *MainController) ShowIndex() {
	o := orm.NewOrm()
	id,_ := c.GetInt("select")


	var articles []models.Article
	_,err := o.QueryTable("Article").All(&articles)
	if err != nil{
		beego.Info("查询所有文章信息出错")
		return
	}



	//分页处理
	//获得数据总数，总页数，当前页码
	count,err := o.QueryTable("Article").RelatedSel("ArticleType").Filter("ArticleType__Id",id).Count()
	if err != nil {
		beego.Info("查询失败",err)
		return
	}
	pagesize := int64(2)  //每页显示数据条目

	index,err := c.GetInt("pageIndex")  //当前页码
	if err != nil{
		index = 1
	}


	pageCount := math.Ceil(float64(count) / float64(pagesize))   //总页数

	if index <=0 || index > int(pageCount){
		index = 1
	}

	start := (int64(index)  -1 ) * pagesize
	// inner   left join
	var artis []models.Article
	o.QueryTable("Article").Limit(pagesize,start).RelatedSel("ArticleType").Filter("ArticleType__Id",id).All(&artis)


	//获取类型数据
	var artiTypes []models.ArticleType
	_,err = o.QueryTable("ArticleType").All(&artiTypes)
	if err != nil{
		beego.Info("获取类型错误")
		return
	}

	c.Data["articleType"] = artiTypes

	c.Data["pageCount"] = pageCount
	c.Data["count"] = count
	c.Data["articles"] = artis
	c.Data["pageIndex"] = index
	c.Data["typeid"] = id    //文章类型ID

	c.TplName = "index.html"
}
*/
