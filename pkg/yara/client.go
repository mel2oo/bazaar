package yara

import (
	"bazaar/pkg/yara/proto"
	"context"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type Client struct {
	proto.YaraScannerClient
}

func NewClient(address string) (*Client, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &Client{
		YaraScannerClient: proto.NewYaraScannerClient(conn),
	}, nil
}

func (c *Client) ScanBlobs(ctx context.Context, data []byte) (*proto.Result, error) {
	stream, err := c.ScanBlob(ctx)
	if err != nil {
		return nil, err
	}

	limit := 1024*1024*3 - 512
	times := len(data) / limit
	if len(data)%limit > 0 {
		times += 1
	}

	for i := 0; i < times; i++ {
		if i == times-1 {
			if err = stream.Send(&proto.RequestBlob{Blob: data[i*limit:]}); err != nil {
				break
			}
			continue
		}
		if err := stream.Send(&proto.RequestBlob{Blob: data[i*limit : (i+1)*limit]}); err != nil {
			break
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, err
	}

	return res.Result, nil
}

func (c *Client) ScanTags(data []byte) ([]string, error) {
	mtags := make(MapTag)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	res, err := c.ScanBlobs(ctx, data)
	if err != nil {
		logrus.Error("yara scan error", err)
		return nil, err
	}

	if res.Hit {
		for _, m := range res.Matches {
			mtags.ParseTags(m.Tags)
		}
	}

	return mtags.Format(), nil
}

type MapTag map[string]bool

func (m MapTag) ParseTags(s string) {
	tlist := strings.Split(s, "|")

	for _, tstr := range tlist {
		list := strings.Split(tstr, "_")

		if len(list) == 1 {
			m[list[0]] = true
		}

		if len(list) == 2 {
			m[list[0]] = true
			m[list[1]] = true
		}

		if len(list) == 3 {
			m[list[1]] = true
			m[list[2]] = true
		}
	}
}

func (m MapTag) Format() []string {
	tags := make([]string, 0)
	for t := range m {
		tags = append(tags, t)
	}
	return tags
}
