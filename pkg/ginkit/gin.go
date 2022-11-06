package ginkit

import (
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"

	"github.com/amikai/line-go-prac/pkg/logkit"
)

func Default(logger *logkit.Logger) *gin.Engine {
	r := gin.New()
	r.Use(ginzap.Ginzap(logger.Logger, time.RFC3339, true))
	r.Use(ginzap.RecoveryWithZap(logger.Logger, true))
	return r
}
