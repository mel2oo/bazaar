package cb

import (
	"bazaar/internal/domain/db"
	"errors"
	"time"

	"github.com/couchbase/gocb/v2"
)

const (
	// Duration to wait until memd connections have been established with
	// the server and are ready.
	timeout = 30 * time.Second
)

type Client struct {
	bucket     *gocb.Bucket
	cluster    *gocb.Cluster
	collection *gocb.Collection
}

func New(c db.Config) (*Client, error) {
	options := gocb.ClusterOptions{
		Authenticator: gocb.PasswordAuthenticator{
			Username: c.Username,
			Password: c.Password,
		},
	}

	// Initialize the Connection
	cluster, err := gocb.Connect(c.Server, options)
	if err != nil {
		return nil, err
	}

	// Get a bucket reference.
	bucket := cluster.Bucket(c.BucketName)

	// We wait until the bucket is definitely connected and setup.
	err = bucket.WaitUntilReady(timeout, nil)
	if err != nil {
		return nil, err
	}

	// Get a collection reference.
	collection := bucket.DefaultCollection()

	return &Client{
		bucket:     bucket,
		cluster:    cluster,
		collection: collection,
	}, nil
}

// Create saves a new document into the collection.
func (c *Client) Create(key string, value interface{}) error {

	// Insert document with timeout 3s
	_, err := c.collection.Insert(key, value, &gocb.InsertOptions{
		Timeout: time.Second * 3,
	})

	// If document exists, return success
	if errors.Is(err, gocb.ErrDocumentExists) {
		return nil
	}

	return err
}

// Get retrieves the document using its key.
func (c *Client) Get(key string, value interface{}) error {

	// Performs a fetch operation against the collection.
	res, err := c.collection.Get(key, &gocb.GetOptions{})
	if err != nil {
		return err
	}

	// Assigns the value of the result into the valuePtr using default decoding.
	return res.Content(&value)
}
