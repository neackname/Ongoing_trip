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
