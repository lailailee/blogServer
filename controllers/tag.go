package controllers

import (
	"end/core"
	"end/models"
	ogorm "end/plugins/gorm"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetTag(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string

		result []models.Tag
	)

	// find all tags
	if err := models.Dbms.Gcfg.Db.Preload("Articles").Find(&result).Error; err != nil {
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

func AddTag(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error
		params  models.Tag
		result  models.Tag
	)

	if err = c.ShouldBind(&params); err != nil {
		core.Logger.Errorf("binding structure failed, %v", err)
		isOk = false
		code = 2
		message = "binding structure failed"
		goto End
	}

	if err = models.Dbms.Gcfg.Create(&params); err != nil {
		core.Logger.Errorf("save failed, %v", err)
		isOk = false
		code = 3
		message = "save failed"
		goto End
	} else {
		if err = models.Dbms.Gcfg.Db.Preload("Articles").Last(&result).Error; err != nil {
			core.Logger.Errorf("query failed, %v", err)
			isOk = false
			code = 1
			message = "query failed"
			goto End
		} else {
			isOk = true
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

func EditTag(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		params models.Tag
		result models.Tag
		id, _  = strconv.Atoi(c.Params.ByName("id"))
	)

	if err = c.ShouldBind(&params); err != nil {
		core.Logger.Errorf("binding structure failed, %v", err)
		isOk = false
		code = 2
		message = "binding structure failed"
		goto End
	} else {
		params.Id = id
		b := params
		models.Dbms.Gcfg.Db.Model(&params).Association("Articles").Clear()
		if err = models.Dbms.Gcfg.Db.Updates(&b).Error; err != nil {
			core.Logger.Errorf("update failed, %v", err)
			isOk = false
			code = 4
			message = "update failed"
			goto End
		} else {
			if err = models.Dbms.Gcfg.Db.Where(models.Article{
				Id: params.Id,
			}).Preload("Articles").First(&result).Error; err != nil {
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

func DeleteTag(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		id, _ = strconv.Atoi(c.Params.ByName("id"))
	)

	models.Dbms.Gcfg.Db.Where(&models.Tag{Id: id}).Association("Articles").Clear()
	if err = models.Dbms.Gcfg.Delete(
		ogorm.M{"id": id},
		&models.Tag{},
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
