package persistence

import (
	"fmt"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/log"
	"github.com/klwxsrx/budget-tracker/pkg/common/app/realtime"
)

type ScheduledRealtimeClient struct {
	service           realtime.Client
	scheduledMessages map[string][][]byte
	logger            log.Logger
}

func (s *ScheduledRealtimeClient) PublishMessage(channel string, data []byte) error {
	s.scheduledMessages[channel] = append(s.scheduledMessages[channel], data)
	return nil
}

func (s *ScheduledRealtimeClient) Rollback() {
	s.reset()
}

func (s *ScheduledRealtimeClient) Commit() {
	for channel, messages := range s.scheduledMessages {
		for _, messageData := range messages {
			err := s.service.PublishMessage(channel, messageData)
			if err != nil {
				s.logger.WithError(err).Error(fmt.Sprintf("failed to publish message to channel %s", channel))
			}
		}
	}
	s.reset()
}

func (s *ScheduledRealtimeClient) reset() {
	s.scheduledMessages = make(map[string][][]byte)
}

func NewScheduledRealtimeClient(client realtime.Client, logger log.Logger) *ScheduledRealtimeClient {
	return &ScheduledRealtimeClient{client, make(map[string][][]byte), logger}
}
