package database

import (
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	PostgreSQLDB *pgxpool.Pool
	MongoDB      *mongo.Database
)
