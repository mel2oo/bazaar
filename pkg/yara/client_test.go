package yara

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"bazaar/pkg/yara/proto"
)

func TestFileScan(t *testing.T) {
	c, err := NewClient("127.0.0.1:6141")
	if err != nil {
		t.Fail()
		fmt.Println(err)
		return
	}

	err = filepath.Walk("/Users/switch/Desktop/yarascan", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		fmt.Println("scan", path)
		res, err := c.ScanFile(context.Background(),
			&proto.Request{
				Path: path,
			})
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(res)

		return nil
	})

	if err != nil {
		t.Log(err)
	}

}

func TestFileScanBlob(t *testing.T) {
	c, err := NewClient("127.0.0.1:6141")
	if err != nil {
		t.Fail()
		fmt.Println(err)
		return
	}

	file := "/home/whales/samples/53/filedump/1D85B2BB3E2FC95_EE43E519.xls"
	f, err := os.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	stream, err := c.ScanBlob(context.TODO())
	if err != nil {
		stream.CloseSend()
		t.Fatal(err)
	}

	s := make([]byte, 1024*1024*4-512)
	for {
		nr, err := f.Read(s[:])
		if err == io.EOF || nr == 0 {
			break
		}
		stream.Send(&proto.RequestBlob{Blob: s[0:nr]})
	}

	t.Log(stream.CloseAndRecv())
}
