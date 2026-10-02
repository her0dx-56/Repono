package handler

import(
	"net/http"
	"go_dfs/internal/middleware"
	"fmt"
)

func AuthTest(w http.ResponseWriter, r *http.Request){
	userID,ok:=middleware.GetUserID(r)
	if !ok{
		http.Error(w,"user not found",http.StatusInternalServerError)
	}
	
	fmt.Println(w,"Authenticated user:",userID)
}