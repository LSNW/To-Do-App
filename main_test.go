package main

import (
	"bytes"
	//"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	//"github.com/gofiber/fiber/v2"
	//"github.com/gofiber/template/html"
	"log"

	"github.com/stretchr/testify/assert"
)

var app = initApp()

func TestInvalidLogin(t *testing.T) {
	data := url.Values{}
	data.Set("login", "non")
	data.Add("password", "existent")
	req, err := http.NewRequest(
		"POST",
		"/login/auth",
		bytes.NewBufferString(data.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	if err != nil {
		log.Println(err)
	}

	res, err := app.Test(req, -1)
	url, _ := res.Location()

	assert.Equalf(t, false, err != nil, "Failed to create wrong username password request")
	assert.Equalf(t, 302, res.StatusCode, "Failed checking wrong username password")
	assert.Equalf(t, "/login", url.Path, "Failed checking wrong username password")
}

func TestInvalidPassword(t *testing.T) {
	data := url.Values{}
	data.Set("login", "please")
	data.Add("password", "incorrect")
	req, err := http.NewRequest(
		"POST",
		"/login/auth",
		bytes.NewBufferString(data.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	if err != nil {
		log.Println(err)
	}

	res, err := app.Test(req, -1)
	url, _ := res.Location()

	assert.Equalf(t, false, err != nil, "Failed to create wrong password request")
	assert.Equalf(t, 302, res.StatusCode, "Failed checking wrong password")
	assert.Equalf(t, "/login", url.Path, "Failed checking wrong password")
}

func TestValidLogin(t *testing.T) {
	data := url.Values{}
	data.Set("login", "please")
	data.Add("password", "work")
	req, err := http.NewRequest(
		"POST",
		"/login/auth",
		bytes.NewBufferString(data.Encode()),
	)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; param=value")
	if err != nil {
		log.Println(err)
	}

	res, err := app.Test(req, -1)
	url, _ := res.Location()

	assert.Equalf(t, false, err != nil, "Failed to create right username password request")
	assert.Equalf(t, 302, res.StatusCode, "Failed checking right username password")
	assert.Equalf(t, "/", url.Path, "Failed checking right username password")
}