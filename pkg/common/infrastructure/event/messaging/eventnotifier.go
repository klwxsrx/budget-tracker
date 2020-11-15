package messaging

import (
	"errors"
	"fmt"
	app "github.com/klwxsrx/expense-tracker/pkg/common/app/event"
	"github.com/klwxsrx/expense-tracker/pkg/common/infrastructure/mysql"
	"github.com/streadway/amqp"
)

const (
	exchangeName       = "domain_event"
	exchangeType       = "topic"
	routingKey         = "last_event_id"
	messageContentType = "text/plain"
)

type StoredEventNotifier struct {
	store   app.Store
	client  mysql.Client
	channel *amqp.Channel
}

func (en *StoredEventNotifier) NotifyOfCreatedEvents() error {
	if en.channel == nil {
		return errors.New("channel is not connected")
	}

	lock := mysql.NewLock(en.client, "stored_event_notification")
	err := lock.Get()
	if err != nil {
		return err
	}
	defer lock.Release()

	lastSentID, err := en.getLastSentEventID()
	if err != nil {
		return fmt.Errorf("can't get last sent id: %v", err)
	}
	lastID, err := en.store.LastID()
	if err != nil {
		return fmt.Errorf("can't get last id from store: %v", err)
	}
	if lastID == lastSentID {
		return nil
	}
	if err := en.sendNotification(lastID); err != nil {
		return fmt.Errorf("can't send notification: %v", err)
	}
	return en.updateLastSentEventID(lastSentID)
}

func (en *StoredEventNotifier) Connect(connection *amqp.Connection) error {
	var err error
	en.channel, err = connection.Channel()
	if err != nil {
		return err
	}

	err = en.channel.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil) // TODO: test
	if err != nil {
		return err
	}
	return nil
}

func (en *StoredEventNotifier) getLastSentEventID() (app.StoredEventID, error) {
	var id app.StoredEventID
	err := en.client.Get(&id, "SELECT IFNULL(MAX(id), 0) FROM last_notified_event")
	return id, err
}

func (en *StoredEventNotifier) updateLastSentEventID(id app.StoredEventID) error {
	_, err := en.client.Exec("INSERT INTO last_notified_event (id) VALUES (?) ON DUPLICATE KEY IGNORE", id)
	return err
}

func (en *StoredEventNotifier) sendNotification(id app.StoredEventID) error {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  messageContentType,
		Body:         []byte(fmt.Sprintf("%v", id)),
	}
	return en.channel.Publish(exchangeName, routingKey, false, false, msg)
}

func NewStoredEventNotifier(store app.Store, client mysql.Client) *StoredEventNotifier {
	return &StoredEventNotifier{store, client, nil}
}
