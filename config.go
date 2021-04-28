package main

import (
	"github.com/gookit/slog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func defaults() {
	viper.SetDefault("Path", "Z:/PRIVATE/Backup/") // path
	viper.SetDefault("Ext", []string{"zip", "bak"})
	viper.SetDefault("CheckPeriod", 10) // in minutes
}

func logViperSettings() {
	slog.Info("Path:", viper.Get("Path"))
	slog.Info("Extensions:", viper.Get("Ext"))
	slog.Info("CheckPeriod(in minute):", viper.Get("CheckPeriod"))
}

func InitViperConfig(rootCmd *cobra.Command) {

	viper.SetConfigName("config.json") // name of config file (without extension)
	viper.SetConfigType("json")        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath(".")           // optionally look for config in the working directory

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			slog.Error("Read config config file error: %s \n", err)
			defaults()
		} else {
			// Config file was found but another error was produced
			slog.Error("Read config config file error: %s \n", err)
			defaults()
		}
	}

	// rootCmd.PreRun = func(cmd *cobra.Command, args []string) {
	// 	if viper.GetBool("verbose") {
	logViperSettings()
	// }
	// }
}
