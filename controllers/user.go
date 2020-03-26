package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"lovehome/models"
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
}