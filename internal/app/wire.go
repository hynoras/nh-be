package app

import (
	"context"
	"database/sql"
	"nh-be/config"
	"nh-be/internal/email"
	"nh-be/internal/features/auth"
	"nh-be/internal/features/experiment"
	"nh-be/internal/features/experiment/result"
	"nh-be/internal/features/permission"
	"nh-be/internal/features/user"
	"nh-be/mq"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// InitializeServices initializes all the services and returns the database, redis, and rabbitmq publisher channel
func InitializeServices() (*gorm.DB, *sql.DB, *redis.Client, *amqp.Connection, *amqp.Channel, *amqp.Channel, context.CancelFunc, error) {
	db := config.ConnectDatabase()
	db.AutoMigrate(
		&auth.VerificationToken{},
		&user.User{},
		&permission.Permission{},
		&permission.PermissionGroup{},
		&user.UserPermission{},
		&experiment.Experiment{},
		&result.ExperimentResult{},
	)

	rdb := config.NewRedisClient()

	conn, err := mq.NewRabbitMQConnection()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	pubCh, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	conCh, err := conn.Channel()
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, err
	}

	dqErr := mq.DeclareQueues(pubCh, email.SendVerificationEmailQueue)
	if dqErr != nil {
		return nil, nil, nil, nil, nil, nil, nil, dqErr
	}

	deErr := mq.DeclareExchange(conCh, email.AuthExchangeName)
	if deErr != nil {
		return nil, nil, nil, nil, nil, nil, nil, deErr
	}

	bqErr := mq.BindQueue(
		conCh,
		email.SendVerificationEmailQueue,
		email.UserRegisteredRoutingKey,
		email.AuthExchangeName,
	)
	if bqErr != nil {
		return nil, nil, nil, nil, nil, nil, nil, bqErr
	}

	conCtx, conCancel := context.WithCancel(context.Background())
	emailConsumer := email.NewEmailConsumer(conCh)
	go emailConsumer.SendVerificationEmail(conCtx)

	sqlDB, _ := db.DB()

	return db, sqlDB, rdb, conn, pubCh, conCh, conCancel, nil
}
