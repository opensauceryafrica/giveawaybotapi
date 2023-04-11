package typing

type Env struct {
	//Port for the server to listen on
	Port string `env:"PORT"`
	//MongoDBURI is the monogoDB connection string for the app
	MongoDBURI string `env:"MONGO_DB_URI"`
	//MongoDBName is the MongoDB database name for the app
	MongoDBName string `env:"MONGO_DB_NAME"`
	//PostgreSQLURI is the PostgreSQL database uri string for App
	// PostgreSQLURI string `env:"POSTGRE_SQL_URI"`
}
