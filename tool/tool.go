package tool

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"os"
)

func Conf(file string) (conf map[string]string, err error) {

	conf = make(map[string]string)
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &conf)
	if err != nil {
		return
	}
	return
}

//创建图片路径和名称
func GetImgName(md5str string) (name string, err error) {
	path := "./img/qr"
	err = os.MkdirAll(path, 0777)
	if err != nil {
		return
	}
	name = path + "/" + md5str
	return
}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
