package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	"github.com/pillowskiy/gopix/internal/delivery/rest/middlewares"
	"github.com/pillowskiy/gopix/internal/delivery/rest/routes"
	"github.com/pillowskiy/gopix/internal/policy"
	"github.com/pillowskiy/gopix/internal/respository/postgres"
	"github.com/pillowskiy/gopix/internal/respository/redis"
	"github.com/pillowskiy/gopix/internal/respository/s3"
	"github.com/pillowskiy/gopix/internal/usecase"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/storage"
	"github.com/pillowskiy/gopix/pkg/token"

	_ "net/http/pprof"
)

type EchoServer struct {
	echo   *echo.Echo
	cfg    *config.Server
	sh     *storage.StorageHolder
	logger logger.Logger
}

func NewEchoServer(cfg *config.Server, sh *storage.StorageHolder, logger logger.Logger) *EchoServer {
	return &EchoServer{echo: echo.New(), cfg: cfg, sh: sh, logger: logger}
}

func (s *EchoServer) Listen() error {
	server := &http.Server{
		Addr:         s.cfg.Addr,
		ReadTimeout:  time.Second * s.cfg.ReadTimeout,
		WriteTimeout: time.Second * s.cfg.WriteTimeout,
	}

	if err := s.MapHandlers(); err != nil {
		return err
	}

	s.logger.Infof("Server is listening on ADDR: %s", s.cfg.Addr)
	if err := s.echo.StartServer(server); err != nil {
		return err
	}

	return nil
}

func (s *EchoServer) MapHandlers() error {
	userCache := redis.NewUserCache(s.sh.Redis)
	userRepo := postgres.NewUserRepository(s.sh.Postgres)

	jwtTokenGen := token.NewJWTTokenGenerator(
		s.cfg.Session.Secret,
		s.cfg.Session.Expire*time.Second,
	)
	authUC := usecase.NewAuthUseCase(userRepo, userCache, s.logger, jwtTokenGen)
	userUC := usecase.NewUserUseCase(userRepo, userCache, s.logger)

	imageCache := redis.NewImageCache(s.sh.Redis)
	imageRepo := postgres.NewImageRepository(s.sh.Postgres)
	imageStorage := s3.NewImageStorage(s.sh.S3, s.sh.S3.PublicBucket)
	imageACL := policy.NewImageAccessPolicy()
	imageUC := usecase.NewImageUseCase(imageStorage, imageCache, imageRepo, imageACL, s.logger)

	commentRepo := postgres.NewCommentRepository(s.sh.Postgres)
	commentACL := policy.NewCommentAccessPolicy()
	commentUC := usecase.NewCommentUseCase(commentRepo, commentACL, imageUC, s.logger)

	s.echo.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))
	s.echo.Use(middlewares.CORSMiddleware(s.cfg.CORS))

	v1 := s.echo.Group("/api/v1")
	guardMiddlewares := middlewares.NewGuardMiddlewares(authUC, s.logger, s.cfg.Cookie)

	authGroup := v1.Group("/auth")
	authHandlers := handlers.NewAuthHandlers(authUC, s.logger, s.cfg.Cookie)
	routes.MapAuthRoutes(authGroup, authHandlers, guardMiddlewares)

	userGroup := v1.Group("/users")
	userHandlers := handlers.NewUserHandlers(userUC, s.logger)
	routes.MapUserRoutes(userGroup, userHandlers, guardMiddlewares)

	imagesGroup := v1.Group("/images")
	imagesHandlers := handlers.NewImageHandlers(imageUC, s.logger)
	routes.MapImageRoutes(imagesGroup, imagesHandlers, guardMiddlewares)

	commentsGroup := imagesGroup.Group("")
	commentsHandlers := handlers.NewCommentHandlers(commentUC, s.logger)
	routes.MapCommentRoutes(commentsGroup, commentsHandlers, guardMiddlewares)

	albumsGroup := v1.Group("/albums")
	albumsHandlers := handlers.NewAlbumHandlers(albumUC, s.logger)
	routes.MapAlbumRoutes(albumsGroup, albumsHandlers, guardMiddlewares)

	return nil
}
