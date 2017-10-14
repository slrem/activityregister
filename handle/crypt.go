package handle

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"errors"
)

var (
	DecodeBase64Error     = errors.New("DecodeBase64Error")
	IVDecodeBase64Error   = errors.New("IVDecodeBase64Error")
	DataDecodeBase64Error = errors.New("DataDecodeBase64Error")

	DecryptAESError = errors.New("DecryptAESError")
	JSONError       = errors.New("JSONError")
	CheckAppIDError = errors.New("CheckAppIDError")
)

type PKCS7Encoder interface {
	Encode([]byte) []byte
	Decode([]byte) []byte
}

const (
	blockSize = 32
)

type pkcs7Encoder struct {
}

func (p pkcs7Encoder) Encode(src []byte) (dist []byte) {
	byteLen := len(src)
	pad := blockSize - (byteLen % blockSize)
	if pad == 0 {
		pad = blockSize
	}

	b := bytes.Repeat([]byte{byte(pad)}, pad)
	dist = append(src, b...)
	return
}

func (p pkcs7Encoder) Decode(src []byte) (dist []byte) {
	byteLen := len(src)
	pad := int(src[byteLen-1])
	if pad < 1 || pad > 32 {
		dist = src
	} else {
		dist = src[:byteLen-pad]
	}
	return
}

type crypter struct {
	key  []byte
	pkcs PKCS7Encoder
}

func newCrypter(k string) (p *crypter, err error) {
	b, err := base64.StdEncoding.DecodeString(k)
	if err != nil {
		err = DecodeBase64Error
		return
	}

	p = &crypter{
		key:  b,
		pkcs: pkcs7Encoder{},
	}
	return
}

func (c crypter) Decode(data []byte, iv []byte) (ret []byte, err error) {
	b, err := aes.NewCipher(c.key)
	if err != nil {
		err = DecryptAESError
		return
	}

	cp := cipher.NewCBCDecrypter(b, iv)
	cp.CryptBlocks(data, data)

	ret = c.pkcs.Decode(data)
	return
}

type UserInfo struct {
	OpenID    string `json:"openId"`
	NickName  string `json:"nickName"`
	Gender    int    `json:"gender"`
	City      string `json:"city"`
	Province  string `json:"province"`
	Country   string `json:"country"`
	AvatarURL string `json:"avatarUrl"`
	UnionID   string `json:"unionId"`
	Watermark struct {
		AppID     string `json:"appid"`
		Timestamp int    `json:"timestamp"`
	} `json:"watermark"`
}

type WXDecrypter struct {
	appID string
	key   string
}

func NewWXDecrypter(appID, key string) *WXDecrypter {
	return &WXDecrypter{
		appID: appID,
		key:   key,
	}
}

func (w *WXDecrypter) Decode(data, iv string) (userInfo UserInfo, err error) {
	biv, err := base64.StdEncoding.DecodeString(iv)
	if err != nil {
		err = IVDecodeBase64Error
		return
	}

	content, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		err = DataDecodeBase64Error
		return
	}

	c, err := newCrypter(w.key)
	if err != nil {
		return
	}

	content, err = c.Decode(content, biv)
	if err != nil {
		return
	}

	err = json.Unmarshal(content, &userInfo)
	if err != nil {
		return
	}

	if userInfo.Watermark.AppID != w.appID {
		err = CheckAppIDError
		return
	}

	return
}
