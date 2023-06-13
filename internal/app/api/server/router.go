package server

import (
	"context"
	"strings"
	"time"
	"uniswapper/internal/app/constants"
	"uniswapper/internal/app/controller/healthcheck"
	poolController "uniswapper/internal/app/controller/pool"
	"uniswapper/internal/app/db"
	poolDBClient "uniswapper/internal/app/db/repository/pool"
	"uniswapper/internal/app/service/logger"
	uniswapv3_pool "uniswapper/internal/app/service/pool"

	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Init(ctx context.Context, dbConnection *db.DBService) *gin.Engine {
	if strings.EqualFold(constants.Config.Environment, "prod") {
		gin.SetMode(gin.ReleaseMode)
	}
	return NewRouter(ctx, dbConnection)

}
func addCSPHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	}
}

func addReferrerPolicyHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Next()
	}
}

func addPermissionsPolicyHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Permissions-Policy", "default-src 'none'")
		c.Next()
	}
}

func addFeaturePolicyHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Feature-Policy", "none")
		c.Next()
	}
}

func NewRouter(ctx context.Context, dbConnection *db.DBService) *gin.Engine {
	log := logger.Logger(ctx)

	log.Info("setting up service and controllers")

	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(helmet.Default())
	//Content-Security-Policy
	router.Use(addCSPHeader())
	//Referrer-Policy
	router.Use(addReferrerPolicyHeader())
	//Permissions-Policy
	router.Use(addPermissionsPolicyHeader())
	//Feature-Policy
	router.Use(addFeaturePolicyHeader())

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Accept", "Content-Type", constants.AUTHORIZATION, constants.CORRELATION_KEY_ID.String()},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(uuidInjectionMiddleware())

	// DB Clients
	var (
		poolDBClient = poolDBClient.NewPoolLogsRepository(dbConnection)
	)

	//Service
	var (
		uniswapV3Pool = uniswapv3_pool.NewUniswapV3Pool(ctx, poolDBClient)
	)

	// Start Uniswap V3 Pool to store Logs
	go uniswapV3Pool.RunUniswapV3Pool(ctx)

	// Controller
	var (
		poolController        = poolController.NewPoolController(poolDBClient)
		healthCheckController = healthcheck.NewHealthCheckController()
	)

	v1 := router.Group("/v1/api/pool")
	{
		v1.GET(HEALTH_CHECK, healthCheckController.HealthCheck)

		v1.GET(POOL_LOG_BY_ID, poolController.GetPoolLogsById)
		v1.GET(POOL_HISTORY_LOG, poolController.GetPoolLogsHistory)
	}

	return router
}

// uuidInjectionMiddleware injects the request context with a correlation id of type uuid
func uuidInjectionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		correlationId := c.GetHeader(string(constants.CORRELATION_KEY_ID))
		if len(correlationId) == 0 {
			correlationID, _ := uuid.NewUUID()
			correlationId = correlationID.String()
			c.Request.Header.Set(constants.CORRELATION_KEY_ID.String(), correlationId)
		}
		c.Writer.Header().Set(constants.CORRELATION_KEY_ID.String(), correlationId)

		c.Next()
	}
}
