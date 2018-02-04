package tests

import (
	"github.com/revel/revel/testing"
	"strings"
)

type AppTest struct {
	testing.TestSuite
}

func (t *AppTest) Before() {
	println("Set up")
}

func (t *AppTest) TestIndex() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *AppTest) TestRetrieve404() {
	t.Get("/test")
	t.AssertNotFound()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *AppTest) TestResourceBadURLParameter() {
	reader := strings.NewReader("")
	t.Post("/", "Content-Type: application/json", reader)
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
}

func (t *AppTest) After() {
	println("Tear down")
}
