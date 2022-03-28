package web

import (
	"errors"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type Claims struct {
	UserId string `json:"user_id"`
	jwt.StandardClaims
}

const tokenExpireDuration = time.Hour * 5

var mySecret = []byte(os.Getenv("SECRET"))

func GenToken(userId string) (string, error) {
	c := Claims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenExpireDuration).Unix(),
			Issuer:    os.Getenv("ISSUER"),
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	return token.SignedString(mySecret)
}
func ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (i interface{}, err error) {
		return mySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func AdminAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("token")
		if authHeader == "" {
			reportError(c, AuthenticationError, "Empty header")
			c.Abort()
			return
		}
		claims, err := ParseToken(authHeader)
		if err != nil {
			reportError(c, AuthenticationError, "Invalid token "+err.Error())
			c.Abort()
			return
		}
		if time.Now().Unix() > claims.ExpiresAt {
			reportError(c, AuthenticationError, "Token expired")
			c.Abort()
		}
		c.Next()
	}
}
func WxAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("token")
		if authHeader == "" {
			reportError(c, AuthenticationError, "Empty header")
			c.Abort()
			return
		}
		claims, err := ParseToken(authHeader)
		if err != nil {
			reportError(c, AuthenticationError, "Invalid token "+err.Error())
			c.Abort()
			return
		}
		if time.Now().Unix() > claims.ExpiresAt {
			reportError(c, AuthenticationError, "Token expired")
			c.Abort()
		}
		c.Set("user_id", claims.UserId)
		c.Next()
	}
}
