package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"nh-be/internal/config"
	"nh-be/internal/features/permission"
	"nh-be/internal/platform/email"
	"nh-be/internal/platform/mq"
	"nh-be/internal/platform/session"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// Service holds all infrastructure connections and background worker
// lifecycle handles. It is created once at startup by InitializeServices
// and passed through the application for dependency injection and cleanup.
type Service struct {
	DB        *gorm.DB
	SQLDB     *sql.DB
	Redis     *redis.Client
	RabbitMQ  *amqp.Connection
	PubCh     *amqp.Channel
	ConCh     *amqp.Channel
	ConCancel context.CancelFunc
	WG        *sync.WaitGroup
}

// InitializeServices connects to PostgreSQL, Redis, and RabbitMQ, declares
// the message queue topology, and starts background consumers. It returns
// a fully initialized Service or an error if any connection fails.
func InitializeServices(cfg *config.Config) (*Service, error) {
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("postgresql: %w", err)
	}

	rdb, err := config.NewRedisClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}

	conn, err := mq.NewRabbitMQConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: %w", err)
	}

	pubCh, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: failed to open publisher channel: %w", err)
	}

	conCh, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: failed to open consumer channel: %w", err)
	}

	dqErr := mq.DeclareQueues(pubCh, email.SendVerificationEmailQueue)
	if dqErr != nil {
		return nil, fmt.Errorf("rabbitmq: failed to declare queues: %w", dqErr)
	}

	deErr := mq.DeclareExchange(conCh, email.AuthExchangeName)
	if deErr != nil {
		return nil, fmt.Errorf("rabbitmq: failed to declare exchange: %w", deErr)
	}

	bqErr := mq.BindQueue(
		conCh,
		email.SendVerificationEmailQueue,
		email.UserRegisteredRoutingKey,
		email.AuthExchangeName,
	)
	if bqErr != nil {
		return nil, fmt.Errorf("rabbitmq: failed to bind queue: %w", bqErr)
	}

	wg := &sync.WaitGroup{}

	conCtx, conCancel := context.WithCancel(context.Background())
	emailConsumer := email.NewEmailConsumer(conCh, cfg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = emailConsumer.SendVerificationEmail(conCtx)
	}()

	sqlDB, _ := db.DB()

	return &Service{
		DB:        db,
		SQLDB:     sqlDB,
		Redis:     rdb,
		RabbitMQ:  conn,
		PubCh:     pubCh,
		ConCh:     conCh,
		ConCancel: conCancel,
		WG:        wg,
	}, nil
}

// OAuthProviderConfig holds credentials for a specific OAuth provider.
type OAuthProviderConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// SharedDeps holds the shared dependencies that are injected into all
// feature routers. It is created once from Service.NewSharedDeps() and
// ensures a single instance of PermissionService and SessionStore.
type SharedDeps struct {
	DB                *gorm.DB
	Redis             *redis.Client
	PubCh             *amqp.Channel
	SessionStore      session.SessionStore
	PermissionService permission.Service
	FrontendURL       string
	OAuthProviders    map[string]*OAuthProviderConfig
}

// NewSharedDeps constructs the SharedDeps struct by creating singletons
// for SessionStore and PermissionService from the Service's connections.
func (s *Service) NewSharedDeps(cfg *config.Config) *SharedDeps {
	sessionStore := session.NewSessionStore(s.Redis)
	permissionRepo := permission.NewRepository(s.DB)
	permissionCache := permission.NewPermissionCache(s.Redis)
	permissionService := permission.NewService(permissionRepo, permissionCache)

	return &SharedDeps{
		DB:                s.DB,
		Redis:             s.Redis,
		PubCh:             s.PubCh,
		SessionStore:      sessionStore,
		PermissionService: permissionService,
		FrontendURL:       cfg.FrontendURL,
		OAuthProviders: map[string]*OAuthProviderConfig{
			"google": {
				ClientID:     cfg.GoogleClientID,
				ClientSecret: cfg.GoogleClientSecret,
				RedirectURL:  cfg.GoogleRedirectURL,
			},
		},
	}
}

// ShutdownServices performs an orderly teardown of all infrastructure resources.
// It should be called via defer in main() after ListenAndServe returns,
// ensuring the HTTP server has fully drained before connections are closed.
//
// Shutdown order:
//  1. Flush OpenTelemetry traces
//  2. Cancel background consumers and wait for goroutines to finish
//  3. Close RabbitMQ channels and connection
//  4. Close database connection
//  5. Close Redis connection
func (s *Service) ShutdownServices(tracerShutdown func(context.Context) error) {
	if tracerShutdown != nil {
		if err := tracerShutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown tracer", "error", err)
		}
		slog.Info("OpenTelemetry tracer shut down")
	}

	s.ConCancel()
	s.WG.Wait()

	if s.PubCh != nil {
		s.PubCh.Close()
		slog.Info("RabbitMQ publisher channel closed")
	}
	if s.ConCh != nil {
		s.ConCh.Close()
		slog.Info("RabbitMQ consumer channel closed")
	}
	if s.RabbitMQ != nil {
		s.RabbitMQ.Close()
		slog.Info("RabbitMQ connection closed")
	}
	if s.SQLDB != nil {
		s.SQLDB.Close()
		slog.Info("Database connection closed")
	}
	if s.Redis != nil {
		s.Redis.Close()
		slog.Info("Redis connection closed")
	}
}
