//go:build wireinject
// +build wireinject

package main

import (
	"collapp/configs"
	"collapp/infras"
	"collapp/transport/http/middleware"
	"github.com/google/wire"

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
	httpRouter "collapp/transport/http/router"
)

// Wiring for configurations.
var configurations = wire.NewSet(
	configs.Get,
)

// Wiring for database.
var database = wire.NewSet(
	infras.NewMysqlDB,
)

var translationModule = wire.NewSet(
	// TranslationRepository interface and implementation
	translationRepo.NewTranslationRepository,

	// TranslationService interface and implementation
	translationService.NewTranslationService,
)

var userModule = wire.NewSet(
	// UserRepository interface and implementation
	userRepo.NewUserRepository,

	// UserService interface and implementation
	userService.NewUserService,
)

var langModule = wire.NewSet(
	// LangRepository interface and implementation
	langRepo.NewLangRepository,

	// LangService interface and implementation
	langService.NewLangService,
)

var modules = wire.NewSet(
	translationModule,
	userModule,
	langModule,
)

var authMiddleware = wire.NewSet(
	middleware.NewAuthMiddleware,
)

// Wiring for HTTP routing.
var routing = wire.NewSet(
	translationHandler.NewTranslationHandler,
	userHandler.NewUserHandler,
	langHandler.NewLangHandler,
	wire.Struct(new(httpRouter.ModuleHandlers), "TranslationHandler", "UserHandler", "LangHandler"),

	httpRouter.NewRouter,
)

var httpServer = wire.NewSet(
	httpTransport.NewHTTP,
)

func InitializeService() *httpTransport.HTTP {
	wire.Build(
		// configurations
		configurations,
		// persistences
		database,
		// middleware
		authMiddleware,
		// domains
		modules,
		// routing
		routing,
		// http
		httpServer,
	)
	return &httpTransport.HTTP{}
}
