package cmd

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	delivery "github.com/danluki/test-task-8/internal/delivery/http"
	"github.com/danluki/test-task-8/internal/server"
	"github.com/danluki/test-task-8/internal/store"
	"github.com/spf13/cobra"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	appCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start TestTask service",
	Long:  "Start TestTask service",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		db, err := gorm.Open(postgres.Open(cfg.Database.Url), &gorm.Config{})
		if err != nil {
			return err
		}
		err = db.AutoMigrate(&store.User{})
		if err != nil {
			return err
		}

		handlers := delivery.NewHandler(db)

		srv := server.NewServer(cfg, handlers.Init(cfg))

		go func() {
			if err := srv.Run(); !errors.Is(err, http.ErrServerClosed) {
				slog.Error(err.Error())
			}
		}()

		slog.Info("Server started")

		// Graceful Shutdown
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

		<-quit

		const timeout = 5 * time.Second

		ctx, shutdown := context.WithTimeout(ctx, timeout)
		defer shutdown()

		if err := srv.Stop(ctx); err != nil {
			slog.Error(err.Error())
		}

		return nil
	},
}
