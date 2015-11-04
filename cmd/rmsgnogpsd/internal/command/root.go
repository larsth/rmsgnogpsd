package command

import "github.com/spf13/cobra"

var (
	RootCmd = &cobra.Command{
		Use: CommandName,
		Short: CommandName + " is a daemon (service) that tells UNIX domain " +
			"clients about its GPS location.",
		Long: CommandName + " is a daemon (service) that tells UNIX domain " +
			"clients about its GPS location. The GPS location does only come " +
			"from a JSON document/file. The GPS location can only be changed " +
			"by changing the JSON document/file, and then restart the daemon " +
			"(service).",
	}
)
