package consumer

import (
	"context"
	"log/slog"

	"github.com/smexpress/pkg/messaging"
	"github.com/smexpress/services/user-service/internal/domain/entity"
	"github.com/smexpress/services/user-service/internal/usecase"
)

type UserCreatedEvent struct {
	UserID      string `json:"user_id"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	CountryCode string `json:"country_code"`
	FranchiseID string `json:"franchise_id"`
}

func HandleUserCreated(profileUC *usecase.ProfileUseCase, logger *slog.Logger) messaging.Handler {
	return func(ctx context.Context, data []byte) error {
		event, err := messaging.DecodeEvent[UserCreatedEvent](data)
		if err != nil {
			logger.Error("failed to decode user created event", slog.String("error", err.Error()))
			return err
		}

		profile := &entity.UserProfile{
			UserID:      event.UserID,
			CountryCode: event.CountryCode,
			Locale:      "en",
		}

		if err := profileUC.Create(ctx, profile); err != nil {
			logger.Error("failed to create profile from event",
				slog.String("user_id", event.UserID),
				slog.String("error", err.Error()))
			return err
		}

		logger.Info("profile created from event", slog.String("user_id", event.UserID))
		return nil
	}
}
