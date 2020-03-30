package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"lovehome/models"
	"strconv"
)

type HouseController struct {
	beego.Controller
}

func (c *HouseController)ReData(resp map[string]interface{}){
	c.Data["json"] = resp
	c.ServeJSON()
}

func (c *HouseController)GetHouseData() {
	resp:= make(map[string]interface{})
	defer c.ReData(resp)
	//从session获取userid
	userId:=c.GetSession("user_id")
	//从数据库中拿到userid 对应的user
	houses:=[]models.House{}
	o:=orm.NewOrm()
	qs:=o.QueryTable("house")
	num,err:=qs.Filter("user_id",userId.(int)).All(&houses)//user__id是双下划线
	if err!=nil{
		resp["errno"]=models.RECODE_DBERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DBERR)
		return
	}
	if num==0 {
		resp["errno"]=models.RECODE_NODATA
		resp["errmsg"]=models.RecodeText(models.RECODE_NODATA)
		return
	}
	respData:=make(map[string]interface{})
	respData["houses"]=houses
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)

}

func (c *HouseController)PostHouseData() {
	resp := make(map[string]interface{})
	defer c.ReData(resp)
	//1.从前端拿到数据
	reqData:=make(map[string]interface{})
	json.Unmarshal(c.Ctx.Input.RequestBody,&reqData)
	//2.判断前端数据的合法性

	//3.插入数据到数据库
	house:=models.House{}
	house.Title=reqData["title"].(string)
	price,_:=strconv.Atoi(reqData["price"].(string))
	house.Price=price
	house.Address=reqData["address"].(string)
	roomCount,_:=strconv.Atoi(reqData["room_count"].(string))
	house.Room_count=roomCount
	house.Unit=reqData["unit"].(string)
	house.Beds=reqData["beds"].(string)
	minDay,_:=strconv.Atoi(reqData["min_days"].(string))
	maxDay,_:=strconv.Atoi(reqData["max_days"].(string))
	house.Min_days=minDay
	house.Max_days=maxDay
	facilitys:=[]models.Facility{}
	for _,fid:=range reqData["facility"].([]interface{}){
		fId,_:=strconv.Atoi(fid.(string))
		fac:=models.Facility{Id:fId}
		facilitys=append(facilitys,fac)
	}
	areaId,_:=strconv.Atoi(reqData["area_id"].(string))
	area:=models.Area{Id:areaId}
	userId:=c.GetSession("user_id").(int)
	user:=models.User{Id:userId}
	house.User=&user
	house.Area=&area
	o:=orm.NewOrm()
	_,err:=o.Insert(house)
	if err!=nil{
		resp["errno"]=models.RECODE_DBERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DBERR)
		return
	}
	//house.Id=int(houseId)

	m2m:=o.QueryM2M(&house,"Facilities")
	num,errM2M:=m2m.Add(&facilitys)
	if errM2M !=nil||num==0{
		resp["errno"]=models.RECODE_DBERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DBERR)
		return
	}
	respData:=make(map[string]interface{})
	respData["house_id"]=house.Id
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)

}

func (c *HouseController)GetDetailHouseData() {
	resp := make(map[string]interface{})
	defer c.ReData(resp)

	respData:=make(map[string]interface{})
	//1.获取当前用户的user_id
	userId:=c.GetSession("user_id")
	//2.从url中获取房屋id
	houseId:=c.Ctx.Input.Param(":id")
	hId,err:=strconv.Atoi(houseId)
	if err!=nil{
		resp["errno"]=models.RECODE_REQERR
		resp["errmsg"]=models.RecodeText(models.RECODE_REQERR)
		return
	}
	//3.从缓存中获取当前房屋数据 redis

	//4.关联查询
	o:=orm.NewOrm()
	house:=models.House{Id:hId}
	user:=models.User{Id:userId.(int)}
	house.User=&user
	if err:=o.Read(&house);err!=nil{
		resp["errno"]=models.RECODE_DBERR
		resp["errmsg"]=models.RecodeText(models.RECODE_DBERR)
		return
	}
	o.LoadRelated(&house,"Area")
	o.LoadRelated(&house,"User")
	o.LoadRelated(&house,"Images")
	o.LoadRelated(&house,"Facilities")

	fmt.Print("house")
	respData["house"]=house
	resp["data"]=respData
	resp["errno"]=models.RECODE_OK
	resp["errmsg"]=models.RecodeText(models.RECODE_OK)
	o.QueryTable("")
	//5.存入缓存

	//6。打包返回json
}