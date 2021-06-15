package controllers

import (
	"blog/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"strconv"
)

// GetCategory 获取目录列表 get /v1/category
func (h *HTTPAPI) GetCategory(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string

		res      []models.Category
		result   []map[string]interface{}
		orderBy  = c.DefaultQuery("orderBy", "createdAt")
		page, _  = strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "8"))
		search   = c.DefaultQuery("search", "")
		count    int64
	)

	models.Dbms.Db.Model(&models.Category{}).Where("name LIKE ?", "%%"+search+"%%").Count(&count)
	// 直接获取所有
	if err := models.Dbms.Db.Where("name LIKE ?", "%%"+search+"%%").Preload("Articles").Order(orderBy + " desc").Offset((page - 1) * limit).Limit(limit).Find(&res).Error; err != nil {
		isOk = false
		code = 1
		message = fmt.Sprintf("category get failed: %v", err)
		h.logger.Errorf(message)
		goto exit
	} else {
		isOk = true
		result = []map[string]interface{}{}
		for _, v := range res {
			tmp := map[string]interface{}{
				"id":           v.Id,
				"name":         v.Name,
				"articles":     v.Articles,
				"articleCount": len(v.Articles),
				"createdAt":    v.CreatedAt.Format("2006-01-02 15:04"),
				"updatedAt":    v.UpdatedAt.Format("2006-01-02 15:04"),
			}
			result = append(result, tmp)
		}

	}

exit:
	if isOk {
		content = gin.H{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
				"list":  result,
				"count": count,
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

// AddCategory 新增目录 post /v1/category
func (h *HTTPAPI) AddCategory(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error
		res     models.Category
		result  map[string]interface{}
	)

	if body, e0 := ioutil.ReadAll(c.Request.Body); e0 != nil {
		message = "http post body error,err=" + e0.Error()
		h.logger.Errorf(message)
		isOk = false
		code = 1
		goto exit
	} else {
		gbody := gjson.ParseBytes(body)
		name := gbody.Get("name").String()
		s := models.Category{
			Name: name,
		}
		if name == "" {
			isOk = false
			code = 1
			message = fmt.Sprintf("name is empty!")
			h.logger.Errorf(message)
			goto exit
		}
		// models.Dbms.Db.Model(&s).Association("Tags").Clear()
		if err = models.Dbms.Db.Create(&s).Error; err != nil {
			isOk = false
			code = 1
			message = fmt.Sprintf("update category failed, %v", err)
			h.logger.Errorf(message)
			goto exit
		} else {
			if err = models.Dbms.Db.Where(models.Category{
				Name: s.Name,
			}).Last(&res).Error; err != nil {
				message = fmt.Sprintf("query category failed: [%v]", err)
				h.logger.Errorf(message)
				isOk = false
				code = 1
				goto exit
			} else {
				isOk = true
				v := res
				result = map[string]interface{}{}
				result["id"] = v.Id
				result["name"] = v.Name
			}
		}
	}
exit:
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

// EditCategory 编辑目录 put /v1/category
func (h *HTTPAPI) EditCategory(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		res    models.Category
		result map[string]interface{}
		id, _  = strconv.Atoi(c.Params.ByName("id"))
	)

	if body, e0 := ioutil.ReadAll(c.Request.Body); e0 != nil {
		message = "http post body error,err=" + e0.Error()
		h.logger.Errorf(message)
		isOk = false
		code = 1
		goto exit
	} else {
		gbody := gjson.ParseBytes(body)

		s := models.Category{
			Id:   id,
			Name: gbody.Get("name").String(),
		}

		// models.Dbms.Db.Model(&s).Association("Tags").Clear()
		if err = models.Dbms.Db.Updates(&s).Error; err != nil {
			isOk = false
			code = 1
			message = fmt.Sprintf("update category failed, %v", err)
			h.logger.Errorf(message)
			goto exit
		} else {
			if err = models.Dbms.Db.Where(models.Category{
				Id: s.Id,
			}).Last(&res).Error; err != nil {
				message = fmt.Sprintf("query category failed: [%v]", err)
				h.logger.Errorf(message)
				isOk = false
				code = 1
				goto exit
			} else {
				isOk = true
				v := res
				result = map[string]interface{}{}
				result["id"] = v.Id
				result["name"] = v.Name
			}
		}
	}

exit:
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

// DeleteCategory 删除目录 delete /v1/category
func (h *HTTPAPI) DeleteCategory(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		res   models.Category
		id, _ = strconv.Atoi(c.Params.ByName("id"))
	)

	if err = models.Dbms.Db.Where(models.Article{
		Id: id,
	}).Preload("Articles").Last(&res).Error; err != nil {
		message = fmt.Sprintf("query category failed: [%v]", err)
		h.logger.Errorf(message)
		isOk = false
		code = 1
		goto exit
	} else {
		if len(res.Articles) != 0 {
			isOk = false
			code = 1
			message = fmt.Sprintf("category %v have %v articles", res.Name, len(res.Articles))
			goto exit
		} else {
			if err = models.Dbms.Db.Where(map[string]interface{}{"id": id}).Delete(
				&models.Category{},
			).Error; err != nil {
				message = fmt.Sprintf("delete category failed, %v", err)
				h.logger.Errorf(message)
				isOk = false
				code = 1
				goto exit
			} else {
				isOk = true
			}
		}

	}

exit:
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
