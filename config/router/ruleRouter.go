package router

import (
	"engine/app/controller"
	"github.com/gin-gonic/gin"
)

func ruleRouterInit(r *gin.RouterGroup) {
	rule := r.Group("/rule")
	rule.GET("/get", controller.GetScore)
}
