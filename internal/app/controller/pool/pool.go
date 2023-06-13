//go:generate mockgen -package=mock -destination=../../service/util/testutils/mocks/controller/pool/pool_mock.go uniswapper/internal/app/controller/pool IPoolController

package pool

import (
	"net/http"
	"strings"

	"uniswapper/internal/app/constants"
	"uniswapper/internal/app/controller"
	poolDB "uniswapper/internal/app/db/repository/pool"

	"uniswapper/internal/app/service/correlation"
	"uniswapper/internal/app/service/logger"

	"github.com/gin-gonic/gin"
)

// IPoolController represents the interface for PoolController
type IPoolController interface {
	GetPoolLogsById(c *gin.Context)
	GetPoolLogsHistory(c *gin.Context)
}

// PoolController is the implementation of the IPoolController interface
type PoolController struct {
	PoolDBClient poolDB.IPoolLogsRepository
}

// NewPoolController creates a new instance of PoolController
func NewPoolController(
	poolDBClient poolDB.IPoolLogsRepository,
) IPoolController {
	return &PoolController{
		PoolDBClient: poolDBClient,
	}
}

func (u PoolController) GetPoolLogsById(c *gin.Context) {
	ctx := correlation.WithReqContext(c) // Get the request context
	log := logger.Logger(ctx)            // Get the logger

	poolID := c.Param("pool_id")
	if poolID == "" || len(strings.TrimSpace(poolID)) == 0 {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BadRequest)
		return
	}

	block := c.DefaultQuery("block", "latest")

	logs, err := u.PoolDBClient.GetPoolLogs(ctx, poolID, block)
	if err != nil {
		log.Errorf("Error getting user", err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	// Respond with success and the list of users
	controller.RespondWithSuccess(c, http.StatusAccepted, "Requested Log Info", logs)
}

func (u PoolController) GetPoolLogsHistory(c *gin.Context) {
	ctx := correlation.WithReqContext(c) // Get the request context
	log := logger.Logger(ctx)            // Get the logger

	poolID := c.Param("pool_id")
	if poolID == "" {
		controller.RespondWithError(c, http.StatusBadRequest, constants.BadRequest)
		return
	}

	logs, err := u.PoolDBClient.GetPoolLogsHistory(ctx, poolID)
	if err != nil {
		log.Errorf("Error getting user", err)
		controller.RespondWithError(c, http.StatusInternalServerError, constants.InternalServerError)
		return
	}

	// Respond with success and the list of users
	controller.RespondWithSuccess(c, http.StatusAccepted, "Requested Log Info", logs)
}
