package controllers

import (
	"github.com/revel/revel"
	"shuttle/app/models"
	"shuttle/app/services"
	"time"
	"github.com/revel/revel/cache"
	"strings"
	"net/url"
	"net/http"
	"encoding/json"
	"io/ioutil"
)

const VALIDATION_URL_ERROR_MESSAGE = "Does not match an URL pattern."
const VALIDATION_DOMAIN_ERROR_MESSAGE = "Already shortened."
const APPLICATION_JSON = "application/json; charset=utf-8"
const APPLICATION_FORM_ENCODED = "application/x-www-form-urlencoded"
const PAYLOAD_ERROR = "error"
const CACHE_RESOURCE = "resource_"

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	RCaptchaStatus, _ := revel.Config.Bool("recaptcha.enable")
	RCaptchaKey, _ := revel.Config.String("recaptcha.key")
	return c.Render(RCaptchaKey, RCaptchaStatus)
}

func (c App) Resource() revel.Result {

	payload := make(map[string]interface{})

	// check recaptcha
	status, _ := revel.Config.Bool("recaptcha.enable")
	if status {
		ValidRecaptcha(&c)
		if c.Validation.HasErrors() {
			payload[PAYLOAD_ERROR] = c.Validation.Errors
			return c.RenderJSON(payload)
		}
	}

	c.Response.ContentType = APPLICATION_JSON
	longURL := c.Params.Form.Get("url")

	// lightweight url validation, find a stricter way maybe
	u, err := url.ParseRequestURI(longURL)
	if err != nil {
		c.Validation.Error(VALIDATION_URL_ERROR_MESSAGE)
		payload[PAYLOAD_ERROR] = c.Validation.Errors
		return c.RenderJSON(payload)
	}
	if 0 == strings.Compare(u.Host, c.Request.Host) {
		c.Validation.Error(VALIDATION_DOMAIN_ERROR_MESSAGE)
	}
	if c.Validation.HasErrors() {
		payload[PAYLOAD_ERROR] = c.Validation.Errors
		return c.RenderJSON(payload)
	}
	var resource *models.Resource
	if err := cache.Get(CACHE_RESOURCE+longURL, &resource); err != nil {
		resource = services.FindByLongURL(longURL)
		if nil != resource {
			go cache.Set(CACHE_RESOURCE+longURL, resource, 24*time.Hour)
		}
	}
	if nil == resource {
		resource = services.Create(models.Resource {
			LongURL: longURL,
		})
	}
	payload["url"] = services.Encode(resource.Id)
	return c.RenderJSON(payload)
}

func (c App) Retrieve(uri string) revel.Result  {
	id := services.Decode(uri)
	var resource *models.Resource
	if err := cache.Get(CACHE_RESOURCE+string(id), &resource); err != nil {
		resource = services.Get(id)
		if nil != resource {
			go cache.Set(CACHE_RESOURCE+string(id), resource, 24*time.Hour)
		}
	}
	if nil == resource {
		return c.NotFound("Resource not found")
	}
	resource.LastRetrieval = time.Now()
	resource = services.Update(resource)
	return c.Redirect(resource.LongURL)
}

func ValidRecaptcha(c *App) {
	// reCaptcha
	recaptchaSecret, _ := revel.Config.String("recaptcha.secret")
	recaptchaURL, _ := revel.Config.String("recaptcha.url")
	form := url.Values{}
	form.Set("secret", recaptchaSecret)
	form.Set("response", c.Params.Form.Get("response"))
	request, err := http.NewRequest("POST", recaptchaURL, strings.NewReader(form.Encode()))
	request.Header.Set("Content-Type", APPLICATION_FORM_ENCODED)
	services.CheckErr(err)
	client := http.Client{}
	response, err := client.Do(request)
	services.CheckErr(err)
	i, err := ioutil.ReadAll(response.Body)
	services.CheckErr(err)
	tmp := make(map[string]json.RawMessage)
	err = json.Unmarshal(i, &tmp)
	services.CheckErr(err)
	if "false" == string(tmp["success"]) {
		c.Validation.Error("reCaptcha denied")
		var errors []string
		err = json.Unmarshal(tmp["error-codes"], &errors)
		services.CheckErr(err)

		for _, recaptchaErr := range errors {
			c.Validation.Error(recaptchaErr)
		}
	}
}
