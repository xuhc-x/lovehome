package routers

import (
	"github.com/astaxie/beego"
	"lovehome/controllers"
)

func init() {
    beego.Router("/", &controllers.MainController{})
	beego.Router("/api/v1.0/areas", &controllers.AreaController{},"get:GetArea")
	beego.Router("/api/v1.0/houses/index", &controllers.HouseIndexController{},"get:GetHouseIndex")
	beego.Router("/api/v1.0/session", &controllers.SessionController{},"get:GetSessionData;delete:DeleteSessionData")
	beego.Router("/api/v1.0/users", &controllers.UserController{},"post:Reg")
	beego.Router("/api/v1.0/sessions", &controllers.SessionController{},"post:Login")
	beego.Router("/api/v1.0/user/avatar", &controllers.UserController{},"post:PostAvatar")
	beego.Router("/api/v1.0/user", &controllers.UserController{},"get:GetUserData")
	beego.Router("/api/v1.0/user/name", &controllers.UserController{},"put:PutUpdateName")
	beego.Router("/api/v1.0/user/auth", &controllers.UserController{},"get:GetUserData;post:PostAuth")
	beego.Router("/api/v1.0/user/houses", &controllers.HouseController{},"get:GetHouseData")
	beego.Router("/api/v1.0/houses", &controllers.HouseController{},"post:PostHouseData")
	beego.Router("/api/v1.0/houses/?:id", &controllers.HouseController{},"get:GetDetailHouseData")
	//v1.0/user/orders?role=custom
	beego.Router("/api/v1.0/user/orders", &controllers.OrderController{},"get:GetOrderData")
}
