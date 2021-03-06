package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/finove/golibused/net/fsp"
	"github.com/finove/golibused/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	serverIP                       string
	localPort, remotePort          uint
	serverPass, serverNewPass      string
	cmdLS, cmdGet, cmdSave, cmdPut string
	showServerVersion              bool
)

var fspCmd = &cobra.Command{
	Use:     "fsp",
	Version: "0.0.1",
	Short:   "download file using fsp protocol",
	Example: "some example usage",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var fspSession *fsp.Session
		var addr *net.UDPAddr
		var conn *net.UDPConn
		addr, conn, err = getFSPServerIP()
		if err != nil {
			logger.Info("Failed, get fsp server ip %v", err)
			return
		}
		fspSession, err = fsp.NewSessionWithConn(conn, addr.String(), serverPass)
		if err != nil {
			logger.Info("Failed, open fsp session %v", err)
			return
		}
		if showServerVersion {
			fmt.Printf("fsp server version: %s\n", fspSession.Version())
		} else if cmdLS != "" {
			err = fspSession.ShowDir(cmdLS)
			if err != nil {
				logger.Info("Failed, read dir %v", err)
			}
		} else if cmdGet != "" {
			if len(cmdGet) > 0 && os.IsPathSeparator(cmdGet[len(cmdGet)-1]) {
				err = fspSession.DownloadDirectory(cmdGet, cmdSave)
			} else {
				err = fspSession.DwonloadFile(cmdGet, cmdSave, 3)
			}
			if err != nil {
				logger.Info("Failed, get file %s error %v", cmdGet, err)
			}
		} else if serverNewPass != "" {
			err = fspSession.ChangePassword(serverNewPass)
			if err != nil {
				logger.Info("Failed, change password error %v", err)
			}
		} else if cmdPut != "" {
			err = fspSession.UploadFile(cmdPut, cmdSave)
			if err != nil {
				logger.Info("Failed, upload file error %v", err)
			}
		}
		fspSession.Close()
	},
}

func init() {
	rootCmd.AddCommand(fspCmd)
	fspCmd.Flags().StringVar(&serverIP, "ip", "", "fsp server ip:port")
	fspCmd.Flags().UintVar(&localPort, "port", 0, "local port for used")
	fspCmd.Flags().UintVar(&remotePort, "dport", 0, "fsp server port")
	fspCmd.Flags().StringVarP(&serverPass, "password", "p", "", "fsp server password")
	fspCmd.Flags().StringVar(&serverNewPass, "np", "", "change the password of FSP server")
	fspCmd.Flags().StringVar(&cmdPut, "put", "", "upload file path")
	fspCmd.Flags().StringVar(&cmdLS, "ls", "", "fsp command list files")
	fspCmd.Flags().StringVarP(&cmdGet, "get", "g", "", "fsp command get files")
	fspCmd.Flags().StringVarP(&cmdSave, "save", "s", "", "get file save path")
	fspCmd.Flags().BoolVar(&showServerVersion, "server_version", false, "show server version")
}

func getFSPServerIP() (addr *net.UDPAddr, conn *net.UDPConn, err error) {
	var localAddr *net.UDPAddr
	if localPort > 0 {
		localAddr, _ = net.ResolveUDPAddr("udp4", fmt.Sprintf(":%d", localPort))
	}
	if serverIP != "" {
		addr, err = net.ResolveUDPAddr("udp4", serverIP)
		if err != nil {
			return nil, nil, err
		}
		conn, err = net.ListenUDP("udp4", localAddr)
	} else {
		err = fmt.Errorf("miss command line parameter, need -ip or -id or -mac value")
		return nil, nil, err
	}
	return
}
