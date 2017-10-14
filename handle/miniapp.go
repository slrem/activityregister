package handle

import (
	"activityregister/tool"
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	appid  string
	secret string
)

func init() {
	conf, err := tool.Conf("conf.json")
	if err != nil {
		log.Fatal(err)
	}
	appid = conf["appid"]
	secret = conf["secret"]
}

type Evidence struct {
	Openid      string `json:"openid"`
	Session_key string `json:"session_key"`
}

//获取用户信息
func GetUserInfo(code, encryptedData, iv string) (userinfo UserInfo, err error) {
	res, err := http.Get("https://api.weixin.qq.com/sns/jscode2session?appid=" +
		appid + "&secret=" + secret + "&js_code=" +
		code + "&grant_type=authorization_code")
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var e Evidence
	err = json.Unmarshal(b, &e)
	if err != nil {
		return
	}
	if e.Session_key == "" {
		err = errors.New(string(b))
		return
	}
	w := NewWXDecrypter(appid, e.Session_key)
	userinfo, err = w.Decode(encryptedData, iv)
	return
}

//发送模板消息
type Send struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	mtx         sync.Mutex
}

func (s *Send) getAccessToken() (err error) {
	res, err := http.Get("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + appid + "&secret=" + secret)
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var a Send
	err = json.Unmarshal(b, &a)
	if err != nil {
		return
	}
	if a.AccessToken == "" {
		err = errors.New(string(b))
		return
	}
	s.mtx.Lock()
	s.AccessToken = a.AccessToken
	s.ExpiresIn = a.ExpiresIn + time.Now().Unix()
	s.mtx.Unlock()
	return
}

func (s *Send) checkLive() (err error) {
	if s.ExpiresIn-time.Now().Unix() < 600 {
		err = s.getAccessToken()
	}
	return
}

type result struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func (s *Send) SendMessage(content string) (err error) {
	err = s.checkLive()
	if err != nil {
		return
	}
	surl := "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send?access_token=" + s.AccessToken
	res, err := http.Post(surl, "application/json", bytes.NewReader([]byte(content)))
	if err != nil {
		return
	}
	defer res.Body.Close()
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	var r result
	err = json.Unmarshal(b, &r)
	if err != nil {
		return
	}
	if r.Errcode != 0 {
		err = errors.New(string(b))
	}
	return
}
