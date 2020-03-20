package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"lovehome/models"
)

type AreaController struct {
	beego.Controller
}

func (c *AreaController)RetData(resp map[string]interface{}){
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *AreaController) GetArea() {
	fmt.Print("connect success")

	resp:=make(map[string]interface{})
	defer c.RetData(resp)
	//从session拿数据

	//从mysql数据库中拿到数据
	var areas []models.Area
	o:=orm.NewOrm()
	num,err:=o.QueryTable("area").All(&areas)
	if err!=nil {
		fmt.Print("数据错误")
		resp["errno"]=models.RECODE_DBERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DBERR)
		return
	}
	if num==0{
		fmt.Print("数据错误")
		resp["errno"]=models.RECODE_NODATA
		resp["errmsg"]=models.RecodeText(models.RECODE_NODATA)
		return
	}
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)
	resp["data"]=&areas
	//打包数据成json返回前端
	//json_str:=json.Marshal(resp)
	//c.Ctx.WriteString(json_str)
	fmt.Print("query resp=",resp)
}
