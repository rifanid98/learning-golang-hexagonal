package logger

import (
	"bitbucket.org/klopos/majoo-logger/log"
	"bitbucket.org/klopos/majoo-logger/logger"
	"time"

	"learning-golang-hexagonal/utils/config"

	"os"
	"strings"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/spf13/viper"
)

var (
	// MiddlewareLog logger.
	MiddlewareLog *rotatelogs.RotateLogs
	appLog        *rotatelogs.RotateLogs
	// Logger is logger instance.
	//Logger zerolog.Logger
	Logger logger.ILogger
	//Log           *log.Logger
)

func init() {
	logdir := viper.GetString("logdir")
	logMaxAge := viper.GetInt("log_max_age")
	//debug := viper.GetBool("debug")

	// default app log dir setting
	if !strings.HasPrefix(logdir, "/") {
		dir := ""
		if viper.Get("env") == "testing" {
			if viper.Get("env") == nil {
				config.LoadEnvVars()
			}
			dir = viper.GetString("app_path")
		} else {
			dir, _ = os.Getwd()
		}
		logdir = dir + "/log"
	}

	if logMaxAge < 1 {
		// default 15 days
		logMaxAge = 15
	}

	// Set Middleware logging.
	MiddlewareLog, _ = rotatelogs.New(
		logdir+"/access_log.%Y%m%d%H%M",
		rotatelogs.WithLinkName(logdir+"/access_log"),
		rotatelogs.WithMaxAge(time.Duration(logMaxAge)*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	// Set App logging.
	appLog, _ = rotatelogs.New(
		logdir+"/app_log.%Y%m%d%H%M",
		rotatelogs.WithLinkName(logdir+"/app_log"),
		rotatelogs.WithMaxAge(time.Duration(logMaxAge)*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)

	// Set logger, format and level

	// # ZERO LOG
	//zerolog.SetGlobalLevel(zerolog.InfoLevel)
	//if debug {
	//	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	//}
	//zerolog.TimeFieldFormat = "2006-01-02T15:04:05.000000"
	//Logger = zerolog.New(appLog).With().Timestamp().Logger()

	// # MAJOO LOG
	majooLog, err := logger.NewLogger(
		logger.WithAppName(os.Getenv("APP_NAME")),
		logger.WithLevel(logger.Level.DEBUG),
		logger.WithOutput(logger.Output.JSON),
	)
	if err != nil {
		panic(err)
	}
	Logger = majooLog
}

func Debug(message string, options ...log.Option) {
	Logger.Debug(message, options...)
}

func Info(message string, options ...log.Option) {
	Logger.Info(message, options...)
}

func Warn(message string, options ...log.Option) {
	Logger.Warn(message, options...)
}

func Error(message string, options ...log.Option) {
	Logger.Error(message, options...)
}

func Panic(message string, options ...log.Option) {
	Logger.Panic(message, options...)
}
