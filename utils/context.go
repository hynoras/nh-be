package utils

import (
	"context"
	"errors"
	"log"
	"nh-be/constant"

	"github.com/google/uuid"
)

func GetUserIdFromContext(ctx context.Context) (uuid.UUID, error) {
	userID, ok := ctx.Value(constant.CtxUserId).(uuid.UUID)
	log.Println("userID from context", userID)
	log.Println("ok", ok)
	if !ok {
		return uuid.Nil, errors.New("user ID not found in context")
	}
	return userID, nil
}
