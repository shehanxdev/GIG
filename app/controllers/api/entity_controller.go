package api

import (
	"GIG-SDK/libraries"
	"GIG-SDK/models"
	"GIG/app/controllers"
	"GIG/app/repositories"
	"GIG/app/storages"
	"errors"
	"github.com/revel/revel"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type EntityController struct {
	*revel.Controller
}

func (c EntityController) Search() revel.Result {
	var (
		entities []models.Entity
		err      error
	)
	searchKey := c.Params.Values.Get("query")
	categories := c.Params.Values.Get("categories")
	attributes := c.Params.Values.Get("attributes")
	limit, limitErr := strconv.Atoi(c.Params.Values.Get("limit"))
	page, pageErr := strconv.Atoi(c.Params.Values.Get("page"))
	if pageErr != nil || page < 1 {
		page = 1
	}

	c.Response.Out.Header().Set("Access-Control-Allow-Origin", "*")

	if limitErr != nil {
		errResp := controllers.BuildErrResponse(400, errors.New("result limit is required"), )
		c.Response.Status = 400
		return c.RenderJSON(errResp)
	}

	categoriesArray := libraries.ParseCategoriesString(categories)
	attributesArray := libraries.ParseCategoriesString(attributes)

	if searchKey == "" && categories == "" {
		errResp := controllers.BuildErrResponse(400, errors.New("search value or category is required"), )
		c.Response.Status = 400
		return c.RenderJSON(errResp)
	}

	var responseArray []models.SearchResult
	entities, err = repositories.EntityRepository{}.GetEntities(searchKey, categoriesArray, limit, (page-1)*limit)
	if err != nil {
		log.Println(err)
		errResp := controllers.BuildErrResponse(500, err)
		c.Response.Status = 500
		return c.RenderJSON(errResp)
	}

	for _, element := range entities {
		responseArray = append(responseArray, models.SearchResult{}.ResultFrom(element, attributesArray))
	}
	c.Response.Status = 200
	return c.RenderJSON(responseArray)
}

func (c EntityController) Show(title string) revel.Result {
	var (
		entity models.Entity
		err    error
	)
	log.Println("title", title)
	c.Response.Out.Header().Set("Access-Control-Allow-Origin", "*")

	if title == "" {
		errResp := controllers.BuildErrResponse(400, errors.New("invalid entity id format"))
		c.Response.Status = 400
		return c.RenderJSON(errResp)
	}
	dateParam := strings.Split(c.Params.Values.Get("date"), "T")[0]
	entityDate, dateError := time.Parse("2006-01-02", dateParam)
	attributes := c.Params.Values.Get("attributes")
	defaultImageOnly := c.Params.Values.Get("imageOnly")
	attributesArray := libraries.ParseCategoriesString(attributes)

	if dateError != nil || entityDate.IsZero() {
		entity, err = repositories.EntityRepository{}.GetEntityBy("title", title)
	} else {
		entity, err = repositories.EntityRepository{}.GetEntityByPreviousTitle(title, entityDate)
	}

	if err != nil {
		var normalizedName string
		normalizedName, err = repositories.EntityRepository{}.NormalizeEntityTitle(title)
		if err == nil {
			if dateError != nil || entityDate.IsZero() {
				entity, err = repositories.EntityRepository{}.GetEntityBy("title", normalizedName)
			} else {
				entity, err = repositories.EntityRepository{}.GetEntityByPreviousTitle(normalizedName, entityDate)
			}
		}
	}

	if err != nil {
		log.Println(err)
		errResp := controllers.BuildErrResponse(500, err)
		c.Response.Status = 500
		return c.RenderJSON(errResp)
	}

	// return only the default image
	if defaultImageOnly == "true" {
		var (
			localFile *os.File
			err       error
		)
		imageUrl := entity.GetImageURL()
		imagePathArray := strings.Split(imageUrl, "/")
		c.Response.Status = 404
		if len(imagePathArray) != 3 {
			return c.RenderJSON("default image not found")
		}

		localFile, err = storages.FileStorageHandler{}.GetFile(imagePathArray[1], imagePathArray[2])
		if err != nil {
			return c.RenderJSON(err)
		}

		c.Response.Status = 200
		return c.RenderFile(localFile, revel.Inline)

	}

	c.Response.Status = 200
	return c.RenderJSON(models.SearchResult{}.ResultFrom(entity, attributesArray))
}

func (c EntityController) CreateBatch() revel.Result {
	var (
		entities      []models.Entity
		savedEntities []models.Entity
	)
	log.Println("create entity batch request")
	err := c.Params.BindJSON(&entities)
	if err != nil {
		errResp := controllers.BuildErrResponse(403, err)
		c.Response.Status = 403
		return c.RenderJSON(errResp)
	}

	for _, e := range entities {
		entity, _, err := repositories.EntityRepository{}.AddEntity(e)
		if err != nil {
			errResp := controllers.BuildErrResponse(500, err)
			c.Response.Status = 500
			return c.RenderJSON(errResp)
		}
		savedEntities = append(savedEntities, entity)
	}

	c.Response.Status = 200
	return c.RenderJSON(savedEntities)
}

func (c EntityController) Create() revel.Result {
	var (
		entity models.Entity
		err    error
	)
	log.Println("create entity request")
	err = c.Params.BindJSON(&entity)
	if err != nil {
		log.Println("binding error:", err)
		errResp := controllers.BuildErrResponse(403, err)
		c.Response.Status = 403
		return c.RenderJSON(errResp)
	}
	entity, c.Response.Status, err = repositories.EntityRepository{}.AddEntity(entity)
	if err != nil {
		log.Println("entity create error:", err)
		errResp := controllers.BuildErrResponse(500, err)
		c.Response.Status = 500
		return c.RenderJSON(errResp)
	}
	return c.RenderJSON(entity)

}
