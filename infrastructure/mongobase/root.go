package mongobase

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	INIT_SESSION_TIMEOUT      = 10 * time.Second
	TERMINATE_SESSION_TIMEOUT = 2 * time.Second
)

func (m *MongoServiceBase) DB() *mongo.Database {
	return m.Session.Database(m.Database)
}

func (m *MongoServiceBase) InitSession(databaseURI, databaseName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), INIT_SESSION_TIMEOUT)
	defer cancel()

	if mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(databaseURI)); err != nil {
		return err
	} else {
		if err = mongoClient.Ping(ctx, nil); err != nil {
			return err
		}

		m.Session = mongoClient
		m.Database = databaseName
	}

	return nil
}

func (m *MongoServiceBase) TerminateSession() {
	if m.Session != nil {
		ctx, cancel := context.WithTimeout(context.Background(), TERMINATE_SESSION_TIMEOUT)
		defer cancel()
		m.Session.Disconnect(ctx)
	}
}
