package command

func init() {
	RootCmd.AddCommand(daemonCmd)
	RootCmd.AddCommand(versionCmd)
	initDaemonCmd()
}
