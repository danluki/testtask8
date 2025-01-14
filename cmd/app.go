package cmd

import (
	"fmt"
	"os"

	"github.com/danluki/test-task-8/internal/config"
	"github.com/spf13/cobra"
)

var cfg = config.DefaultConfig()

var cfgFile string

var appCmd = &cobra.Command{
	Use:   "testtask",
	Short: "TestTask service",
	Long:  "TestTask service",
}

func Execute() {
	if err := appCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	appCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
}

func initConfig() {
	var err error
	cfg, err = config.Load(cfgFile)
	if err != nil {
		panic(err)
	}
}
