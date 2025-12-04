package jwts

import (
	"errors"
	"fmt"
	"fungo/common/response"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/rest/httpx"
)

type JwtPayLoad struct {
	UserID   uint64 `json:"user_id"`
	UserName string `json:"user_name"`
	Role     int    `json:"role"`
}

type CustomClaims struct {
	JwtPayLoad
	jwt.RegisteredClaims
}

func GenToken(user JwtPayLoad, accessSecret string, expires int64) (string, error) {
	claim := CustomClaims{
		JwtPayLoad: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(expires))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(accessSecret))
}

// ParseToken 解析Token
func ParseToken(tokenStr string, accessSecret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func JwtUnauthorizedResult(w http.ResponseWriter, r *http.Request, err error) {
	fmt.Println(err)
	httpx.WriteJson(w, http.StatusOK, response.Body{Code: 10087, Msg: "鉴权失败"})
}
