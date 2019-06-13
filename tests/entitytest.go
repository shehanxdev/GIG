package tests

import (
	"GIG/app/models"
	"GIG/app/utility/requesthandlers"
	"github.com/revel/revel/testing"
)

var apiUrl = "http://localhost:9000/api/add"

type EntityTest struct {
	testing.TestSuite
}

func (t *EntityTest) Before() {
	println("Set up")
}

func (t *EntityTest) TestThatSearchApiWorks() {
	t.Get("/api/search?query=Sri%20Lanka")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
}

func (t *EntityTest) TestThatGetEntityApiWorks() {
	t.Get("/api/get/Sri%20Lanka")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
}

func (t *EntityTest) TestThatCreateEntityApiWorks() {
	entity := models.Entity{}
	entity.Title = "Sri Lanka"

	//save to db
	result, _ := requesthandlers.PostRequest(apiUrl, entity)
	defer result.Body.Close()
	t.AssertEqual(result.Status,"202 Accepted")
}

func (t *EntityTest) After() {
	println("Tear down")
}