package procedure

import (
	"context"
	"nh-be/constant"
	"nh-be/internal/procedure"
	"time"

	"github.com/google/uuid"
)

func ContextWithUser(userID uuid.UUID) context.Context {
	return context.WithValue(context.Background(), constant.CtxUserId, userID)
}

func TestProcedureList() []procedure.Procedure {
	return []procedure.Procedure{
		{
			ID:          uuid.New(),
			Title:       "Test Procedure",
			Description: "Test Procedure Description",
			Version:     1,
			ParentID:    nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Title:       "Test Procedure 2",
			Description: "Test Procedure Description 2",
			Version:     2,
			ParentID:    nil,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}
}
