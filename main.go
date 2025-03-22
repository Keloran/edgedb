package main

import (
	"context"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/keloran/edgedb/conex"
	"github.com/keloran/edgedb/sig"
	"os"
	"os/signal"
	"time"
)

func main() {
	port := uint32(8500)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logs.Info("Starting...")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		con := conex.NewSystem()
		con.SetContext(ctx)
		con.Interval = time.Second * 5
		con.Port = port
		conns, err := con.CountConnections()
		if err != nil {
			logs.Infof("Error counting the connections: %v", err)
			cancel()
		}

		sigs := sig.NewSystem()
		sigs.SetContext(ctx)
		//sigs.SetFolder("/tmp/node-exporter")
		if err := sigs.SendLogs(conns); err != nil {
			logs.Infof("Error sending logs: %v", err)
			cancel()
		}
	}()
	<-signalChan

	cancel()
	time.Sleep(time.Second * 5)
	logs.Info("Shutting down...")
}
