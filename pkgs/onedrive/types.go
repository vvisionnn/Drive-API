package onedrive

type RefreshResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    uint   `json:"expires_in"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ListResponse struct {
	Value   []ItemInfo `json:"value"`
	Context string     `json:"@odata.context"`
}

type ItemInfo struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Size        uint64  `json:"size"`
	DownloadURL string  `json:"@microsoft.graph.downloadUrl"`
	File        *file   `json:"file"`
	Folder      *folder `json:"folder"`
}

type file struct {
	MimeType string `json:"mimeType"`
}

type folder struct {
	ChildCount int `json:"childCount"`
}
