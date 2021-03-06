// Generated with the code-generator
//
// Modifications in code regions will be lost during regeneration!

package main

import (
	"context"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gookit/slog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "CleanBackup",
	Short: "R6S backup files clean tool",
	Long:  "Clean R6S backup files periodically.",
	Args:  cobra.NoArgs,
}

func startCleanUp(ctx context.Context) error {
	slog.Info("startCleanUp")

	path := viper.Get("Path")
	in := viper.Get("Ext")

	switch ext := in.(type) {
	case []interface{}:
		for _, value := range ext {
			pattern := path.(string) + "*." + value.(string)
			files, err := filepath.Glob(pattern)
			if err != nil {
				slog.Error(err)
			}
			for _, f := range files {
				slog.Info("Clean up:", f)
				if err := os.Remove(f); err != nil {
					slog.Error(err)
				}
			}
		}
	}

	return nil
}

var cleanCmd = &cobra.Command{
	Use:   "runOnce",
	Short: "Clean R6S backup files once",
	Long:  "Clean R6S backup files once.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		return startCleanUp(ctx)
	},
}

func deamonCleanUp(ctx context.Context) error {
	slog.Info("deamonCleanUp")

	period := viper.Get("CheckPeriod").(float64)
	if period < 1 {
		slog.Error("'CheckPeriod' is too small, ", period)
		period = 1
	}

	tick := time.Tick(time.Duration(period) * time.Minute)
	for {
		select {
		case <-ctx.Done():
			slog.Info(ctx.Err())
			return nil
		case <-tick:
			startCleanUp(ctx)
		default:
			time.Sleep(time.Second)
		}
	}
}

var deamonCmd = &cobra.Command{
	Use:   "deamon",
	Short: "Clean R6S backup files periodically",
	Long:  "Run in background, clean R6S backup files periodically.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		// defer cancel()

		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		go func() {
			for sig := range c {
				switch sig {
				case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
					slog.Info("Exit with:", sig)
					cancel()
					return
				}
			}
		}()

		return deamonCleanUp(ctx)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(deamonCmd)

	//cleanCmd.PersistentFlags().StringVar(&env, "env", env, "The env of the importer will run.")
	InitViperConfig(rootCmd)
}

func main() {

	defer func() {
		viper.WriteConfig()
		viper.WriteConfigAs("config_templ.json")
	}()

	slog.Configure(func(logger *slog.SugaredLogger) {
		f := logger.Formatter.(*slog.TextFormatter)
		f.EnableColor = true
	})

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

}
