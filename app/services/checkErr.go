package services

import (
	"shuttle/app"
)

func CheckErr(err error) {
	if err != nil {
		app.AppLog.Fatalf(err.Error())
	}
}
