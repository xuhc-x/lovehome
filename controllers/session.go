package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"lovehome/models"
)

type SessionController struct {
	beego.Controller
}

func (c *SessionController)ReData(resp map[string]interface{}){
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *SessionController)GetSessionData(){
	resp:=make(map[string]interface{})
	defer c.ReData(resp)
	user:=models.User{}
	//user.Name="wyj"
	//resp["errno"]=models.RECODE_OK
	//resp["errmsg"]=models.RecodeText(models.RECODE_OK)
	//resp["data"]=user

	resp["errno"]=models.RECODE_DBERR
	resp["errmsg"]=models.RecodeText(models.RECODE_DBERR)

	name:=c.GetSession("name")
	if name!=nil{
		user.Name=name.(string)
		resp["errno"]=models.RECODE_OK
		resp["errmsg"]=models.RecodeText(models.RECODE_OK)
		resp["data"]=user
	}
}

func (c *SessionController)DeleteSessionData(){
	resp:=make(map[string]interface{})
	defer c.ReData(resp)
	c.DelSession("name")
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)
}

func (c *SessionController)Login(){
	//得到用户信息
	resp:=make(map[string]interface{})
	defer c.ReData(resp)

	_ = json.Unmarshal(c.Ctx.Input.RequestBody, &resp) //用户注册，设置json为交互格式，获取request body 里的json数据
	beego.Info("=====name=",resp["mobile"],"======password=",resp["password"])
	//判断是否合法
	if resp["mobile"]==nil||resp["password"]==nil {
		resp["errno"]=models.RECODE_DATAERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DATAERR)
		fmt.Print("111")
		return
	}

	//与数据库匹配账号密码
	o:=orm.NewOrm()
	user:=models.User{Name:resp["mobile"].(string)}

	qs:=o.QueryTable("user")
	err:=qs.Filter("mobile",resp["mobile"].(string)).One(&user)
	if err!= nil {
		resp["errno"]=models.RECODE_DATAERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DATAERR)
		fmt.Print("222")
		return
	}
	if user.Password_hash!=resp["password"] {
		resp["errno"]=models.RECODE_DATAERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DATAERR)
		fmt.Print("333")
		return
	}
	//添加session

	c.SetSession("name",user.Name)
	c.SetSession("mobile",resp["mobile"])
	c.SetSession("user_id",user.Id)
	//返回json数据给前端
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)


}
