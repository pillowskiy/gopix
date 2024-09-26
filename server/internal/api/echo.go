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
	"github.com/pillowskiy/gopix/internal/respository/httprepo"
	"github.com/pillowskiy/gopix/internal/respository/postgres"
	"github.com/pillowskiy/gopix/internal/respository/redis"
	"github.com/pillowskiy/gopix/internal/respository/s3"
	"github.com/pillowskiy/gopix/internal/usecase"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/metric"
	"github.com/pillowskiy/gopix/pkg/storage"
	"github.com/pillowskiy/gopix/pkg/token"

	_ "net/http/pprof"
)

type EchoServer struct {
	echo   *echo.Echo
	cfg    *config.Config
	sh     *storage.StorageHolder
	logger logger.Logger
}

func NewEchoServer(cfg *config.Config, sh *storage.StorageHolder, logger logger.Logger) *EchoServer {
	return &EchoServer{echo: echo.New(), cfg: cfg, sh: sh, logger: logger}
}

func (s *EchoServer) Listen() error {
	server := &http.Server{
		Addr:         s.cfg.Server.Addr,
		ReadTimeout:  time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout: time.Second * s.cfg.Server.WriteTimeout,
	}

	if err := s.MapHandlers(); err != nil {
		return err
	}

	s.logger.Infof("Server is listening on ADDR: %s", s.cfg.Server.Addr)
	if err := s.echo.StartServer(server); err != nil {
		return err
	}

	return nil
}

func (s *EchoServer) MapHandlers() error {
	followingRepo := postgres.NewFollowingRepository(s.sh.Postgres)
	followingUC := usecase.NewFollowingUseCase(followingRepo)

	userCache := redis.NewUserCache(s.sh.Redis)
	userRepo := postgres.NewUserRepository(s.sh.Postgres)

	jwtTokenGen := token.NewJWTTokenGenerator(
		s.cfg.Server.Session.Secret,
		s.cfg.Server.Session.Expire*time.Second,
	)
	authUC := usecase.NewAuthUseCase(userRepo, userCache, s.logger, jwtTokenGen)
	userUC := usecase.NewUserUseCase(userRepo, userCache, followingUC, s.logger)

	subscriptionUC := usecase.NewSubscriptionUseCase(followingUC, userUC)

	vecRepo := httprepo.NewVectorizationRepository(s.cfg.VecService.URL)
	imageCache := redis.NewImageCache(s.sh.Redis)
	imageRepo := postgres.NewImageRepository(s.sh.Postgres)
	imageStorage := s3.NewImageStorage(s.sh.S3, s.sh.S3.PublicBucket)
	imageACL := policy.NewImageAccessPolicy()
	imageUC := usecase.NewImageUseCase(
		imageStorage, imageCache, imageRepo, vecRepo, imageACL, s.logger,
	)

	commentRepo := postgres.NewCommentRepository(s.sh.Postgres)
	commentACL := policy.NewCommentAccessPolicy()
	commentUC := usecase.NewCommentUseCase(commentRepo, commentACL, imageUC, s.logger)

	albumRepo := postgres.NewAlbumRepository(s.sh.Postgres)
	albumACL := policy.NewAlbumAccessPolicy()
	albumUC := usecase.NewAlbumUseCase(albumRepo, albumACL, imageUC)

	tagRepo := postgres.NewTagRepository(s.sh.Postgres)
	tagACL := policy.NewTagAccessPolicy()
	tagUC := usecase.NewTagUseCase(tagRepo, tagACL, imageUC)

	metrics, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.Name)
	if err != nil {
		s.logger.Error("CreateMetrics", err.Error())
	}

	s.echo.Use(middlewares.MetricsMiddleware(metrics))
	s.echo.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))
	s.echo.Use(middlewares.CORSMiddleware(s.cfg.Server.CORS))

	v1 := s.echo.Group("/api/v1")
	guardMiddlewares := middlewares.NewGuardMiddlewares(authUC, s.logger, s.cfg.Server.Cookie)

	authGroup := v1.Group("/auth")
	authHandlers := handlers.NewAuthHandlers(authUC, s.logger, s.cfg.Server.Cookie)
	routes.MapAuthRoutes(authGroup, authHandlers, guardMiddlewares)

	userGroup := v1.Group("/users")
	userHandlers := handlers.NewUserHandlers(userUC, s.logger)
	routes.MapUserRoutes(userGroup, userHandlers, guardMiddlewares)

	subscriptionGroup := v1.Group("/subscriptions")
	subscriptionHandlers := handlers.NewSubscriptionHandlers(subscriptionUC, s.logger)
	routes.MapSubscriptionRoutes(subscriptionGroup, subscriptionHandlers, guardMiddlewares)

	imagesGroup := v1.Group("/images")
	imagesHandlers := handlers.NewImageHandlers(imageUC, s.logger)
	routes.MapImageRoutes(imagesGroup, imagesHandlers, guardMiddlewares)

	commentsGroup := imagesGroup.Group("")
	commentsHandlers := handlers.NewCommentHandlers(commentUC, s.logger)
	routes.MapCommentRoutes(commentsGroup, commentsHandlers, guardMiddlewares)

	tagsGroup := imagesGroup.Group("")
	tagsHandlers := handlers.NewTagHandlers(tagUC, s.logger)
	routes.MapTagRoutes(tagsGroup, tagsHandlers, guardMiddlewares)

	albumsGroup := v1.Group("/albums")
	albumsHandlers := handlers.NewAlbumHandlers(albumUC, s.logger)
	routes.MapAlbumRoutes(albumsGroup, albumsHandlers, guardMiddlewares)

	return nil
}
