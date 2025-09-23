package log

import (
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)


var Logger *slog.Logger

func init(){
	Logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

type LoggerKeyType struct{}

func With(c *gin.Context) *gin.Context {
	c.Set(LoggerKeyType{}, Logger)
	return c
}

func From(c *gin.Context) *slog.Logger {
	return c.Value(LoggerKeyType{}).(*slog.Logger)
}

func InfoF(msg string, args ...any) {
	Logger.Info(msg, args...)
}