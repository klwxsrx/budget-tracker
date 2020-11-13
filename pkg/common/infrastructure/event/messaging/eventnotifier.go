package messaging

import (
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

type StoredEventNotifier interface {
	NotifyOfCreatedEvents() error
}

type storedEventNotifier struct {
	store   app.Store
	client  mysql.Client
	channel *amqp.Channel
}

func (en *storedEventNotifier) NotifyOfCreatedEvents() error {
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

func (en *storedEventNotifier) Connect(connection *amqp.Connection) error {
	ch, err := connection.Channel()
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(exchangeName, exchangeType, true, false, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (en *storedEventNotifier) getLastSentEventID() (app.StoredEventID, error) {
	var id app.StoredEventID
	err := en.client.Get(&id, "SELECT id FROM last_notified_event LIMIT 1")
	return id, err
}

func (en *storedEventNotifier) updateLastSentEventID(id app.StoredEventID) error {
	_, err := en.client.Exec("UPDATE last_notified_event SET id = ?", id)
	return err
}

func (en *storedEventNotifier) sendNotification(id app.StoredEventID) error {
	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  messageContentType,
		Body:         []byte(fmt.Sprintf("%v", id)),
	}
	return en.channel.Publish(exchangeName, routingKey, false, false, msg)
}

func NewStoredEventNotifier(store app.Store, client mysql.Client, channel *amqp.Channel) StoredEventNotifier {
	return &storedEventNotifier{store, client, channel}
}
