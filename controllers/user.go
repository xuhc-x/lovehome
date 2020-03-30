package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/tedcy/fdfs_client"
	"lovehome/models"
	"path"
)

type UserController struct {
	beego.Controller
}

func (c *UserController)ReData(resp map[string]interface{}){
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *UserController)Reg(){
	resp:=make(map[string]interface{})
	defer c.ReData(resp)

	json.Unmarshal(c.Ctx.Input.RequestBody, &resp)  //用户注册，设置json为交互格式，获取request body 里的json数据
	fmt.Print("resp=",&resp)
	o:=orm.NewOrm()
	user:=models.User{}
	user.Password_hash=resp["password"].(string)
	user.Name=resp["name"].(string)
	user.Mobile=resp["mobile"].(string)
	id,err:=o.Insert(&user)
	if err!=nil{
		resp["errno"]=models.RECODE_NODATA
		resp["errmsg"]=models.RecodeText(models.RECODE_NODATA)
		return
	}
	fmt.Print("reg success,id=",id)
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)
	c.SetSession("name",user.Name)
	c.SetSession("userId",user.Id)
	c.SetSession("mobile",user.Mobile)
}

func (c *UserController)PostAvatar(){
	resp:=make(map[string]interface{})
	defer c.ReData(resp)
	//获取上传文件
	fileData,hd,err:=c.GetFile("avatar")
	if err!=nil{
		resp["errno"]=models.RECODE_REQERR
		resp["errmsg"]=models.RecodeText(models.RECODE_REQERR)
		fmt.Print("111")
		return
	}
	//2.得到文件后缀
	suffix:=path.Ext(hd.Filename)//专门用来取文件后缀
	//3.存储文件到fastdfs上
	fdfsClient,err:=fdfs_client.NewClientWithConfig("conf/client.conf")
	if err!=nil{
		resp["errno"]=models.RECODE_REQERR
		resp["errmsg"]=models.RecodeText(models.RECODE_REQERR)
		fmt.Print("222")
		return
	}
	fileBuffer:=make([]byte,hd.Size)
	_,err=fileData.Read(fileBuffer)
	if err!=nil{
		resp["errno"]=models.RECODE_REQERR
		resp["errmsg"]=models.RecodeText(models.RECODE_REQERR)
		fmt.Print("333")
		return
	}
	DataResponse,err:=fdfsClient.UploadByBuffer(fileBuffer,suffix[1:])
	if err!=nil{
		resp["errno"]=models.RECODE_REQERR
		resp["errmsg"]=models.RecodeText(models.RECODE_REQERR)
		fmt.Print(err.Error())
		fmt.Print("4444")
		return
	}
	//4.从session李拿到user_id
	userId:=c.GetSession("user_id")
	var user models.User
	//5.更新用户数据库中的内容
	o:=orm.NewOrm()
	qs:=o.QueryTable("user")
	qs.Filter("Id",userId).One(&user)
	user.Avatar_url=DataResponse
	_,errUpdate:=o.Update(&user)
	if errUpdate!=nil {
		resp["errno"]=models.RECODE_REQERR
		resp["errmsg"]=models.RecodeText(models.RECODE_REQERR)
		fmt.Print("55555")
		return
	}
	urlMap:=make(map[string]string)
	urlMap["avatar_url"]="192.168.159.130:8098/"+DataResponse
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)
	resp["data"]=urlMap
}
func (c *UserController)GetUserData(){
	resp:=make(map[string]interface{})
	defer c.ReData(resp)
	//从session获取userid
	userId:=c.GetSession("user_id")
	//从数据库中拿到userid 对应的user
	user:=models.User{Id:userId.(int)}
	o:=orm.NewOrm()
	err:=o.Read(&user)
	if err!=nil{
		resp["errno"]=models.RECODE_DBERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DBERR)
		return
	}
	resp["data"]=&user
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)
}


func (c *UserController)PutUpdateName() {
	resp := make(map[string]interface{})
	defer c.ReData(resp)
	//获得session中的user_id
	userId:=c.GetSession("user_id")
	//获取前端传过来的数据
	UserName:=make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody,&UserName)
	//更新user_id对应的name
	o:=orm.NewOrm()
	user:=models.User{Id:userId.(int)}
	if err:=o.Read(&user);err!=nil{
		resp["errno"]=models.RECODE_DBERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DBERR)
		return
	}

	user.Name=UserName["name"]
	if _,err :=o.Update(&user);err!=nil{
		resp["errno"]=models.RECODE_DBERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DBERR)
		return
	}
	//把session中的name自断更新
	c.SetSession("name",UserName["name"])
	//把数据打包返回给前端
	resp["data"]=UserName
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)
}

func (c *UserController)PostAuth() {
	resp := make(map[string]interface{})
	defer c.ReData(resp)

	//1.从session获得userID
	userId:=c.GetSession("user_id")
	//2.获取前端传过来的数据
	UserName:=make(map[string]string)
	json.Unmarshal(c.Ctx.Input.RequestBody,&UserName)

	//3.更新数据库中userID对应的表的信息
	o:=orm.NewOrm()
	user:=models.User{Id:userId.(int)}
	if err:=o.Read(&user);err!=nil{
		resp["errno"]=models.RECODE_DBERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DBERR)
		return
	}
	user.Real_name=UserName["real_name"]
	user.Id_card=UserName["id_card"]
	if _,err:=o.Update(&user);err!=nil{
		resp["errno"]=models.RECODE_DBERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DBERR)
		return
	}
	c.SetSession("user_id",user.Id)
	//4.打包json数据返回前端
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)

}