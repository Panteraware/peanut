package peanut

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"strings"
)

func proxyPrivateDownload(token string, url string, c echo.Context) error {
	headers := []CustomHeaders{{Label: "Accept", Value: "application/octet-stream"}}

	strings.Replace(url, "https://api.github.com/", fmt.Sprintf("https://%s@api.github.com/", token), 1)

	getRes := getRequest(url, headers)

	c.Response().Header().Set("Location", getRes.Header.Get("Location"))

	return nil
}

func stringLength(str string) int {
	return strings.Count(str, "")
}

func getJson(body io.ReadCloser, target interface{}) error {
	defer body.Close()

	return json.NewDecoder(body).Decode(target)
}

func checkAlias(s string) (string, error) {
	if s == "win" || s == "windows" || s == "win32" {
		return "exe", nil
	} else if s == "debian" {
		return "deb", nil
	} else if s == "appimage" {
		return "AppImage", nil
	} else if s == "dmg" {
		return "dmg", nil
	} else if s == "mac" || s == "macos" || s == "osx" {
		return "darwin", nil
	} else if s == "fedora" {
		return "rpm", nil
	} else {
		return "", errors.New("invalid alias")
	}
}
