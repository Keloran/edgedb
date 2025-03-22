package main

import (
	"context"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/keloran/edgedb/conex"
	"github.com/keloran/edgedb/sig"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"time"
)

func main() {
	var (
		port     = uint32(8500)
		testMode = false
	)

	rootCmd := &cobra.Command{
		Use: "edgedb",
		Run: func(cmd *cobra.Command, args []string) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			logs.Infof("Listening on port: %d", port)

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
				if !testMode {
					sigs.SetFolder("/tmp/node-exporter")
				}
				if err := sigs.SendLogs(conns); err != nil {
					logs.Infof("Error sending logs: %v", err)
					cancel()
				}
			}()
			<-signalChan

			cancel()
			time.Sleep(time.Second * 5)
			logs.Info("Shutting down...")
		},
	}

	rootCmd.Flags().Uint32VarP(&port, "port", "p", 8500, "Port to use for connections")
	rootCmd.Flags().BoolVarP(&testMode, "test", "t", false, "Run in test mode (don't set folder)")

	if err := rootCmd.Execute(); err != nil {
		logs.Fatal(err)
		os.Exit(1)
	}
}
