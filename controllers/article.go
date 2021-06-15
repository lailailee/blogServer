package controllers

import (
	"blog/models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"io/ioutil"
	"strconv"
	"time"
)

// GetArticleList 获取文章列表 /v1/article
func (h *HTTPAPI) GetArticleList(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error
		res     []models.Article

		search        = c.DefaultQuery("search", "")
		orderBy       = c.DefaultQuery("orderBy", "createdAt")
		page, _       = strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _      = strconv.Atoi(c.DefaultQuery("limit", "8"))
		categoryId, _ = strconv.Atoi(c.DefaultQuery("categoryId", "0"))
		seriesId, _   = strconv.Atoi(c.DefaultQuery("seriesId", "0"))
		tagId, _      = strconv.Atoi(c.DefaultQuery("tagId", "0"))
		condition     = map[string]interface{}{}
		result        []map[string]interface{}
		count         int64
	)
	if categoryId != 0 {
		condition["categoryId"] = categoryId
	} else if seriesId != 0 {
		condition["seriesId"] = seriesId
	}
	var querySentence *gorm.DB
	if tagId == 0 {
		models.Dbms.Db.Model(&models.Article{}).Where("title LIKE ?", "%%"+search+"%%").Count(&count)
		querySentence = models.Dbms.Db.Where("title LIKE ?", "%%"+search+"%%").Where(condition).Preload("Category").Preload("Series").Preload("Tags").Order(orderBy + " desc").Offset((page - 1) * limit).Limit(limit).Find(&res)
	} else {
		var r models.Tag
		querySentence = models.Dbms.Db.Where(map[string]interface{}{"id": tagId}).Preload("Articles", func(db *gorm.DB) *gorm.DB {
			db.Model(&models.Article{}).Where("title LIKE ?", "%%"+search+"%%").Count(&count)
			return db.Where(condition).Where("title LIKE ?", "%%"+search+"%%").Order(orderBy + " desc").Offset((page - 1) * limit).Limit(limit)
		}).Preload("Articles.Tags").Preload("Articles.Series").Preload("Articles.Category").First(&r)
		res = r.Articles
	}
	if err = querySentence.Error; err != nil {
		isOk = false
		code = 1
		message = fmt.Sprintf("article get faileds: %v", err)
		h.logger.Errorf(message)
		goto exit
	} else {
		result = []map[string]interface{}{}
		for _, v := range res {
			tmp := map[string]interface{}{}
			tmp["id"] = v.Id
			tmp["title"] = v.Title
			tmp["content"] = v.Content
			tmp["overview"] = v.Overview
			tmp["categoryId"] = v.CategoryId
			tmp["category"] = v.Category
			tmp["seriesId"] = v.SeriesId
			tmp["series"] = v.Series
			tmp["tags"] = v.Tags
			tmp["seriesIndex"] = v.SeriesIndex
			tmp["viewCount"] = v.ViewCount
			tmp["createdAt"] = v.CreatedAt.Format("2006-01-02 15:04")
			tmp["updatedAt"] = v.UpdatedAt.Format("2006-01-02 15:04")
			result = append(result, tmp)
		}
		isOk = true
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

// GetArticle 获取文章详情 /v1/article/:id
func (h *HTTPAPI) GetArticle(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		res    models.Article
		id, _  = strconv.Atoi(c.Params.ByName("id"))
		result map[string]interface{}
	)

	querySentence := models.Dbms.Db.Preload("Category").Preload("Tags").Where(models.Article{Id: id}).First(&res)
	if err = querySentence.Error; err != nil {
		isOk = false
		code = 1
		message = fmt.Sprintf("article get faileds: %v", err)
		h.logger.Errorf(message)
		goto exit
	} else {
		isOk = true
		v := res
		result = map[string]interface{}{}
		result["id"] = v.Id
		result["title"] = v.Title
		result["content"] = v.Content
		result["overview"] = v.Overview
		result["categoryId"] = v.CategoryId
		result["category"] = v.Category
		result["tags"] = v.Tags
		result["seriesId"] = v.SeriesId
		result["series"] = v.Series
		result["viewCount"] = v.ViewCount
		result["createdAt"] = v.CreatedAt.Format("20060102T150405Z")
		result["updatedAt"] = v.UpdatedAt.Format("20060102T150405Z")

		h.lock.Lock()
		s := models.Article{
			Id:        id,
			ViewCount: v.ViewCount + 1,
		}
		models.Dbms.Db.Updates(&s)
		h.lock.Unlock()
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

// AddArticle 新增文章 /v1/article
func (h *HTTPAPI) AddArticle(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error
		res     models.Article
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
		title := gbody.Get("title").String()
		content := gbody.Get("content").String()
		overview := gbody.Get("overview").String()
		categoryId := int(gbody.Get("categoryId").Int())
		seriesId := int(gbody.Get("seriesId").Int())
		seriesIndex := int(gbody.Get("seriesIndex").Int())
		loc, _ := time.LoadLocation("Local") // 获取时区
		createdAt, _ := time.ParseInLocation("2006-01-02 15:04:05", gbody.Get("createdAt").String(), loc)
		tags := gbody.Get("tags").String()
		s := models.Article{
			Title:      title,
			Content:    content,
			Overview:   overview,
			CategoryId: categoryId,
			// SeriesId:   seriesId,
			SeriesIndex: seriesIndex,
			CreatedAt:   createdAt,
		}
		if seriesId != 0 {
			s.SeriesId = seriesId
		}
		json.Unmarshal([]byte(tags), &s.Tags)
		if err = models.Dbms.Db.Create(&s).Error; err != nil {
			isOk = false
			code = 1
			message = fmt.Sprintf("save article data failed: [%v]", err)
			h.logger.Errorf(message)
			goto exit
		} else {
			if err = models.Dbms.Db.Where(models.Article{
				Title: s.Title,
			}).Preload("Category").Preload("Tags").Last(&res).Error; err != nil {
				message = fmt.Sprintf("query article failed: [%v]", err)
				h.logger.Errorf(message)
				isOk = false
				code = 1
				goto exit
			} else {
				isOk = true
				v := res
				result = map[string]interface{}{}
				result["id"] = v.Id
				result["title"] = v.Title
				result["content"] = v.Content
				result["overview"] = v.Overview
				result["categoryId"] = v.CategoryId
				result["category"] = v.Category
				result["seriesId"] = v.SeriesId
				result["series"] = v.Series
				result["tags"] = v.Tags
				result["viewCount"] = v.ViewCount
				result["createdAt"] = v.CreatedAt.Format("2006-01-02 15:04")
				result["updatedAt"] = v.UpdatedAt.Format("2006-01-02 15:04")
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

// EditArticle 编辑文章 put /v1/article/:id
func (h *HTTPAPI) EditArticle(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		// params models.Article
		res    models.Article
		result []map[string]interface{}
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
		title := gbody.Get("title").String()
		content := gbody.Get("content").String()
		overview := gbody.Get("overview").String()
		categoryId := int(gbody.Get("categoryId").Int())
		seriesId := int(gbody.Get("seriesId").Int())
		seriesIndex := int(gbody.Get("seriesIndex").Int())
		loc, _ := time.LoadLocation("Local") // 获取时区
		created := gbody.Get("createdAt").String()
		createdAt, _ := time.ParseInLocation("2006-01-02 15:04:05", created, loc)

		tags := gbody.Get("tags").String()
		s := models.Article{
			Id:          id,
			Title:       title,
			Content:     content,
			Overview:    overview,
			CategoryId:  categoryId,
			SeriesIndex: seriesIndex,
			CreatedAt:   createdAt,
		}
		if seriesId != 0 {
			s.SeriesId = seriesId
		}
		json.Unmarshal([]byte(tags), &s.Tags)

		models.Dbms.Db.Model(&models.Article{}).Association("Tags").Clear()
		if err = models.Dbms.Db.Updates(&s).Error; err != nil {
			isOk = false
			code = 1
			message = fmt.Sprintf("update article failed, %v", err)
			h.logger.Errorf(message)
			goto exit
		} else {
			if err = models.Dbms.Db.Where(models.Article{
				Id: s.Id,
			}).Preload("Category").Preload("Tags").Last(&res).Error; err != nil {
				message = fmt.Sprintf("query article failed: [%v]", err)
				h.logger.Errorf(message)
				isOk = false
				code = 1
				goto exit
			} else {
				isOk = true
				v := res
				result := map[string]interface{}{}
				result["id"] = v.Id
				result["title"] = v.Title
				result["content"] = v.Content
				result["overview"] = v.Overview
				result["categoryId"] = v.CategoryId
				result["category"] = v.Category
				result["seriesId"] = v.SeriesId
				result["series"] = v.Series
				result["tags"] = v.Tags
				result["viewCount"] = v.ViewCount
				result["createdAt"] = v.CreatedAt.Format("20060102T150405Z")
				result["updatedAt"] = v.UpdatedAt.Format("20060102T150405Z")
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

// DeleteArticle  删除文章  delete /v1/article/:id
func (h *HTTPAPI) DeleteArticle(c *gin.Context) {
	var (
		isOk    bool
		content gin.H
		code    int
		message string
		err     error

		id, _ = strconv.Atoi(c.Params.ByName("id"))
	)

	models.Dbms.Db.Model(&models.Article{Id: id}).Association("Tags").Clear()
	if err = models.Dbms.Db.Where(
		map[string]interface{}{"id": id},
	).Delete(&models.Article{}).Error; err != nil {
		message = fmt.Sprintf("article delete failed, %v", err)
		h.logger.Errorf(message)
		isOk = false
		code = 5
		goto End
	} else {
		isOk = true
	}

End:
	if isOk {
		content = gin.H{
			"code":    0,
			"message": "success",
			"data": map[string]interface{}{
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

// // GetArticleTimeline 获取文章时间线 /v1/article/timeline
// func (h *HTTPAPI) GetArticleTimeline(c *gin.Context) {
// 	var (
// 		isOk    bool
// 		content gin.H
// 		code    int
// 		message string
// 		err     error
// 		res     []models.Article
//
//
// 		result []map[string]interface{}
// 	)
//
// 	querySentence := models.Dbms.Db.Preload("Category").Preload("Tags").Order("createdAt").Find(&res)
// 	if err = querySentence.Error; err != nil {
// 		isOk = false
// 		code = 1
// 		message = fmt.Sprintf("article get faileds: %v", err)
// 		h.logger.Errorf(message)
// 		goto exit
// 	} else {
// 		result = []map[string]interface{}{}
// 		for _, v := range res {
// 			tmp := map[string]interface{}{}
// 			tmp["id"] = v.Id
// 			tmp["title"] = v.Title
// 			tmp["viewCount"] = v.ViewCount
// 			tmp["createdAt"] = v.CreatedAt.Format("20060102T150405Z")
// 			result = append(result, tmp)
// 		}
// 		isOk = true
// 	}
//
// exit:
// 	if isOk {
// 		content = gin.H{
// 			"code":    0,
// 			"message": "success",
// 			"data": map[string]interface{}{
// 				"list":  result,
// 				"count": len(result),
// 			},
// 		}
// 	} else {
// 		content = gin.H{
// 			"code":    code,
// 			"message": message,
// 		}
// 	}
// 	c.JSON(200, content)
// }
