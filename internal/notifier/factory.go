package notifier

import (
	"log/slog"

	"github.com/giannuccilli/user-api/internal/config"
	"github.com/giannuccilli/user-api/internal/domain"
)

func NewNotifier(cfg *config.Config, logger *slog.Logger, failedEventRepo domain.FailedEventRepository) domain.UserNotifier {
	if cfg.KafkaBrokers == "" {
		return NewNoopNotifier(logger)
	}
	return NewKafkaNotifier(cfg.KafkaBrokers, cfg.KafkaTopic, logger, failedEventRepo)
}
