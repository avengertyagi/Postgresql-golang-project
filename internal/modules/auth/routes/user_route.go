package routes

import (
	usercontroller "github.com/akshit_tyagi/postgresql_project/internal/modules/auth/controllers/user"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.RouterGroup) {
	auth := r.Group("auth")
	auth.POST("/register", usercontroller.Register)
	auth.POST("/verify-otp", usercontroller.VerifyOTP)
}
