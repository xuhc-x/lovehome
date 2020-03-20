package controllers

import (
	"github.com/astaxie/beego"
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