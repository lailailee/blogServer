package controllers

import (
	"end/core"
	"end/models"
	ogorm "end/plugins/gorm"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func GetArticles(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		result []models.Article
	)
	if err = models.Dbms.Gcfg.Db.Preload("Category").Preload("Tags").Find(&result).Error; err != nil {
		core.Logger.Errorf("query failed: %v", err)
		isOk = false
		code = 1
		message = "query failed"
		goto End
	} else {
		isOk = true
	}

End:
	if isOk {
		content = gin.H{
			"code":    0,
			"message": "success",
			"data":    result,
		}
	} else {
		content = gin.H{
			"code":    code,
			"message": message,
		}
	}
	c.JSON(200, content)
}

func GetArticle(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		result models.Article
		id, _  = strconv.Atoi(c.Params.ByName("id"))
	)

	if err = models.Dbms.Gcfg.Db.Where(models.Article{Id: id}).Preload("Category").Preload("Tags").First(&result).Error; err != nil {
		core.Logger.Errorf("query failed: %v", err)
		isOk = false
		code = 1
		message = "query failed"
		goto End
	} else {
		isOk = true
	}

End:
	if isOk {
		content = gin.H{
			"code":    0,
			"message": "success",
			"data":    result,
		}
	} else {
		content = gin.H{
			"code":    code,
			"message": message,
		}
	}
	c.JSON(200, content)
}

func AddArticle(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error
		params  models.Article
		result  models.Article
	)

	if err = c.ShouldBind(&params); err != nil {
		core.Logger.Errorf("binding structure failed, %v", err)
		isOk = false
		code = 2
		message = "binding structure failed"
		goto End
	} else {
		if params.CreateTime == "" {
			params.CreateTime = time.Now().Format("2006-01-02 03:04:05")
		}

		if err = models.Dbms.Gcfg.Create(&params); err != nil {
			core.Logger.Errorf("save failed, %v", err)
			isOk = false
			code = 3
			message = "save failed"
			goto End
		} else {
			if err = models.Dbms.Gcfg.Db.Preload("Category").Preload("Tags").Last(&result).Error; err != nil {
				core.Logger.Errorf("query failed, %v", err)
				isOk = false
				code = 1
				message = "query failed"
				goto End
			} else {
				isOk = true
			}
		}
	}

End:
	if isOk {
		content = gin.H{
			"code":    0,
			"message": "success",
			"data":    result,
		}
	} else {
		content = gin.H{
			"code":    code,
			"message": message,
		}
	}
	c.JSON(200, content)
}

func EditArticle(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		params models.Article
		result models.Article
		id, _  = strconv.Atoi(c.Params.ByName("id"))
	)

	if err = c.ShouldBind(&params); err != nil {
		core.Logger.Errorf("binding structure failed, %v", err)
		isOk = false
		code = 2
		message = "binding structure failed"
		goto End
	} else {
		params.Id = (id)
		b := params
		models.Dbms.Gcfg.Db.Model(&params).Association("Tags").Clear()
		if err = models.Dbms.Gcfg.Db.Updates(&b).Error; err != nil {
			core.Logger.Errorf("update failed, %v", err)
			isOk = false
			code = 4
			message = "update failed"
			goto End
		} else {
			if err = models.Dbms.Gcfg.Db.Where(models.Article{
				Id: params.Id,
			}).Preload("Category").Preload("Tags").First(&result).Error; err != nil {
				core.Logger.Errorf("query failed: [%v]", err)
				isOk = false
				code = 10004
				message = "query failed"
				goto End
			} else {
				isOk = true
			}
		}
	}
End:
	if isOk {
		content = gin.H{
			"code":    0,
			"message": "success",
			"data":    result,
		}
	} else {
		content = gin.H{
			"code":    code,
			"message": message,
		}
	}
	c.JSON(200, content)
}

func DeleteArticle(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		id, _ = strconv.Atoi(c.Params.ByName("id"))
	)

	models.Dbms.Gcfg.Db.Where(&models.Article{Id: id}).Association("Tags").Clear()
	if err = models.Dbms.Gcfg.Delete(
		ogorm.M{"id": id},
		&models.Article{},
	); err != nil {
		core.Logger.Errorf("delete failed, %v", err)
		isOk = false
		code = 5
		message = "delete failed"
		goto End
	} else {
		isOk = true
	}

End:
	if isOk {
		content = gin.H{
			"code":    0,
			"message": "success",
			"result": map[string]interface{}{
				"id": id,
			},
		}
	} else {
		content = gin.H{
			"code":    code,
			"message": message,
		}
	}
	c.JSON(200, content)
}
