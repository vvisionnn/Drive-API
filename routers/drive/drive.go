package drive

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vvisionnn/Drive-API/pkgs/onedrive"
	"github.com/vvisionnn/Drive-API/pkgs/response"
	"github.com/vvisionnn/Drive-API/settings"
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

	response.SuccessWithMessage(ctx, status)
}

func UrlHandler(ctx *gin.Context) {
	var oauthUrl = ""
	finalUrl := ctx.DefaultQuery("path", "")
	if len(finalUrl) == 0 {
		oauthUrl = Drive.GetOAuthURI()
	} else {
		oauthUrl = Drive.GetOAuthURI(finalUrl)
	}
	response.SuccessWithData(ctx, oauthUrl)
}

func CallbackHandler(ctx *gin.Context) {
	code := ctx.DefaultQuery("code", "")
	finalUrl := ctx.DefaultQuery("state", "")

	if len(code) == 0 {
		response.InternalServerError(ctx, "get code error")
		return
	}
	if err := Drive.UpdateCredential(code); err != nil {
		response.InternalServerError(ctx, err.Error())
		return
	}
	response.RedirectTemporary(ctx, finalUrl)
}

func ListHandler(ctx *gin.Context) {
	if !Drive.LoginStatus() {
		response.SuccessWithMessage(ctx, "login first")
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
		response.SuccessWithMessage(ctx, "get access token error")
		return
	}

	url := fmt.Sprintf("%s/%s", settings.CONF.OnedriveEndpoint, suffix)

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("err:", err)
		response.SuccessWithMessage(ctx, "request graph api error")
		return
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	items := onedrive.ListResponse{}
	if err := json.Unmarshal(respBody, &items); err != nil {
		// todo: add log
		fmt.Println("list error: ", err)
	}
	response.SuccessWithData(ctx, items)
}
