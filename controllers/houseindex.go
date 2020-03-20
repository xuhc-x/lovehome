package controllers

import (
	"github.com/astaxie/beego"
	"lovehome/models"
)

type HouseIndexController struct {
	beego.Controller
}

func (c *HouseIndexController)ReData(resp map[string]interface{}){
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *HouseIndexController)GetHouseIndex(){
	resp:=make(map[string]interface{})

	resp["errno"]=models.RECODE_DATAERR
	resp["errmsg"]=models.RecodeText(models.RECODE_DATAERR)
	c.ReData(resp)
}