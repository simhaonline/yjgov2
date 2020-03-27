package api

import (
	"yj-app/app/controller/api/login"
	"yj-app/app/service/middleware/jwt"
	"yj-app/app/yjgframe/router"
)

func init() {
	group1 := router.New("api", "/v1")
	group1.POST("/login", "", login.Login)
	group2 := router.New("api", "/v1/api", jwt.JWTAuthMiddleware())
	group2.POST("/test", "api", login.Test)
}
