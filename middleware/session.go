package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"os"
)

// Session 初始化session
func Session(secret string) gin.HandlerFunc {
	redisAddr := os.Getenv("REDIS_ADDR")
	store, _ := redis.NewStore(10, "tcp", redisAddr, os.Getenv("REDIS_PW"), []byte(secret))
	//Also set Secure: true if using SSL, you should though
	store.Options(sessions.Options{HttpOnly: true, MaxAge: 900, Path: "/"})
	return sessions.Sessions("camp-seesion", store)
}
