package drive

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/vvisionnn/Drive-API/pkgs/onedrive"
	"github.com/vvisionnn/Drive-API/pkgs/response"
	"github.com/vvisionnn/Drive-API/settings"
	"io/ioutil"
	"log"
	"os"
)

var Drive *onedrive.Client

func SetConfiguration(ctx *gin.Context) {
	var err error
	// todo: check if curr and previous conf are same one
	data := struct {
		ClientId     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RedirectUrl  string `json:"redirect_url"`
		Path         string `json:"path"`
	}{}
	if err = ctx.ShouldBindJSON(&data); err != nil {
		log.Println(err)
		response.SuccessWithMessage(ctx, "params error")
		return
	}
	settings.CONF = &settings.Configuration{
		ClientId:         data.ClientId,
		ClientSecret:     data.ClientSecret,
		RedirectUrl:      data.RedirectUrl,
		Scopes:           []string{"Files.Read.All"},
		OauthEndpoint:    "https://login.microsoftonline.com/common/oauth2/v2.0",
		Authority:        "https://login.microsoftonline.com/common",
		OnedriveEndpoint: "https://graph.microsoft.com/v1.0/me",
	}
	settings.CONF.Save()

	Drive = onedrive.NewClient(
		settings.CONF.ClientId,
		settings.CONF.ClientSecret,
		settings.CONF.OauthEndpoint,
		settings.CONF.RedirectUrl,
		settings.CONF.Scopes,
	)
	// get oauth url and return
	response.SuccessWithData(ctx, Drive.GetOAuthURI(data.Path))
}

func InitialDrive() error {
	// if conf doesn't exist, just return and wait
	if settings.CONF == nil { return nil }
	log.Println("initial drive from previous settings")

	// if previous config found, initial drive
	Drive = onedrive.NewClient(
		settings.CONF.ClientId,
		settings.CONF.ClientSecret,
		settings.CONF.OauthEndpoint,
		settings.CONF.RedirectUrl,
		settings.CONF.Scopes,
	)

	// check if the token exist
	if _, err := os.Stat(".tokens.json"); os.IsNotExist(err) {
		log.Println("tokens not found, creating file...")
		// isn't exist, just create file, update after login
		_, err := os.OpenFile(".tokens.json", os.O_CREATE|os.O_RDWR, 0766)
		if err != nil {
			return err
		}
	} else {
		log.Println("found previous tokens file, update drive status from file...")
		// exist, initial token from file
		content, err := ioutil.ReadFile(".tokens.json")
		if err != nil {
			return err
		}
		ts := onedrive.Tokens{}
		if err := json.Unmarshal(content, &ts); err != nil {
			return err
		}
		Drive.AccessToken = ts.AccessToken
		Drive.RefreshToken = ts.RefreshToken
		Drive.AccessTokenExpireTime = ts.AccessTokenExpireTime
		log.Println("update from file done.")
	}
	return nil
}

func StatusHandler(ctx *gin.Context) {
	var status string
	confStat := settings.CONF != nil
	driveStat := Drive != nil && Drive.LoginStatus()
	if confStat && driveStat {
		status = "ok"
	} else {
		status = "login first"
	}

	response.SuccessWithMessage(ctx, status)
}

//func UrlHandler(ctx *gin.Context) {
//	var oauthUrl = ""
//	finalUrl := ctx.DefaultQuery("path", "")
//	if len(finalUrl) == 0 {
//		oauthUrl = Drive.GetOAuthURI()
//	} else {
//		oauthUrl = Drive.GetOAuthURI(finalUrl)
//	}
//	response.SuccessWithData(ctx, oauthUrl)
//}

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

func ListRootHandler(ctx *gin.Context) {
	items, err := Drive.ListRootChildren()
	if err != nil {
		log.Println(err)
		response.SuccessWithMessage(ctx, "list children error")
		return
	}
	response.SuccessWithData(ctx, *items)
}

func ListHandler(ctx *gin.Context) {
	items, err := Drive.ListItemChildren(ctx.Param("id"))
	if err != nil {
		log.Println(err)
		response.SuccessWithMessage(ctx, "list children error")
		return
	}
	response.SuccessWithData(ctx, *items)
}
