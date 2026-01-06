package serviceinterface

import (
	"challenge-app/internal/domain/model"
	"context"
)

type Notifier interface {
	Push(ctx context.Context, userID uint, n model.Notification) error
}
