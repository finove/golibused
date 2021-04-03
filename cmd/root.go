package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/finove/golibused/pkg/logger"
	"github.com/finove/golibused/pkg/vconfig"
	"github.com/spf13/cobra"
)

const golibusedVersion = "0.0.1"

var (
	configFile string
	logFile    string
)

var rootCmd = &cobra.Command{
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		vconfig.Viper.Set(fmt.Sprintf("%s_version", cmd.Root().Use), cmd.Root().Version)
		// vconfig.Viper.WriteConfig()
	},
	Use:     "golibused",
	Version: golibusedVersion,
}

// Execute 执行命令行主程序
func Execute() {
	var err error
	if rootCmd.HasSubCommands() {
		var appName = filepath.Base(os.Args[0])
		for _, cmd := range rootCmd.Commands() {
			if cmd.Use == appName {
				var newArgs = make([]string, len(os.Args))
				copy(newArgs, os.Args)
				newArgs[0] = appName
				rootCmd.SetArgs(newArgs)
			}
		}
	}
	if err = rootCmd.Execute(); err != nil {
		logger.Fatal("execute fail:%v", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(loadConfig)
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "set config file")
	rootCmd.PersistentFlags().StringVar(&logFile, "logpath", "", "set log file")
}

func loadConfig() {
	vconfig.LoadConfigFile(configFile, true)
}
