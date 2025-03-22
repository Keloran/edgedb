package conex

import (
	"context"
	"github.com/shirou/gopsutil/net"
	"time"
)

type System struct {
	context context.Context

	Count    uint32
	Port     uint32
	Interval time.Duration
}

func NewSystem() *System {
	return &System{
		context: context.Background(),
		Count:   0,
		Port:    8500,
	}
}

func (s *System) SetContext(ctx context.Context) {
	s.context = ctx
}

func (s *System) CountConnections() (uint32, error) {
	// get all tcp connections that way we can listen for ipv4 and ipv6
	conn, err := net.Connections("tcp")
	if err != nil {
		return 0, err
	}

	count := uint32(0)
	for _, conn := range conn {
		if conn.Raddr.Port == s.Port {
			// we only care about established connections
			if conn.Status == "ESTABLISHED" {
				count++
			}
		}
	}

	return count, nil
}
