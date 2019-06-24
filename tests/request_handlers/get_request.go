package request_handlers

import (
	"GIG/app/utility/request_handlers"
)

func (t *TestRequestHandlers) TestThatGetRequestWorks() {
	link := "http://www.buildings.gov.lk/index.php"
	result, _ := request_handlers.GetRequest(link)
	defer result.Body.Close()
	t.AssertEqual(result.Status,"200 OK")
}