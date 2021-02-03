package drive

import (
	"Drive-API/pkgs/onedrive"
	"Drive-API/settings"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

var Drive = onedrive.NewClient(
	settings.CONF.AppId,
	settings.CONF.AppSecret,
	settings.CONF.OauthEndpoint,
	settings.CONF.Redirect,
	settings.CONF.Scopes,
)

func StatusHandler(ctx *gin.Context) {
	var status string
	if Drive.LoginStatus() {
		status = "ok"
	} else {
		status = "login first"
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": status,
	})
}

func UrlHandler(ctx *gin.Context) {
	finalUrl := ctx.DefaultQuery("path", "")
	if len(finalUrl) == 0 {
		ctx.String(http.StatusOK, Drive.GetOAuthURI())
		return
	}
	//fmt.Println(Drive.GetOAuthURI(finalUrl))
	ctx.String(http.StatusOK, Drive.GetOAuthURI(finalUrl))
}

func CallbackHandler(ctx *gin.Context) {
	code := ctx.DefaultQuery("code", "")
	fmt.Println("code: " + code)
	if len(code) == 0 {
		ctx.String(http.StatusOK, "code error")
		return
	}
	if err := Drive.UpdateCredential(code); err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	finalUrl := ctx.DefaultQuery("state", "")
	ctx.Redirect(http.StatusTemporaryRedirect, finalUrl)
}

func ListHandler(ctx *gin.Context) {
	if !Drive.LoginStatus() {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "login first",
		})
		return
	}

	var suffix string
	if itemId := ctx.DefaultQuery("id", ""); len(itemId) == 0 {
		suffix = "drive/root/children"
	} else {
		suffix = fmt.Sprintf("drive/items/%s/children", itemId)
	}

	accessToken, err := Drive.GetAccessToken()
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "get access token error",
		})
		return
	}

	client := &http.Client{}
	url := fmt.Sprintf("%s/%s", settings.CONF.OnedriveEndpoint, suffix)
	fmt.Println("url: ", url)
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err:", err)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "request error",
		})
		return
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	items := onedrive.ListResponse{}
	if err := json.Unmarshal(respBody, &items); err != nil {
		fmt.Println("list error: ", err)
	}
	ctx.JSON(http.StatusOK, items)
}
