package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jimmy/server/header"
	"github.com/jimmy/server/jwt"
	"github.com/jimmy/server/model"
	"log"
	"net/http"
	"time"
)

type LoginDatabase interface {
	GetUserAuthByBasic(userName, password string) *model.UserAuth
}

type LoginAPI struct {
	DB LoginDatabase
}

func (u *LoginAPI) Login(c *gin.Context) {
	//c.JSON(http.StatusOK, map[string]string{ "status": "Logged in." })
	type reqAuth struct {
		UserName string
		Password string
	}
	header.HeaderWrite(c)

	var req reqAuth
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		log.Println("Auth json decode failed: ", err)
	}
	fmt.Println(req.UserName, req.Password, "login.")

	ua := u.DB.GetUserAuthByBasic(req.UserName, req.Password)
	/****** if query result just have one match *********/
	if ua != nil {
		j := jwt.UserJwt{
			Token: ua.Token,
			Exp: time.Now(),
		}
		if jwtString, ok := j.JwtSignedString(); ok {
			c.JSON(http.StatusOK, map[string]string{"token": jwtString})
		}
	}else{
		/****** Fail Auth, return null and ok status *********/
		log.Println("Auth failed: ", req)
		c.JSON(http.StatusUnauthorized, map[string]string{"error": "Wrong Username or Password."})
	}
}