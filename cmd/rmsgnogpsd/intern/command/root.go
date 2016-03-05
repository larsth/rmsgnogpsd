package command

import (
	"strconv"

	"github.com/spf13/cobra"
)

var (
	RootCmd = &cobra.Command{
		Use: CommandName,
		Short: CommandName + " is a daemon (service) that listens for " +
			"inbound TCP connections at the port number specified in " +
			"the configuration file, and tell TCP clients about its GPS " +
			"location by sending a JSON document as the response.",
		Long: CommandName + " is a daemon (service) that listens for " +
			"incoming TCP connection at the port number specified in " +
			"the configuration file, and tell TCP clients about its GPS " +
			"location by sending a JSON document as the response." +
			"The GPS location can only be changed by changing the JSON " +
			"configuration document/file, and then restart the daemon" +
			"(service).\n" +
			"If the amount of workers is not set or is set to zero (0), then " +
			strconv.FormatUint(uint64(defaultClientConnectionWorkers), 10) +
			" TCP client connection workers will be used.",
	}
)
