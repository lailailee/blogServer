package controllers

import (
	"blog/core"
	"blog/models"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"time"
)

func (h *HTTPAPI) Login(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		rsp     interface{}
	)

	if body, e0 := ioutil.ReadAll(c.Request.Body); e0 != nil {
		code = 1 // 读取参数出错
		message = "http post body error,err=" + e0.Error()
		h.logger.Errorf(message)
		goto exit
	} else {
		gbody := gjson.ParseBytes(body)
		name := gbody.Get("name")
		password := gbody.Get("password")
		if !name.Exists() {
			h.logger.Errorf("user http post miss key=[name] error")
			code = 1 // 参数不存在
			message = "user http post miss key=[name] error"
			goto exit
		}
		if !password.Exists() {
			h.logger.Errorf("user http post miss key=[password] error")
			code = 1 // 参数不存在
			message = "user http post miss key=[password] error"
			goto exit
		}

		var muser models.User
		if err := models.Dbms.Db.Where(models.User{Name: name.String(), Password: password.String()}).First(&muser).Error; err != nil {
			code = 1 // 参数不存在
			message = "user name or password error"
			h.logger.Errorf("get all user records error,error=%v", err)
			goto exit
		} else {
			isOk = true
			code = 0
			message = "success"
			token := jwt.New(jwt.SigningMethodHS256)
			claims := make(jwt.MapClaims)

			// 5小时过期
			claims["exp"] = time.Now().Add(time.Hour * time.Duration(5)).Unix()
			claims["iat"] = time.Now().Unix()
			claims["name"] = muser.Name
			token.Claims = claims
			tokenString, err := token.SignedString([]byte(core.SecretKey))
			if err != nil {
				c.Status(500)
				return
			}
			rsp = gin.H{
				"id":    muser.Id,
				"name":  muser.Name,
				"token": tokenString,
			}
			isOk = true
			goto exit
		}
	}

exit:
	if isOk {
		content = gin.H{
			"code":   code,
			"result": rsp,
		}
	} else {
		content = gin.H{
			"code":    code,
			"message": message,
		}
	}
	c.JSON(200, content)
}

func (h *HTTPAPI) CreateAccount(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string

		params models.User
	)
	if err := c.ShouldBind(&params); err != nil {
		h.logger.Errorf("data bind failed，v%", err)
		isOk = false
		message = "data bind failed"
		code = 1
		goto exit
	} else {
		var b models.User
		if err = models.Dbms.Db.Where(&models.User{
			Name: params.Name,
		}).First(&b).Error; err != nil {
			if err = models.Dbms.Db.Create(&params).Error; err != nil {
				h.logger.Errorf("save data fail: [%v]", err)
				isOk = false
				code = 1
				message = "save data fail"
				goto exit
			} else {
				isOk = true
				code = 0
				message = "success"
				goto exit
			}
		} else {
			h.logger.Errorf("name repeat")
			isOk = false
			code = 10003
			message = fmt.Sprintf("name %v repeat", params.Name)
			goto exit
		}
	}

exit:
	if isOk {
		content = gin.H{
			"code":    code,
			"message": message,
			// "result":  result,
		}
	} else {
		content = gin.H{
			"code":    code,
			"message": message,
		}
	}
	c.JSON(200, content)
}
