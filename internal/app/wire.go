package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"nh-be/internal/config"
	"nh-be/internal/features/auth"
	"nh-be/internal/features/experiment"
	"nh-be/internal/features/experiment/result"
	"nh-be/internal/features/permission"
	"nh-be/internal/features/user"
	"nh-be/internal/platform/email"
	"nh-be/internal/platform/mq"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

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

// InitializeServices initializes all the services and returns the database, redis, and rabbitmq publisher channel
func InitializeServices(cfg *config.Config) (*Service, error) {
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		return nil, fmt.Errorf("postgresql: %w", err)
	}
	if cfg.AppEnv == "dev" {
		if err := db.AutoMigrate(
			&auth.VerificationToken{},
			&user.User{},
			&permission.Permission{},
			&permission.PermissionGroup{},
			&user.UserPermission{},
			&experiment.Experiment{},
			&result.ExperimentResult{},
		); err != nil {
			return nil, fmt.Errorf("automigrate failed: %w", err)
		}
		slog.Info("Running AutoMigrate in dev mode")
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
