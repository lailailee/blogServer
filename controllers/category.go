package controllers

import (
	"end/core"
	"end/models"
	ogorm "end/plugins/gorm"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetCategory(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string

		result []models.Category
	)

	if err := models.Dbms.Gcfg.GetAll(&result); err != nil {
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
			"result":  result,
		}
	} else {
		content = gin.H{
			"code":    code,
			"message": message,
		}
	}
	c.JSON(200, content)
}

func AddCategory(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error
		params  models.Category
		result  models.Category
	)

	if err = c.ShouldBind(&params); err != nil {
		core.Logger.Errorf("binding structure failed, %v", err)
		isOk = false
		code = 2
		message = "binding structure failed"
		goto End
	} else {
		if err = models.Dbms.Gcfg.Create(&params); err != nil {
			core.Logger.Errorf("save failed, %v", err)
			isOk = false
			code = 3
			message = "save failed"
			goto End
		} else {
			if err = models.Dbms.Gcfg.GetLast(&result); err != nil {
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

func EditCategory(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		params models.Category
		result models.Category
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
		if err = models.Dbms.Gcfg.Db.Updates(&params).Error; err != nil {
			core.Logger.Errorf("update failed, %v", err)
			isOk = false
			code = 4
			message = "update failed"
			goto End
		} else {
			if err = models.Dbms.Gcfg.Db.Where(&models.Category{Id: params.Id}).Find(&result).Error; err != nil {
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

func DeleteCategory(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		id, _ = strconv.Atoi(c.Params.ByName("id"))
	)

	if err = models.Dbms.Gcfg.Delete(
		ogorm.M{"id": id},
		&models.Category{},
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
