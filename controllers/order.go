package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"lovehome/models"
)

type OrderController struct {
	beego.Controller
}

func (c *OrderController)ReData(resp map[string]interface{}){
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *OrderController)GetOrderData() {
	resp := make(map[string]interface{})
	defer c.ReData(resp)

	//1.根据session获得当前用户ID
	userId:=c.GetSession("user_id")

	//2.根据url获得当前操作的角色
	role:=c.GetString("role")
	if role=="custom" {
		//orders:=[]OrderController{}
		orders:=[]models.OrderHouse{}
		o:=orm.NewOrm()
		qs:=o.QueryTable("OrderHouse")
		user:=models.User{Id:userId.(int)}
		qs.Filter("user__id",userId.(int)).All(&orders)
		for _,order:=range orders{
			order.User=&user
			o.LoadRelated(&order,"User")
		}
		respData:=make(map[string]interface{})
		respData["order"]=orders
		resp["data"]=respData
		resp["errno"]=models.RECODE_OK
		resp["errmsg"]=models.RecodeText(models.RECODE_OK)
		return

	}
	if role=="landlord" {

	}
	if role=="" {
		resp["errno"]=models.RECODE_REQERR
		resp["errmsg"]="请求url错误"
		return
	}
}