package main

import (
	"activityregister/auth"
	"activityregister/db"
	"activityregister/handle"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
)

var s handle.Send
var msg = make(chan string, 20)

func sendMsg(content string) (err error) {
	err = s.SendMessage(content)
	return
}

func main() {
	ioutil.WriteFile("server.pid", []byte(fmt.Sprintf("%d", os.Getpid())), 0666)

	go func() {
		for {
			list, err := db.GetMessage()
			if err != nil {
				log.Println(err)
			}
			for _, v := range list {
				msg <- v
			}
			time.Sleep(30000)
		}
	}()

	go func() {
		for m := range msg {
			err := s.SendMessage(m)
			if err != nil {
				log.Println(err)
			}
		}
	}()

	e := echo.New()
	e.SetHTTPErrorHandler(func(err error, c echo.Context) {
		log.Println(err)
		r := &struct {
			Ret int
			Msg string
		}{-1, err.Error()}
		c.JSON(200, r)
		return
	})
	e.POST("/login", handle.SaveUserInfo())

	api := e.Group("/api", auth.Auth)
	api.POST("/createactivtty", handle.CreateActivity())
	api.POST("/cancelactivity", handle.CancelActivity())
	api.POST("/getcreateactivity", handle.GetCreateActivity())
	api.POST("/joinactivity", handle.JoinActivity())
	api.POST("/canceljoinactivity", handle.CancelJoinActivity())
	api.POST("/getjoinactivity", handle.GetJoinActivity())
	api.POST("/getactivitybyid", handle.GetActivityByID())

	e.Run(standard.New(":7701"))
}
