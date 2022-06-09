package mongobase

import "go.mongodb.org/mongo-driver/mongo"

type MongoServiceBase struct {
	Session  *mongo.Client
	Database string
}
