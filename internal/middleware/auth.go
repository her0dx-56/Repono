package middleware

import (
	"context"
	"go_dfs/internal/auth"
	"net/http"
	"strings"
)
type contextKey string
const userIDKey contextKey = "user_id"
func Auth(jwtManager *auth.JWTManager) func(http.Handler) http.Handler{
	return func(next http.Handler) http.Handler{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			authHeader:=r.Header.Get("Authorization")
			parts:= strings.SplitN(authHeader," ",2)
			if len(parts)!=2 || parts[0]!="Bearer"{
				http.Error(w,"unauthorized",http.StatusUnauthorized)
				return
			} 
			tokenString:=parts[1]
			claims,err:=jwtManager.ValidateToken(tokenString)
			if err!=nil{
				http.Error(w,"unauthorized",http.StatusUnauthorized)
				return
			}
			
			ctx:=context.WithValue(r.Context(),userIDKey,claims.UserID)
			r=r.WithContext(ctx)
			next.ServeHTTP(w,r)

		})
	}
}

func GetUserID(r *http.Request) (string,bool) {
	userID,ok:=r.Context().Value(userIDKey).(string)
	return userID,ok
}