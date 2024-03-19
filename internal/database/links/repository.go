package links

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/ptsypyshev/gb-golang-level3-new/internal/database"
)

const collection = "links"

func New(db *mongo.Database, timeout time.Duration) *Repository {
	return &Repository{db: db, timeout: timeout}
}

type Repository struct {
	db      *mongo.Database
	timeout time.Duration
}

func (r *Repository) Create(ctx context.Context, req CreateReq) (database.Link, error) {
	var l database.Link
	l.ID = req.ID
	l.URL = req.URL
	l.Title = req.Title
	l.Tags = req.Tags
	l.Images = req.Images
	l.UserID = req.UserID
	l.CreatedAt = time.Now()
	l.UpdatedAt = l.CreatedAt

	_, err := r.db.Collection(collection).InsertOne(ctx, l)
	if err != nil {
		return database.Link{}, err
	}

	return l, nil
}

func (r *Repository) FindByUserAndURL(ctx context.Context, link, userID string) (database.Link, error) {
	var l database.Link
	
	filter := bson.M{"url": bson.M{"$eq": link}, "userID": bson.M{"$eq": userID}}
	cursor, err := r.db.Collection(collection).Find(ctx, filter)
	if err != nil {
		return l, err
	}
	var filteredLinks []database.Link
	if err = cursor.All(ctx, &filteredLinks); err != nil || len(filteredLinks) == 0 {
		return l, err
	}
	return filteredLinks[0], nil
}

func (r *Repository) FindByCriteria(ctx context.Context, criteria Criteria) ([]database.Link, error) {
	return nil, nil
}
