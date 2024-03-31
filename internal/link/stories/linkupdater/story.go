package linkupdater

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ptsypyshev/gb-golang-level3-new/internal/database"
	"github.com/ptsypyshev/gb-golang-level3-new/internal/link/models"

	"github.com/ptsypyshev/gb-golang-level3-new/pkg/scrape"
)

func New(repository repository, consumer amqpConsumer, queueName string) *Story {
	return &Story{
		repository: repository,
		consumer:   consumer,
		queueName:  queueName,
	}
}

type Story struct {
	repository repository
	consumer   amqpConsumer
	queueName  string
}

func (s *Story) Run(ctx context.Context) error {
	// Слушаем очередь
	ch, err := s.consumer.Consume(s.queueName, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case m, ok := <-ch:
			if !ok {
				return errors.New("rabbitmq queue is closed")
			}
			err := s.processMsg(ctx, m)
			if err != nil {
				slog.Error("process message error", slog.Any("err", err))
			}
		}
	}
}

func (s *Story) processMsg(ctx context.Context, msg amqp091.Delivery) error {
	var m models.Message
	err := json.Unmarshal(msg.Body, &m)
	if err != nil {
		return err
	}

	id, err := primitive.ObjectIDFromHex(m.ID)
	if err != nil {
		return err
	}

	// Получаем текущий объект ссылки
	link, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Добавляем данные из scrape
	parsed, err := scrape.Parse(ctx, link.URL)
	if err != nil {
		return err
	}

	if parsed.Title != "" {
		link.Title = parsed.Title
	}

	if len(parsed.Tags) > 0 {
		link.Tags = append(link.Tags, parsed.Tags...)
	}

	req := database.UpdateLinkReq{
		ID:     id,
		Title:  link.Title,
		URL:    link.URL,
		Images: link.Images,
		Tags:   link.Tags,
		UserID: link.UserID,
	}

	// Обновляем данные в DB
	_, err = s.repository.Update(ctx, req)
	return err
}
