package main

import (
	"github.com/labstack/echo/v4"
	"github.com/ybkuroki/go-webapp-sample/config"
	"github.com/ybkuroki/go-webapp-sample/logger"
	"github.com/ybkuroki/go-webapp-sample/middleware"
	"github.com/ybkuroki/go-webapp-sample/migration"
	"github.com/ybkuroki/go-webapp-sample/mycontext"
	"github.com/ybkuroki/go-webapp-sample/repository"
	"github.com/ybkuroki/go-webapp-sample/router"
)

func main() {
	e := echo.New()

	conf, env := config.Load()
	logger := logger.NewLogger(env)
	logger.GetZapLogger().Infof("Loaded this configuration : application." + env + ".yml")

	rep := repository.NewBookRepository(logger, conf)
	context := mycontext.NewContext(rep, conf, logger)

	migration.CreateDatabase(context)
	migration.InitMasterData(context)

	router.Init(e, context)
	middleware.InitLoggerMiddleware(e, context)
	middleware.InitSessionMiddleware(e, context)

	if conf.StaticContents.Path != "" {
		e.Static("/", conf.StaticContents.Path)
		logger.GetZapLogger().Infof("Served the static contents. path: " + conf.StaticContents.Path)
	}

	if err := e.Start(":8080"); err != nil {
		logger.GetZapLogger().Errorf(err.Error())
	}

	defer rep.Close()
}
