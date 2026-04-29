package router

import (
	"encoding/json"

	dbpkg "backend/internal/infra/db"
	"backend/internal/infra/web/gin/handler"
	"backend/internal/infra/web/gin/middleware"
	"backend/internal/infra/web/gin/presenter"
	"backend/internal/usecase"
	"backend/pkg"
	"backend/pkg/logger"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/gin-gonic/gin"

	ginMiddleware "github.com/oapi-codegen/gin-middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"
	"gorm.io/gorm"
)

func setupSwagger(router *gin.Engine) (*openapi3.T, error) {
	swagger, err := presenter.GetSwagger()
	if err != nil {
		return nil, err
	}

	env := pkg.GetEnvDefault("APP_ENV", "development")
	if env == "development" {
		swaggerJson, _ := json.Marshal(swagger)
		var SwaggerInfo = &swag.Spec{
			InfoInstanceName: "swagger",
			SwaggerTemplate:  string(swaggerJson),
		}
		swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	return swagger, nil
}

func NewGinRouter(db *gorm.DB, corsAllowOrigins []string) (*gin.Engine, error) {
	router := gin.Default()

	router.Use(middleware.CorsMiddleware(corsAllowOrigins))
	swagger, err := setupSwagger(router)
	if err != nil {
		logger.Warn(err.Error())
		return nil, err
	}

	router.Use(middleware.GinZap())
	router.Use(middleware.RecoveryWithZap())

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/", handler.Index)
		apiGroup.GET("/health", handler.Health)
		v1 := apiGroup.Group("/v1")
		{

			csrfHandler := handler.NewCsrfHandler()

			base := dbpkg.NewBaseRepository(db)
			txManager := dbpkg.NewTxManager(db)

			userRepository := dbpkg.NewUserRepository(base)
			userUseCase := usecase.NewUserUsecase(txManager, userRepository)
			userHandler := handler.NewUserHandler(userUseCase)

			taskRepository := dbpkg.NewTaskRepository(base)
			statusRepository := dbpkg.NewStatusRepository(base)
			taskUseCase := usecase.NewTaskUsecase(txManager, taskRepository, statusRepository)
			taskHandler := handler.NewTaskHandler(taskUseCase)

			serverHandler := handler.NewHandler().
				Register(csrfHandler).
				Register(userHandler).
				Register(taskHandler)

			wrapper := presenter.ServerInterfaceWrapper{
				Handler: serverHandler,
			}

			csrfTokenGenerator := middleware.CsrfTokenGenerator()
			v1.GET("/csrf", csrfTokenGenerator, wrapper.GetCsrfToken)

			useCsrf := v1.Group("")
			{

				// useCsrfではCSRF検証->OAPIバリデータ
				useCsrf.Use(middleware.CsrfValidator())
				useCsrf.Use(ginMiddleware.OapiRequestValidator(swagger))

				useCsrf.POST("/signup", wrapper.PostSignUp)
				useCsrf.POST("/login", wrapper.PostLogin)
				useCsrf.POST("/logout", wrapper.PostLogout)

				useJwt := useCsrf.Group("")
				{
					// useJwtではCSRF検証->OAPIバリデータ->JWT認証
					// 処理が軽いものからすることで負荷を軽減
					useJwt.Use(middleware.JwtAuthMiddleware())

					useJwt.POST("/tasks", wrapper.CreateTask)
					useJwt.GET("/tasks/:id", wrapper.GetTaskById)
					useJwt.GET("/tasks", wrapper.GetAllTasks)
					useJwt.PATCH("/tasks/:id", wrapper.UpdateTaskById)
					useJwt.DELETE("/tasks/:id", wrapper.DeleteTaskById)
				}
			}
		}
	}
	return router, err
}
