package browse

import (
	"bazaar/config"
	"bazaar/internal/domain/db"
	"bazaar/internal/domain/db/cb"
	"bazaar/internal/domain/filetype"
	"bazaar/internal/domain/storage"
	"bazaar/pkg/crypto"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/sirupsen/logrus"
)

type Browse struct {
	c       *config.Config
	db      db.Client
	storage *storage.Client
}

func New(c *config.Config) (*Browse, error) {
	db, err := cb.New(c.Counchbase)
	if err != nil {
		return nil, err
	}

	st, err := storage.New(c.Storage)
	if err != nil {
		return nil, err
	}

	return &Browse{
		c:       c,
		db:      db,
		storage: st,
	}, nil
}

type Malware struct {
	File *multipart.FileHeader `json:"-" form:"file" binding:"required"`

	Date string   `json:"date,omitempty" form:"date"`
	Name string   `json:"name,omitempty" form:"name"`
	Path string   `json:"path,omitempty"`
	Type string   `json:"type,omitempty" form:"type"`
	Tags []string `json:"tags,omitempty" form:"tags"`
	crypto.Hash
}

func (b *Browse) MalwareCreate(m *Malware) (*Malware, error) {

	if len(m.Date) == 0 {
		m.Date = time.Now().Local().String()
	}

	if len(m.Name) == 0 {
		m.Name = m.File.Filename
	}

	logrus.Infof("recevie upload file: %s", m.Name)

	fi, err := m.File.Open()
	if err != nil {
		return nil, err
	}
	defer fi.Close()

	data, err := io.ReadAll(fi)
	if err != nil {
		return nil, err
	}

	if len(m.Type) == 0 {
		m.Type, err = filetype.ScanData(data)
		if err != nil {
			return nil, err
		}
	}

	if len(m.Tags) == 0 {
		m.Tags = make([]string, 0)
	}

	m.Hash = crypto.HashBytes(data)
	m.Path, err = b.storage.Create(data, m.Hash.MD5, m.Type)
	if err != nil {
		return nil, err
	}

	return m, b.db.Create(m.Hash.MD5, m)
}

type QueryMeta struct {
	Hash string `form:"hash"`
	Type string `form:"type"`
	Tags string `form:"tags"`
	Size int    `form:"size"`
}

func (b *Browse) MalwareQuery(q *QueryMeta) ([]interface{}, error) {
	var statement string

	if q.Size == 0 {
		q.Size = 1000
	}

	if len(q.Hash) > 0 {
		statement = fmt.Sprintf("SELECT date,name,type,md5,sha256,tags FROM %s WHERE md5='%s' OR sha256='%s'",
			b.c.Counchbase.BucketName, q.Hash, q.Hash,
		)
	} else if len(q.Type) > 0 && len(q.Tags) == 0 {
		statement = fmt.Sprintf("SELECT date,name,type,md5,sha256,tags FROM %s WHERE type='%s' limit %d",
			b.c.Counchbase.BucketName, q.Type, q.Size,
		)
	} else if len(q.Type) > 0 && len(q.Tags) > 0 {
		statement = fmt.Sprintf("SELECT date,name,type,md5,sha256,tags FROM %s WHERE '%s' IN tags AND type='%s' limit %d",
			b.c.Counchbase.BucketName, q.Tags, q.Type, q.Size)
	} else if len(q.Type) == 0 && len(q.Tags) > 0 {
		statement = fmt.Sprintf("SELECT date,name,type,md5,sha256,tags FROM %s WHERE '%s' IN tags limit %d",
			b.c.Counchbase.BucketName, q.Tags, q.Size)
	} else {
		statement = fmt.Sprintf("SELECT date,name,type,md5,sha256,tags FROM %s limit %d",
			b.c.Counchbase.BucketName, q.Size)
	}

	res, err := b.db.Query(statement, nil)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	result := make([]interface{}, 0)
	for res.Next() {
		meta := make(map[string]interface{})

		if err := res.Row(&meta); err != nil {
			return nil, err
		}

		result = append(result, meta)
	}

	return result, nil
}

func (b *Browse) MalwareCount(q *QueryMeta) (map[string]interface{}, error) {
	var statement string

	if len(q.Type) > 0 && len(q.Tags) == 0 {
		statement = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE type='%s'",
			b.c.Counchbase.BucketName, q.Type)
	} else if len(q.Type) == 0 && len(q.Tags) > 0 {
		statement = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE '%s' IN tags",
			b.c.Counchbase.BucketName, q.Tags)
	} else if len(q.Type) > 0 && len(q.Tags) > 0 {
		statement = fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE '%s' IN tags AND type='%s'",
			b.c.Counchbase.BucketName, q.Tags, q.Type)
	} else {
		statement = fmt.Sprintf("SELECT COUNT(*) FROM %s", b.c.Counchbase.BucketName)
	}

	res, err := b.db.Query(statement, nil)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var result map[string]interface{}
	if err := res.One(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func (b *Browse) MalwareDownload(md5 string) (*Malware, error) {
	var m Malware

	if err := b.db.Get(md5, &m); err != nil {
		return nil, err
	}

	return &m, nil
}
