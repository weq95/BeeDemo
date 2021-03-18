package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/scrypt"
	"io"
	"time"
)

const (
	SecretKey              = "JWT-Secret-key"
	DEFAULT_EXPIRE_SECONDS = 600 //过期时间
	PasswordHashBytes      = 20  //hash 加密长度
)

type MyCustomClaims struct {
	UserID int `json:"userID"`
	jwt.StandardClaims
}

type JwtPayload struct {
	UserName  string `json:"user_name"`
	UserId    int    `json:"user_id"`
	IssueAt   int64  `json:"issue_at"`
	ExpiresAt int64  `json:"expires_at"`
}

//token过期时间
func expireAt() int64 {
	return time.Now().Add(time.Second * time.Duration(DEFAULT_EXPIRE_SECONDS)).Unix()
}

//创建token
func GenerateToken(uname string, uid int) (string, error) {
	// Create the token using your claims
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		MyCustomClaims{
			uid,
			jwt.StandardClaims{
				Issuer:    uname,
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: expireAt(),
			},
		},
	)

	// Signs the token with a secret
	return token.SignedString([]byte(SecretKey))
}

//验证token
func ValidateToken(tokenStr string) (*JwtPayload, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok || !token.Valid {
		return nil, errors.New("failed to validate token ")
	}

	return &JwtPayload{
		UserName:  claims.StandardClaims.Issuer,
		UserId:    claims.UserID,
		IssueAt:   claims.StandardClaims.IssuedAt,
		ExpiresAt: claims.StandardClaims.ExpiresAt,
	}, nil
}

//刷新token
func RefreshToken(tokenStr string) (string, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&MyCustomClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	)

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok || !token.Valid {
		return "", errors.New("refresh token failed")
	}

	newToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		MyCustomClaims{
			claims.UserID,
			jwt.StandardClaims{
				Issuer:    claims.StandardClaims.Issuer,
				IssuedAt:  time.Now().Unix(),
				ExpiresAt: expireAt(),
			},
		},
	)

	return newToken.SignedString([]byte(SecretKey))
}

//生成盐
func GenerateSalt() (string, error) {
	buf := make([]byte, PasswordHashBytes)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		return "", errors.New("error: failed to generate user's salt")
	}

	return fmt.Sprintf("%x", buf), nil
}

//密码hash处理
func GeneratePassHash(pwd, salt string) (string, error) {
	h, err := scrypt.Key([]byte(pwd), []byte(salt),
		16384, 8, 1, PasswordHashBytes)

	if err != nil {
		return "", errors.New("error: failed to generate password hash")
	}

	return fmt.Sprintf("%x", h), nil
}

//获取加密字符串
func GetPwdString(pwd string) (pwdStr, saltKey string, err error) {
	//获取盐值
	saltKey, err = GenerateSalt()
	if err != nil {
		return "", "", err
	}

	//加密密码
	pwdStr, err = GeneratePassHash(pwd, saltKey)

	return pwdStr, saltKey, err
}
