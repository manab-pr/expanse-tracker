package routes

import (
	handler "expanse-tracker/handlers"

	"github.com/gin-gonic/gin"
)

func EmailRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.POST("/send-otp", handler.SendOTPHandler())
	incomingRoutes.POST("/verify-otp", handler.VerifyOTPHandler())
	incomingRoutes.POST("/forgot-password", handler.ForgotPasswordHandler())
	incomingRoutes.POST("/reset-password", handler.ResetPasswordHandler())

}
