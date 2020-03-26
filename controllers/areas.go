package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/astaxie/beego/orm"
	"lovehome/models"
	"time"
)

type AreaController struct {
	beego.Controller
}

func (c *AreaController)RetData(resp map[string]interface{}){
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *AreaController) GetArea() {
	//fmt.Print("connect success")
	resp:=make(map[string]interface{})
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)
	defer c.RetData(resp)
	//从redis缓存中拿数据
	cacheCon, err := cache.NewCache("redis", `{"key":"lovehome","conn":":6399","dbNum":"0"}`)
	if areaData:=cacheCon.Get("area");areaData!=nil{
		resp["data"]=areaData
		fmt.Print("get data from cache===",resp["data"])
		return
	}

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
	//resp["errno"]=models.RECODE_OK
	//resp["errmsg"]=models.RecodeText(models.RECODE_OK)

	resp["data"]=areas
	//把数据转换成json格式存入缓存
	jsonStr,err:=json.Marshal(areas)
	fmt.Print("jsonStr====",jsonStr)
	if err !=nil{
		fmt.Print("encoding err")
		return
	}
	cacheCon.Put("area",jsonStr, 3600*time.Second)
	//打包数据成json返回前端
	//json_str:=json.Marshal(resp)
	//c.Ctx.WriteString(json_str)
	fmt.Print("query resp=",resp)
}
