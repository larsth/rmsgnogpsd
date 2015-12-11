package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
	
	binmsg "github.com/larsth/go-rmsggpsbinmsg"

	"github.com/spf13/cobra"
)

var (
	daemonCmd = &cobra.Command{
		Use: "daemon",
		Short: "Sub-command 'daemon' starts this program as a daemon (service) " +
			"that tells TCP clients about its GPS location.",
		Long: "Sub-command 'daemon' starts this program as a " +
			"daemon (service) that tells TCP clients about its GPS " +
			"location. The GPS location does only come " +
			"from a JSON document/file. The GPS location can only be changed " +
			"by changing the JSON document/file, and then restart the daemon " +
			"(service).\n"+
			"A Timestamp (date and time) is set once, when the daemon starts.",
		Example: "nohup rmsgnogpsd daemon -c ./example_rmsgnogpsd_configuration_file.json",
		RunE: runDaemonE,
	}
)

func handleConnection(connChan chan *net.TCPConn) {
	var (
		conn    *net.TCPConn
		err     error
		buf     bytes.Buffer
		message []byte = make([]byte, 0, os.Getpagesize())
	)

	select {
	case conn = <-connChan:
		if message, err = json.Marshal(data.Payload); err != nil {
			conn.CloseWrite()
			log.Fatalln(err.Error())
		}
		if err = json.Indent(&buf, message, "", "\t"); err != nil {
			conn.CloseWrite()
			log.Fatalln(err.Error())
		}
		if _, err = buf.WriteTo(conn); err != nil {
			conn.CloseWrite()
			log.Println(err.Error())
		}
		conn.CloseWrite()
		buf.Reset()
		message = message[0:0:cap(message)] //reuse slice and it's backing array
		err = nil
		conn = nil
	}
}

func listenTCP(connChan chan *net.TCPConn) error {
	var (
		tcpConn *net.TCPConn
		err     error
	)
	//Listen for an inbound TCP connection:
	if tcpConn, err = data.Tcp.Listener.AcceptTCP(); err != nil {
		log.Fatalln(err.Error())
		return err
	}

	//Set some sane defaults ...
	//	We do not read!
	//	Creating a TCP connection is enough to trigger an answer ...
	if err = tcpConn.CloseRead(); err != nil {
		outErr := fmt.Errorf("INFO: %s", err.Error())
		log.Println(outErr.Error())
		return outErr
	}
	//	Use TCP keep alive ...
	if err = tcpConn.SetKeepAlive(true); err != nil {
		outErr := fmt.Errorf("INFO: %s", err.Error())
		log.Println(outErr.Error())
		return outErr
	}
	// ... for a period of up to 30 minutes ...
	if err = tcpConn.SetKeepAlivePeriod(time.Minute * 30); err != nil {
		outErr := fmt.Errorf("INFO: %s", err.Error())
		log.Println(outErr.Error())
		return outErr
	}
	//Use Nagle's algorithm (=send fewer TCP packets)...
	if err = tcpConn.SetNoDelay(false); err != nil {
		outErr := fmt.Errorf("INFO: %s", err.Error())
		log.Println(outErr.Error())
		return outErr
	}

	//Defer the real work to a TCP worker go routine:
	connChan <- tcpConn

	return nil
}

func runDaemonE(_ *cobra.Command, _ []string) error {
	var (
		configFile *os.File
		err        error
		decoder    *json.Decoder
		workers, i uint

		connChan chan *net.TCPConn
	)

	//Open the JSON configuration fiel:
	if configFile, err = os.Open(flags.ConfigFileName); err != nil {
		log.Println(err.Error())
		return err
	}

	//Parse the JSON document (configuration file):
	decoder = json.NewDecoder(configFile)
	if err = decoder.Decode(&data.Config); err != nil {
		log.Println(err.Error())
		configFile.Close()
		return err
	}

	//File and JSON ressources clean-up:
	decoder = nil
	configFile.Close()
	configFile = nil

	//Choose which host string to use:
	//	if len(data.Config.Host) > 0 {
	data.Host = data.Config.Host
	//	} else {
	//Host string:
	//data.Host = fmt.Sprintf("%s:%d", DefaultHost, flags.Port)
	//}
	
	//Create binmsg.Payload structure.
	//Don't use *Binary methods on Paylaod structure
	// - won't work because we initialize Payload without a HMACKey and a salt.
	data.Payload = &binmsg.Payload{}
	
	//Copy data from 'data.Config.Gps' to the Payload:
	data.Payload.Message.Gps.SetAlt(data.Config.Gps.Alt)
	data.Payload.Message.Gps.SetLat(data.Config.Gps.Lat)
	data.Payload.Message.Gps.SetLon(data.Config.Gps.Lon)
	
	//Set the TimeStamp
	data.Payload.Message.TimeStamp.Time = time.Now().UTC()
	//Set FixMode
	if data.Config.Gps.Alt == float64(0.0) {
		data.Payload.Message.Gps.FixMode = binmsg.Fix2D
	} else {
		data.Payload.Message.Gps.FixMode = binmsg.Fix3D
	}

	//Create the *net.TCPAddr:
	if data.Tcp.Addr, err = net.ResolveTCPAddr("tcp6", data.Host); err != nil {
		log.Println(err.Error())
		return err
	}

	//Create the TCP listener:
	data.Tcp.Listener, err = net.ListenTCP("tcp6", data.Tcp.Addr)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	//Create the connection channel:
	connChan = make(chan *net.TCPConn, 0)
	
	if data.Config.Workers != nil {
		if (*data.Config.Workers) > 0 {
			workers = (*data.Config.Workers)
		} else {
			workers = defaultClientConnectionWorkers
		}
	} else {
		workers = defaultClientConnectionWorkers
	}
	//Start 'workers' number of TCP client connection workers:
	for i = 1; i <= workers; i++ {
		go handleConnection(connChan)
	}

	//Listen for an inbound TCP connection:
	for {
		if err = listenTCP(connChan); err != nil {
			outErr := fmt.Errorf("INFO: %s", err.Error())
			log.Println(outErr.Error())
			return outErr
		}
		//rinse, repeat, ...
	}

	return nil
}

func initDaemonCmd() {
	var pf = daemonCmd.PersistentFlags()
	pf.StringVarP(&flags.ConfigFileName, "config", "c", DefaultConfigFileName,
		"The name of the JSON configuration file.")
	//	pf.Uint16VarP(&flags.Port, "port", "p", DefaultPort,
	//		"TCP port number, range: 1024 - 65535.")
}
