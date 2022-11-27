package browse

import (
	"bazaar/config"
	"bazaar/internal/domain/db"
	"bazaar/internal/domain/db/cb"
	"bazaar/internal/domain/filetype"
	"bazaar/internal/domain/storage"
	"bazaar/pkg/crypto"
	"bazaar/pkg/yara"
	"io/ioutil"
	"mime/multipart"
	"time"

	"github.com/sirupsen/logrus"
)

type Browse struct {
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
		db:      db,
		storage: st,
		yara:    yara,
	}, nil
}

type Malware struct {
	File *multipart.FileHeader `json:"-" form:"file" binding:"required"`

	Date string        `json:"date,omitempty" form:"date"`
	Name string        `json:"name,omitempty" form:"name"`
	Path string        `json:"path,omitempty"`
	Type string        `json:"type,omitempty" form:"type"`
	Tags []string      `json:"tags,omitempty" form:"tags"`
	Hash crypto.Result `json:"hash,omitempty"`
}

func (b *Browse) MalwareCreate(m *Malware) error {

	if len(m.Date) == 0 {
		m.Date = time.Now().Local().String()
	}

	if len(m.Name) == 0 {
		m.Name = m.File.Filename
	}

	logrus.Infof("recevie upload file: %s", m.Name)

	fi, err := m.File.Open()
	if err != nil {
		return err
	}
	defer fi.Close()

	data, err := ioutil.ReadAll(fi)
	if err != nil {
		return err
	}

	if len(m.Type) == 0 {
		m.Type, err = filetype.ScanData(data)
		if err != nil {
			logrus.Warnf("yara tags scan error", err)
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
		return err
	}

	return b.db.Create(m.Hash.MD5, m)
}
