package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ServerInterface interface {
	ListHosts(ctx *gin.Context)
	CreateHost(ctx *gin.Context)
	GetHost(ctx *gin.Context, id int)
	UpdateHost(ctx *gin.Context, id int)
	DeleteHost(ctx *gin.Context, id int)
	ListPackages(ctx *gin.Context)
	CreatePackage(ctx *gin.Context)
	GetPackage(ctx *gin.Context, id int)
	DeletePackage(ctx *gin.Context, id int)
	ListDeployments(ctx *gin.Context)
	CreateDeployment(ctx *gin.Context)
	GetDeployment(ctx *gin.Context, id int)
	DeleteDeployment(ctx *gin.Context, id int)
	RunDeployment(ctx *gin.Context, id int)
	GetDeploymentLogs(ctx *gin.Context, id int)
	HealthCheck(ctx *gin.Context)
}

func RegisterHandlers(router *gin.Engine, si ServerInterface) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", si.HealthCheck)

		hosts := v1.Group("/hosts")
		{
			hosts.GET("", si.ListHosts)
			hosts.POST("", si.CreateHost)
			hosts.GET("/:id", si.GetHost)
			hosts.PUT("/:id", si.UpdateHost)
			hosts.DELETE("/:id", si.DeleteHost)
		}

		packages := v1.Group("/packages")
		{
			packages.GET("", si.ListPackages)
			packages.POST("", si.CreatePackage)
			packages.GET("/:id", si.GetPackage)
			packages.DELETE("/:id", si.DeletePackage)
		}

		deployments := v1.Group("/deployments")
		{
			deployments.GET("", si.ListDeployments)
			deployments.POST("", si.CreateDeployment)
			deployments.GET("/:id", si.GetDeployment)
			deployments.DELETE("/:id", si.DeleteDeployment)
			deployments.POST("/:id/run", si.RunDeployment)
			deployments.GET("/:id/logs", si.GetDeploymentLogs)
		}
	}
}

type Error struct {
	Message string `json:"message"`
}

func ErrorJSON(ctx *gin.Context, code int, message string) {
	ctx.JSON(code, Error{Message: message})
}

func GetIDParam(ctx *gin.Context) (int, bool) {
	var id int
	if _, err := fmt.Sscanf(ctx.Param("id"), "%d", &id); err != nil {
		return 0, false
	}
	return id, true
}
