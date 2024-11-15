package vo

type Code2SessionResponse struct {
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
}

type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type IdentifyCode struct {
	EncryptedData string `json:"encrypted_data" binding:"required"`
	Iv            string `json:"iv" binding:"required"`
	SessionKey    string `json:"session_id" binding:"required"`
}

type UpdateUserRequest struct {
	Telephone string `json:"telephone"`
	NickName  string `json:"nick_name"`
	Motto     string `json:"motto"`
	Gender    int    `json:"gender"` //0表示男、1表示女
}

type ShowUserRequest struct {
	Telephone string `json:"telephone"`
	NickName  string `json:"nick_name"`
	Motto     string `json:"motto"`
	Gender    int    `json:"gender"` //0表示男、1表示女
}
