package utils

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

// 需要在配置文件app.conf中设置Tokenexp的值
func CreateToken(userName string) string {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	//添加令牌关键信息
	TokenExp, _ := strconv.Atoi(beego.AppConfig.String("TokenExp"))
	//添加令牌期限，TokenExp变量的值为多少，期限就是多少小时
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(TokenExp)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["userName"] = userName
	token.Claims = claims
	tokenString, _ := token.SignedString([]byte(beego.AppConfig.String("TokenSecrets")))
	return tokenString
}

func CheckToken(tokenString string) string {
	userName := ""
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(beego.AppConfig.String("TokenSecrets")), nil
	})

	if token != nil && token.Valid {
		claims, _ := token.Claims.(jwt.MapClaims)
		userName = claims["userName"].(string)
	}
	return userName
}

// return this result to client then all later request should have header "Authorization: Bearer <token> "
func GetHeaderTokenValue(tokenString string) string {
	//Authorization: Bearer <token>
	return fmt.Sprintf("Bearer %s", tokenString)
}
