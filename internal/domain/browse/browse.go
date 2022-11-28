package browse

import (
	"bazaar/config"
	"bazaar/internal/domain/db"
	"bazaar/internal/domain/db/cb"
	"bazaar/internal/domain/filetype"
	"bazaar/internal/domain/storage"
	"bazaar/pkg/crypto"
	"bazaar/pkg/yara"
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
	yara    *yara.Client
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

	yara, err := yara.NewClient(c.Yara.Address)
	if err != nil {
		return nil, err
	}

	return &Browse{
		c:       c,
		db:      db,
		storage: st,
		yara:    yara,
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

	tags, err := b.yara.ScanTags(data)
	if err == nil {
		m.Tags = append(m.Tags, tags...)
	}

	m.Hash = crypto.HashBytes(data)
	m.Path, err = b.storage.Create(data, m.Hash.MD5, m.Type)
	if err != nil {
		return nil, err
	}

	return m, b.db.Create(m.Hash.MD5, m)
}

type QueryMeta struct {
	Hash string   `form:"hash"`
	Type string   `form:"type"`
	Tags []string `form:"tags"`
}

func (b *Browse) MalwareQuery(q *QueryMeta) ([]interface{}, error) {
	var statement string

	if len(q.Hash) > 0 {
		statement = fmt.Sprintf("SELECT data,name,type,md5,sha256,tags FROM %s WHERE md5='%s'",
			b.c.Counchbase.BucketName, q.Hash,
		)
	} else if len(q.Type) > 0 {
		statement = fmt.Sprintf("SELECT data,name,type,md5,sha256,tags FROM %s WHERE type='%s' limit 1000",
			b.c.Counchbase.BucketName, q.Type,
		)
	}

	res, err := b.db.Query(statement, nil)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	result := make([]interface{}, 0)
	meta := make(map[string]interface{})
	for res.Next() {
		if err := res.Row(&meta); err != nil {
			return nil, err
		}

		result = append(result, meta)
	}

	return result, nil
}
