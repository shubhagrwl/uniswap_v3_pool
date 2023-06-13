package constants

import "uniswapper/internal/config"

var Config *config.ServiceConfig

const (
	//Header constants
	AUTHORIZATION      = "Authorization"
	BEARER             = "Bearer "
	CTK_CLAIM_KEY      = CONTEXT_KEY("claims")
	CORRELATION_KEY_ID = CORRELATION_KEY("X-Correlation-ID")
)

type (
	ENVIRONMENT     string
	CONTEXT_KEY     string
	CORRELATION_KEY string
)

func (c CONTEXT_KEY) String() string {
	return string(c)
}

func (c CORRELATION_KEY) String() string {
	return string(c)
}
