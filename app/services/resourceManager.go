package services

import (
	"shuttle/app"
	"shuttle/app/models"
	"math"
	"strings"
)

const ALPHABET = "123456789abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ"
const BASE = float64(len(ALPHABET))

func Get(id int64) *models.Resource {
	resource, err := app.DB.Get(models.Resource{}, id)
	CheckErr(err)
	if nil == resource {
		return nil
	}
	return resource.(*models.Resource)
}

func FindByLongURL(longURL string) *models.Resource {
	var resource *models.Resource
	app.DB.SelectOne(&resource, "select * from resources where long_url=?", longURL)
	if nil == resource {
		return nil
	}
	return resource
}

func Create(resource models.Resource) *models.Resource {
	err := app.DB.Insert(&resource)
	CheckErr(err)
	return &resource
}

func Update(resource *models.Resource) *models.Resource {
	_, err := app.DB.Update(resource)
	CheckErr(err)
	return resource
}

func Encode(i int64) string {
	var encoded = ""
	for i != 0 {
		var remainder = math.Mod(float64(i), BASE)
		i = int64(math.Floor(float64(i)/BASE))
		encoded = string([]rune(ALPHABET)[int(remainder)]) + encoded
	}
	return encoded
}

func Decode(s string) int64 {
	var decoded = 0
	for s != "" {
		i := strings.Index(ALPHABET, string([]rune(s)[0]))
		power := len(s)-1
		decoded += i*int(math.Pow(BASE, float64(power)))
		s = s[1:]
	}
	return int64(decoded)
}
