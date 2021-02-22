package onedrive

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vvisionnn/Drive-API/settings"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	ClientID              string
	ClientSecret          string
	Endpoint              string
	RedirectURI           string
	AccessToken           string
	AccessTokenExpireTime time.Time
	RefreshToken          string
	OauthURI              string
	Scopes                []string

	HttpClient				http.Client
}

// NewClient create a new onedrive instance
func NewClient(clientID, clientSecret, endpoint, redirectURI string, scopes []string) *Client {
	return &Client{
		ClientID:              clientID,
		ClientSecret:          clientSecret,
		Endpoint:              endpoint,
		RedirectURI:           redirectURI,
		AccessToken:           "",
		AccessTokenExpireTime: time.Time{},
		RefreshToken:          "",
		Scopes:                scopes,
		OauthURI: fmt.Sprintf("%s/authorize?"+
			"client_id=%s"+
			"&response_type=code"+
			"&redirect_uri=%s"+
			"&response_mode=query"+
			"&scope=offline_access %s",
			endpoint,
			clientID,
			redirectURI,
			strings.Join(scopes, " ")),
		HttpClient: http.Client{Timeout: time.Second * 10},
	}
}

type Tokens struct {
	AccessToken           string    `json:"access_token"`
	AccessTokenExpireTime time.Time `json:"access_token_expire_time"`
	RefreshToken          string    `json:"refresh_token"`
}

// LoginStatus return if onedrive instance logged in
func (drive *Client) LoginStatus() bool {
	return drive != nil && len(drive.RefreshToken) > 0
}

// GetAccessToken return the onedrive access token, refresh if needed
func (drive *Client) GetAccessToken() (string, error) {
	if drive.AccessTokenExpireTime.Before(time.Now()) {
		// update access token by refresh token
		if err := drive.UpdateCredential(); err != nil {
			return "", err
		}
	}
	return drive.AccessToken, nil
}

// UpdateCredential update refresh token by code if provide, or just refresh it
func (drive *Client) UpdateCredential(code ...string) error {
	headers := map[string]string{
		"Content-Type": "application/x-www-form-urlencoded",
	}
	tokenURI := fmt.Sprintf("%s/token", drive.Endpoint)
	data := url.Values{
		"client_id":     {drive.ClientID},
		"client_secret": {drive.ClientSecret},
		"redirect_uri":  {drive.RedirectURI},
		"scope":         {"offline_access " + strings.Join(drive.Scopes, " ")},
	}

	switch len(code) {
	case 0:
		data.Add("refresh_token", drive.RefreshToken)
		data.Add("grant_type", "refresh_token")
	case 1:
		data.Add("code", code[0])
		data.Add("grant_type", "authorization_code")
	default:
		return errors.New("data length error")
	}

	// build request
	req, err := http.NewRequest(
		"POST",
		tokenURI,
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return errors.New("create request error")
	}

	// add headers
	for key, val := range headers {
		req.Header.Add(key, val)
	}

	resp, _ := drive.HttpClient.Do(req)
	respBody, _ := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	_resp := RefreshResponse{}
	if err := json.Unmarshal(respBody, &_resp); err != nil {
		return err
	}

	// update .tokens.json file
	ts := Tokens{
		AccessToken:           _resp.AccessToken,
		AccessTokenExpireTime: time.Now().Add(time.Duration(_resp.ExpiresIn) * time.Second),
		RefreshToken:          _resp.RefreshToken,
	}
	tsStr, err := json.Marshal(ts)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./data/tokens.json", tsStr, 0644)
	if err != nil {
		return err
	}

	drive.AccessToken = _resp.AccessToken
	drive.RefreshToken = _resp.RefreshToken
	drive.AccessTokenExpireTime = time.Now().Add(time.Duration(_resp.ExpiresIn) * time.Second)

	return nil
}

func (drive *Client) GetOAuthURI(state ...string) string {
	if len(state) == 0 {
		return drive.OauthURI
	}
	return fmt.Sprintf("%s&state=%s", drive.OauthURI, state[0])
}

func (drive *Client) ListRootChildren() (*ListResponse, error) {
	suffix := "drive/root/children"
	_url := fmt.Sprintf("%s/%s", settings.CONF.OnedriveEndpoint, suffix)
	return drive.ListChildren(_url)
}

func (drive *Client) ListItemChildren(itemId string) (*ListResponse, error) {
	suffix := fmt.Sprintf("drive/items/%s/children", itemId)
	_url := fmt.Sprintf("%s/%s", settings.CONF.OnedriveEndpoint, suffix)
	return drive.ListChildren(_url)
}

func (drive *Client) ListChildren(url string) (*ListResponse, error) {
	accessToken, err := drive.GetAccessToken()
	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := drive.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(respBody))
	items := ListResponse{}
	if err := json.Unmarshal(respBody, &items); err != nil {
		return nil, err
	}
	return &items, nil
}
