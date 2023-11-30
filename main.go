package main

import (
	"github.com/Yu-Qi/exchange/api"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()
	registerCommonAPI(r)
	registerV1API(r)

	r.Run(":8040")
}
func registerCommonAPI(r *gin.Engine) {
	r.GET("/ok", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})
}

func registerV1API(r *gin.Engine) {
	v1 := r.Group("/v1")

	// accounts
	account := v1.Group("/accounts")
	account.POST("/register", api.AccountRegister)
	account.POST("/verify", api.AccountVerify)
	account.POST("/login", api.AccountLogin)
	account.POST("/update_password", api.AccountUpdatePassword)
	account.POST("/forget_password", api.AccountForgetPassword)
	account.GET("/me", api.AccountMe)
}
