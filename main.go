package main

import (
	"collapp/configs"
	"collapp/infras"
	langHandler "collapp/module/lang/handler"
	langRepo "collapp/module/lang/repository"
	langService "collapp/module/lang/service"
	translationHandler "collapp/module/translation/handler"
	translationRepo "collapp/module/translation/repository"
	translationService "collapp/module/translation/service"
	userHandler "collapp/module/user/handler"
	userRepo "collapp/module/user/repository"
	userService "collapp/module/user/service"

	httpTransport "collapp/transport/http"
	"collapp/transport/http/middleware"
	httpRouter "collapp/transport/http/router"

	_ "github.com/go-sql-driver/mysql"
)

var config *configs.Config

func main() {
	// Initialize config
	config = configs.Get()

	db := infras.NewMysqlDB(config)

	translationRepository := translationRepo.NewTranslationRepository(db)
	translationService := translationService.NewTranslationService(db, translationRepository)
	translationHandler := translationHandler.NewTranslationHandler(db, config, translationService)

	userRepository := userRepo.NewUserRepository(db)
	userService := userService.NewUserService(db, userRepository)
	userHandler := userHandler.NewUserHandler(db, config, userService, translationService)

	langRepository := langRepo.NewLangRepository(db)
	langService := langService.NewLangService(db, langRepository)
	langHandler := langHandler.NewLangHandler(db, config, langService, translationService)

	moduleHandlers := httpRouter.ModuleHandlers{
		TranslationHandler: translationHandler,
		UserHandler:        userHandler,
		LangHandler:        langHandler,
	}
	router := httpRouter.NewRouter(moduleHandlers)

	authMiddleware := middleware.NewAuthMiddleware(config, translationService, userService)
	httpServer := httpTransport.NewHTTP(config, router, authMiddleware)
	httpServer.SetupAndServe()
}
