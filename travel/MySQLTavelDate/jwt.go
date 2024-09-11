package MySQLTavelDate

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"travel/TravelModel"
)

var jwtKey = []byte("I_am_Snactop")

type Claim struct {
	UserId uint
	jwt.StandardClaims
}

func ReleaseToken(user TravelModel.TraUser) (string, error) {
	//TODO 设定TOKEN过期时间
	expiresTime := time.Now().Add(24 * time.Hour)

	//TODO 创建声明
	claim := &Claim{
		UserId: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiresTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "Snactop",
			Subject:   "user token",
		},
	}

	//TODO 签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	//TODO 加密
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	//TODO 发放TOKEN
	return tokenString, nil
}

// ParseToken @title ParseToken
// @description	解析token
// @auth	Snactop		2023-11-13	19:27
// @param	tokenString	string	传入token字符串
// @return	*jwt.Token, *Claim, error	传出标准的jwt认证的token、解析后的token内容以及错误信息
func ParseToken(tokenString string) (*jwt.Token, *Claim, error) {
	claim := &Claim{}

	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})

	return token, claim, err
}
