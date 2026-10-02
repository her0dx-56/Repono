package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
	
)

type Claims struct{
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type JWTManager struct{
	secretKey []byte
}

func NewJWTManager(secret string) *JWTManager {
	return &JWTManager{
		secretKey: []byte(secret),
	}
}

func(m *JWTManager) GenerateToken(userID string) (string,error){
	now:=time.Now()

	claims:= Claims{
		UserID:userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(24*time.Hour)),
		},
	}
	token:=jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)
	return token.SignedString(m.secretKey)
}

func(m *JWTManager) ValidateToken(tokenString string) (*Claims, error){
	token,err := jwt.ParseWithClaims(
		tokenString,
		&Claims{},
		func(token *jwt.Token) (interface{},error) {
			if token.Method!= jwt.SigningMethodHS256{
				return nil,jwt.ErrTokenSignatureInvalid
			}
			return m.secretKey,nil
		},
	)
	if err!=nil{
		return nil,err
	}
	if !token.Valid{
		return nil,jwt.ErrTokenInvalidClaims
	}

	claims,ok:=token.Claims.(*Claims)
	if !ok{
		return nil,jwt.ErrTokenInvalidClaims
	}
	return claims,nil
}