package cmd

import (
	"github.com/finove/golibused/pkg/logger"
	"github.com/finove/golibused/pkg/vconfig"
	"github.com/spf13/cobra"
)

var (
	tryWhat        string
	tryCfgFileName string
)

var tryWhatList = []string{"log", "cfg"}

var tryCmd = &cobra.Command{
	Use:   "try",
	Short: "try functions",
	Long:  `try test functions`,
	Run: func(cmd *cobra.Command, args []string) {
		switch tryWhat {
		case "cfg":
			testCfg()
		default:
			testLog()
		}
	},
}

func init() {
	rootCmd.AddCommand(tryCmd)
	tryCmd.Flags().StringVar(&tryWhat, "what", "", "what you want test")
	tryCmd.Flags().StringVar(&tryCfgFileName, "cfg", "", "config file to load for test")
	tryCmd.MarkFlagRequired("what")
}

func testCfg() {
	logger.Info("load config1 %v", vconfig.LoadConfigFile(tryCfgFileName, true))
	vconfig.Viper.SetDefault("cfg1", "value2")
	logger.Info("get cfg1 = %s", vconfig.Viper.GetString("cfg1"))
	vconfig.Viper.WriteConfig()
	logger.Info("work directory:%s, app directory:%s", vconfig.WorkDirectory(), vconfig.AppDirectory())
}

func testLog() {
	logger.Setup(true, "debug", "t.log", `{"tosyslog":true, "appname":"goutils"}`)
	logger.SetLevel("info")
	logger.Fatal("this is fatal log %d", 1)
	logger.Error("this is error log %d", 2)
	logger.Warning("this is warning log %d", 3)
	logger.Info("this is info log %d", 4)
	logger.Debug("this is debug log %d", 5)
	logger.Info("support trywhat %v", tryWhatList)
	logger.NewLogFile("tf.log", 3, "toFile")
	logger.ErrorFor("toFile", "to file log info %d", 11)
	logger.WarningFor("toFile", "to file log info %d", 12)
	logger.InfoFor("toFile", "to file log info %d", 13)
}
