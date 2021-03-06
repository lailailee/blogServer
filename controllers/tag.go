package controllers

import (
	"blog/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"strconv"
)

// GetTag 获取标签列表 get /v1/tag
func (h *HTTPAPI) GetTag(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string

		res      []models.Tag
		result   []map[string]interface{}
		orderBy  = c.DefaultQuery("orderBy", "createdAt")
		page, _  = strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "8"))
		search   = c.DefaultQuery("search", "")
		count    int64
	)
	// 直接获取所有
	models.Dbms.Db.Model(&models.Tag{}).Where("name LIKE ?", "%%"+search+"%%").Count(&count)
	if err := models.Dbms.Db.Where("name LIKE ?", "%%"+search+"%%").Preload("Articles").Order(orderBy + " desc").Offset((page - 1) * limit).Limit(limit).Find(&res).Error; err != nil {
		isOk = false
		code = 1
		message = fmt.Sprintf("tag get failed: %v", err)
		h.logger.Errorf(message)
		goto exit
	} else {
		isOk = true
		result = []map[string]interface{}{}
		for _, v := range res {
			tmp := map[string]interface{}{
				"id":           v.Id,
				"name":         v.Name,
				"articleCount": len(v.Articles),
				"article":      v.Articles,
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

// AddTag 新增标签 post /v1/tag
func (h *HTTPAPI) AddTag(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error
		res     models.Tag
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
		s := models.Tag{
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
			message = fmt.Sprintf("update tag failed, %v", err)
			h.logger.Errorf(message)
			goto exit
		} else {
			if err = models.Dbms.Db.Where(models.Tag{
				Name: s.Name,
			}).Last(&res).Error; err != nil {
				message = fmt.Sprintf("query tag failed: [%v]", err)
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

// EditTag 编辑标签 put /v1/tag
func (h *HTTPAPI) EditTag(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		res    models.Tag
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

		name := gbody.Get("name").String()
		s := models.Tag{
			Id:   id,
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
		if err = models.Dbms.Db.Updates(&s).Error; err != nil {
			isOk = false
			code = 1
			message = fmt.Sprintf("update tag failed, %v", err)
			h.logger.Errorf(message)
			goto exit
		} else {
			if err = models.Dbms.Db.Where(models.Tag{
				Id: s.Id,
			}).Last(&res).Error; err != nil {
				message = fmt.Sprintf("query tag failed: [%v]", err)
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

// DeleteTag 删除标签 delete /v1/tag
func (h *HTTPAPI) DeleteTag(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		res   models.Tag
		id, _ = strconv.Atoi(c.Params.ByName("id"))
	)

	if err = models.Dbms.Db.Where(models.Article{
		Id: id,
	}).Preload("Articles").Last(&res).Error; err != nil {
		message = fmt.Sprintf("query tag failed: [%v]", err)
		h.logger.Errorf(message)
		isOk = false
		code = 1
		goto exit
	} else {
		if len(res.Articles) != 0 {
			isOk = false
			code = 1
			message = fmt.Sprintf("tag %v have %v articles", res.Name, len(res.Articles))
			goto exit
		} else {
			if err = models.Dbms.Db.Where(map[string]interface{}{"id": id}).Delete(
				&models.Tag{},
			).Error; err != nil {
				message = fmt.Sprintf("delete tag failed, %v", err)
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
