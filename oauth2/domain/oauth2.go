package domain

type Oauth2Info struct {
	Uid         string `json:"id"`
	Name        string `json:"name"`
	Username    string `json:"username"`
	AccessToken string `json:"access_token"`
}
