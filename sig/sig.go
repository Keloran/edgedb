package sig

import (
	"context"
	"fmt"
	"github.com/bugfixes/go-bugfixes/logs"
	"os"
	"time"
)

type System struct {
	context context.Context
	folder  string
}

func NewSystem() *System {
	return &System{
		context: context.Background(),
	}
}

func (s *System) SetContext(ctx context.Context) {
	s.context = ctx
}

func (s *System) SetFolder(folder string) {
	s.folder = folder

	// make the folder if it doesn't exist
	if _, err := os.Stat(s.folder); os.IsNotExist(err) {
		if err := os.Mkdir(s.folder, 0755); err != nil {
			logs.Fatalf("Failed to create folder: %v", err)
		}
	}
}

func (s *System) SendLogs(count uint32) error {
	if s.folder == "" {
		logs.Infof("consul_open_http_connections %d\n", count)
		return nil
	}

	// create a temp file so that we can write data to it and the exporter not block
	tmp := fmt.Sprintf("%s/temp.prom", os.TempDir())
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	_, err = f.WriteString(fmt.Sprintf("consul_open_http_connections %d\n", count))
	if err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	// rename the template file to the actual file
	if err := os.Rename(f.Name(), fmt.Sprintf("%s/%d.prom", s.folder, time.Now().Unix())); err != nil {
		return err
	}

	return nil
}
