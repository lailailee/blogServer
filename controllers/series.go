package controllers

import (
	"blog/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"strconv"
)

// GetSeries 获取系列列表 get /v1/series
func (h *HTTPAPI) GetSeries(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string

		res      []models.Series
		result   []map[string]interface{}
		orderBy  = c.DefaultQuery("orderBy", "createdAt")
		page, _  = strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ = strconv.Atoi(c.DefaultQuery("limit", "8"))
		search   = c.DefaultQuery("search", "")
		count    int64
	)

	models.Dbms.Db.Model(&models.Series{}).Where("name LIKE ?", "%%"+search+"%%").Count(&count)
	if err := models.Dbms.Db.Where("name LIKE ?", "%%"+search+"%%").Preload("Articles").Order(orderBy + " desc").Offset((page - 1) * limit).Limit(limit).Find(&res).Error; err != nil {
		isOk = false
		code = 1
		message = fmt.Sprintf("series get failed: %v", err)
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

// AddSeries 新增系列 post /v1/series
func (h *HTTPAPI) AddSeries(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error
		res     models.Series
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
		s := models.Series{
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
			message = fmt.Sprintf("update series failed, %v", err)
			h.logger.Errorf(message)
			goto exit
		} else {
			if err = models.Dbms.Db.Where(models.Series{
				Name: s.Name,
			}).Last(&res).Error; err != nil {
				message = fmt.Sprintf("query series failed: [%v]", err)
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

// EditSeries 编辑系列 put /v1/series
func (h *HTTPAPI) EditSeries(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		res    models.Series
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

		s := models.Series{
			Id:   id,
			Name: gbody.Get("name").String(),
		}

		// models.Dbms.Db.Model(&s).Association("Tags").Clear()
		if err = models.Dbms.Db.Updates(&s).Error; err != nil {
			isOk = false
			code = 1
			message = fmt.Sprintf("update series failed, %v", err)
			h.logger.Errorf(message)
			goto exit
		} else {
			if err = models.Dbms.Db.Where(models.Series{
				Id: s.Id,
			}).Last(&res).Error; err != nil {
				message = fmt.Sprintf("query series failed: [%v]", err)
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

// DeleteSeries 删除系列 delete /v1/series
func (h *HTTPAPI) DeleteSeries(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		res   models.Series
		id, _ = strconv.Atoi(c.Params.ByName("id"))
	)

	if err = models.Dbms.Db.Where(models.Article{
		Id: id,
	}).Preload("Articles").Last(&res).Error; err != nil {
		message = fmt.Sprintf("query series failed: [%v]", err)
		h.logger.Errorf(message)
		isOk = false
		code = 1
		goto exit
	} else {
		if len(res.Articles) != 0 {
			isOk = false
			code = 1
			message = fmt.Sprintf("series %v have %v articles", res.Name, len(res.Articles))
			goto exit
		} else {
			if err = models.Dbms.Db.Where(map[string]interface{}{"id": id}).Delete(
				&models.Series{},
			).Error; err != nil {
				message = fmt.Sprintf("delete series failed, %v", err)
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
