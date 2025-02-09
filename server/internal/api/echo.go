package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/pillowskiy/gopix/internal/config"
	"github.com/pillowskiy/gopix/internal/delivery/rest/handlers"
	"github.com/pillowskiy/gopix/internal/delivery/rest/middlewares"
	"github.com/pillowskiy/gopix/internal/delivery/rest/routes"
	"github.com/pillowskiy/gopix/internal/domain"
	"github.com/pillowskiy/gopix/internal/infrastructure/features"
	"github.com/pillowskiy/gopix/internal/infrastructure/oauth"
	"github.com/pillowskiy/gopix/internal/policy"
	"github.com/pillowskiy/gopix/internal/repository/httprepo"
	"github.com/pillowskiy/gopix/internal/repository/postgres"
	"github.com/pillowskiy/gopix/internal/repository/redis"
	"github.com/pillowskiy/gopix/internal/repository/s3"
	"github.com/pillowskiy/gopix/internal/usecase"
	"github.com/pillowskiy/gopix/pkg/logger"
	"github.com/pillowskiy/gopix/pkg/metric"
	"github.com/pillowskiy/gopix/pkg/signal"
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
	s.prepareMiddlewares()

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

	notifRepo := postgres.NewNotificationRepository(s.sh.Postgres)
	notifUC := usecase.NewNotificationUseCase(notifRepo, signal.NewSignal[*domain.Notification](), s.logger)

	oauthRepo := postgres.NewOAuthRepository(s.sh.Postgres)
	oauthClient := oauth.NewOAuthClient(&s.cfg.OAuth)
	oauthUC := usecase.NewOAuthUseCase(oauthRepo, authUC, oauthClient)

	subscriptionUC := usecase.NewSubscriptionUseCase(followingUC, userUC)

	vecRepo := httprepo.NewVectorizationRepository(s.cfg.VecService.URL)
	featExtractor := features.NewBasicFeatureExtractor()
	imagePropsRepo := postgres.NewImagePropsRepository(s.sh.Postgres)
	imageFeatUC := usecase.NewImageFeaturesUseCase(vecRepo, imagePropsRepo, featExtractor, s.logger)

	imageCache := redis.NewImageCache(s.sh.Redis)
	imageRepo := postgres.NewImageRepository(s.sh.Postgres)
	imageStorage := s3.NewImageStorage(s.sh.S3, s.sh.S3.PublicBucket)
	imageACL := policy.NewImageAccessPolicy()
	imageUC := usecase.NewImageUseCase(
		imageStorage, imageCache, imageRepo, imageFeatUC, imageACL, notifUC, s.logger,
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

	v1 := s.echo.Group("/api/v1")
	guardMiddlewares := middlewares.NewGuardMiddlewares(authUC, s.logger, s.cfg.Server.Cookie)

	authGroup := v1.Group("/auth")
	authHandlers := handlers.NewAuthHandlers(authUC, s.logger, s.cfg.Server.Cookie)
	routes.MapAuthRoutes(authGroup, authHandlers, guardMiddlewares)

	oauthHandlers := handlers.NewOAuthHandlers(oauthUC, s.cfg.Server.Cookie, s.logger)
	routes.MapOAuthRoutes(authGroup, oauthHandlers)

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

	notifGroup := v1.Group("/notifications")
	notifHandlers := handlers.NewNotificationHandlers(notifUC, s.logger)
	routes.MapNotificationRoutes(notifGroup, notifHandlers, guardMiddlewares)
	s.echo.File("/", "./index.html")

	return nil
}

func (s *EchoServer) prepareMiddlewares() {
	s.echo.Use(middlewares.CORSMiddleware(s.cfg.Server.CORS))

	if s.cfg.Server.Mode == config.ProductionMode {
		metrics, err := metric.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.Name)
		if err != nil {
			s.logger.Error("CreateMetrics", err.Error())
		}

		s.echo.Use(middlewares.MetricsMiddleware(metrics))
	}

	if s.cfg.Server.Mode == config.DevelopmentMode {
		s.echo.Use(middlewares.MetricsLoggingMiddleware(s.logger))
		s.echo.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))
	}
}
