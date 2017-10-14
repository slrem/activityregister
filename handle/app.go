package handle

import (
	"activityregister/db"
	"activityregister/tool"
	"fmt"
	"log"
	"strconv"

	"github.com/labstack/echo"
)

type resp struct {
	Ret  int
	Msg  string
	Data interface{}
}

//保存用户信息
type Ticket struct {
	Uid     int
	Session string
}

func SaveUserInfo() echo.HandlerFunc {
	return func(c echo.Context) error {
		code := c.FormValue("code")
		encryptedData := c.FormValue("encryptedData")
		iv := c.FormValue("iv")
		a, err := GetUserInfo(code, encryptedData, iv)
		if err != nil {
			log.Println(err)
			r := &resp{Ret: -1, Msg: err.Error()}
			return c.JSON(200, r)
		}
		uid, err := db.SaveUserInfo(a.OpenID, a.NickName, a.Gender, a.City, a.Province,
			a.Country, a.AvatarURL, a.UnionID)
		if err != nil {
			log.Println(err)
			r := &resp{Ret: -2, Msg: err.Error()}
			return c.JSON(200, r)
		}
		var t Ticket
		t.Uid = uid
		t.Session = tool.MD5("activityregister" + fmt.Sprint(uid) + "activityregister")
		r := &resp{Ret: 0, Msg: "成功", Data: t}
		return c.JSON(200, r)
	}
}

//创建活动
func CreateActivity() echo.HandlerFunc {
	return func(c echo.Context) error {

		title := c.FormValue("title")
		descriptiom := c.FormValue("descriptiom")
		act_start_time := c.FormValue("act_start_time")
		act_end_time := c.FormValue("act_end_time")
		join_start_time := c.FormValue("join_start_time")
		join_end_time := c.FormValue("join_end_time")
		originator, _ := strconv.Atoi(c.FormValue("uid"))

		err := db.CreateActivity(title, descriptiom, act_start_time, act_end_time, join_start_time, join_end_time, originator)
		if err != nil {
			log.Println(err)
			r := &resp{Ret: -1, Msg: err.Error()}
			return c.JSON(200, r)
		}
		r := &resp{Ret: 0, Msg: "成功"}
		return c.JSON(200, r)
	}
}

//发起人取消活动
func CancelActivity() echo.HandlerFunc {
	return func(c echo.Context) error {
		//activityId, userId
		activityid, _ := strconv.Atoi(c.FormValue("activityid"))
		uid, _ := strconv.Atoi(c.FormValue("uid"))
		err := db.CancelActivity(activityid, uid)
		if err != nil {
			log.Println(err)
			r := &resp{Ret: -1, Msg: err.Error()}
			return c.JSON(200, r)
		}
		r := &resp{Ret: 0, Msg: "成功"}
		return c.JSON(200, r)
	}
}

//查看自己发起的活动
func GetCreateActivity() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid, _ := strconv.Atoi(c.FormValue("uid"))
		list, err := db.GetCreateActivity(uid)
		if err != nil {
			log.Println(err)
			r := &resp{Ret: -1, Msg: err.Error()}
			return c.JSON(200, r)
		}
		r := &resp{Ret: 0, Msg: "成功", Data: list}
		return c.JSON(200, r)
	}
}

//报名参加活动
func JoinActivity() echo.HandlerFunc {
	return func(c echo.Context) error {

		activityid, _ := strconv.Atoi(c.FormValue("activityid"))
		uid, _ := strconv.Atoi(c.FormValue("uid"))
		formid := c.FormValue("formid")
		err := db.JoinActivity(activityid, uid, formid)
		if err != nil {
			log.Println(err)
			r := &resp{Ret: -1, Msg: err.Error()}
			return c.JSON(200, r)
		}
		r := &resp{Ret: 0, Msg: "成功"}
		return c.JSON(200, r)
	}
}

//取消报名
func CancelJoinActivity() echo.HandlerFunc {
	return func(c echo.Context) error {
		activityid, _ := strconv.Atoi(c.FormValue("activityid"))
		uid, _ := strconv.Atoi(c.FormValue("uid"))
		err := db.CancelJoinActivity(activityid, uid)
		if err != nil {
			log.Println(err)
			r := &resp{Ret: -1, Msg: err.Error()}
			return c.JSON(200, r)
		}
		r := &resp{Ret: 0, Msg: "成功"}
		return c.JSON(200, r)
	}
}

//查看自己报名的活动
func GetJoinActivity() echo.HandlerFunc {
	return func(c echo.Context) error {
		uid, _ := strconv.Atoi(c.FormValue("uid"))
		list, err := db.GetJoinActivity(uid)
		if err != nil {
			log.Println(err)
			r := &resp{Ret: -1, Msg: err.Error()}
			return c.JSON(200, r)
		}
		r := &resp{Ret: 0, Msg: "成功", Data: list}
		return c.JSON(200, r)
	}
}

//查看活动具体内容
func GetActivityByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		activityid, err := strconv.Atoi(c.FormValue("activityid"))
		if err != nil {
			r := &resp{Ret: -2, Msg: "activityid需是数字"}
			return c.JSON(200, r)
		}
		uid, _ := strconv.Atoi(c.FormValue("uid"))
		list, err := db.GetActivityByID(activityid, uid)
		if err != nil {
			log.Println(err)
			r := &resp{Ret: -1, Msg: err.Error()}
			return c.JSON(200, r)
		}
		r := &resp{Ret: 0, Msg: "成功", Data: list}
		return c.JSON(200, r)
	}
}
