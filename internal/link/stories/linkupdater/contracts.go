package linkupdater

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ptsypyshev/gb-golang-level3-new/internal/database"
)

type repository interface {
	FindByID(ctx context.Context, id primitive.ObjectID) (database.Link, error)
	Update(ctx context.Context, req database.UpdateLinkReq) (database.Link, error)
}

type amqpConsumer interface {
	Consume(queue, consumer string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (
		<-chan amqp.Delivery,
		error,
	)
}
