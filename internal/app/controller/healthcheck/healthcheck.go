package healthcheck

import (
	"net/http"
	"uniswapper/internal/app/constants"
	"uniswapper/internal/app/controller"

	"github.com/gin-gonic/gin"
)

type IHealthCheckController interface {
	HealthCheck(c *gin.Context)
}

type HealthCheckController struct {
}

func NewHealthCheckController() IHealthCheckController {
	return &HealthCheckController{}
}

func (h *HealthCheckController) HealthCheck(c *gin.Context) {
	controller.RespondWithSuccess(c, http.StatusOK, "version", gin.H{"version": constants.Config.ProjectVersion})
}
