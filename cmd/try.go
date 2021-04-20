package cmd

import (
	"fmt"
	"time"

	"github.com/finove/golibused/pkg/errormsg"
	"github.com/finove/golibused/pkg/logger"
	"github.com/finove/golibused/pkg/timedo"
	"github.com/finove/golibused/pkg/vconfig"
	"github.com/spf13/cobra"
)

var (
	tryWhat        string
	tryCfgFileName string
)

var tryWhatList = []string{"log", "cfg", "errmsg"}

var tryCmd = &cobra.Command{
	Use:   "try",
	Short: "try functions",
	Long:  `try test functions`,
	Run: func(cmd *cobra.Command, args []string) {
		switch tryWhat {
		case "cfg":
			testCfg()
		case "errmsg":
			testErrorMsg()
		case "timedo":
			testTimeDO()
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

func testErrorMsg() {
	errormsg.AddErrorMessages(map[int][]string{
		2001: {"403", "fail3"},
		2002: {"404", "fail4"},
		2003: {"494", "fail5"},
	})
	errormsg.AddErrorMessage(1002, "401", "未登录", "认证失败或没有授权")
	logger.Info("code 0 = %d %s", errormsg.HTTPStatus(0), errormsg.Message(0))
	logger.Info("code 1001 = %d %s", errormsg.HTTPStatus(1001), errormsg.Message(1001, fmt.Errorf("more %s", "info")))
	logger.Info("code 1002 = %d %s", errormsg.HTTPStatus(1002), errormsg.Message(1002, fmt.Errorf("补充失败信息 %d", 1002)))
	logger.Info("code 2001 = %d %s", errormsg.HTTPStatus(2001), errormsg.Message(2001))
	logger.Info("code 2002 = %d %s", errormsg.HTTPStatus(2002), errormsg.Message(2002))
	logger.Info("code 2003 = %d %s", errormsg.HTTPStatus(2003), errormsg.Message(2003))
}

func testTimeDO() {
	var s1 = time.Now()
	var s2 = s1.Add(time.Hour)
	doCheck(s1, s2, 30*time.Minute, 30*time.Minute)
	doCheck(s1, s2, 60*time.Minute, 30*time.Minute)
	doCheck(s1, s2, 75*time.Minute, 30*time.Minute)
	doCheck(s1, s2, 90*time.Minute, 30*time.Minute)
	doCheck(s2, s1, 30*time.Minute, 30*time.Minute)
	doCheck(s2, s1, 30*time.Minute, 60*time.Minute)
	doCheck(s2, s1, 30*time.Minute, 75*time.Minute)
	doCheck(s2, s1, 30*time.Minute, 90*time.Minute)
	f, t, err := timedo.DateDuration("2021-04-11", "2021-04-18")
	logger.Info("timedo dateduration %v - %v, err %v", f, t, err)
	f, t = timedo.Week(time.Now(), true)
	logger.Info("timedo week %v - %v", f, t)
}

func doCheck(s1, s2 time.Time, d1, d2 time.Duration) {
	logger.Info("timedo %s-%s %s-%s %v", s1.Format("15:04:05"), s1.Add(d1).Format("15:04:05"), s2.Format("15:04:05"), s2.Add(d2).Format("15:04:05"), timedo.Intersection(s1, d1, s2, d2))
}
