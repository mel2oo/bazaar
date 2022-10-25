package db

// Config represents the database config.
type Config struct {
	// the data source name (DSN) for connecting to the database.
	Server string `mapstructure:"server"`
	// Username used to access the db.
	Username string `mapstructure:"username"`
	// Password used to access the db.
	Password string `mapstructure:"password"`
	// Name of the couchbase bucket.
	BucketName string `mapstructure:"bucket-name"`
}

type Client interface {
	Create(key string, value interface{}) error
	Get(key string, value interface{}) error
}
