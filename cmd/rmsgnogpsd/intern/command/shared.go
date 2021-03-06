package command

import (
	"net"

	binmsg "github.com/larsth/go-rmsggpsbinmsg"
)

type Gps struct {
	Lat float64 `json:"latitude"`
	Lon float64 `json:"longitude"`
	Alt float64 `json:"altitude"`
}

type Config struct {
	Workers *uint  `json: workers,omitempty`
	Host    string `json:"tcp6host,omitempty"`
	Gps     Gps    `json:"wgs84_nogps"`
}

type Flags struct {
	ConfigFileName string
	//Port           uint16
}

type TcpData struct {
	Addr     *net.TCPAddr
	Listener *net.TCPListener
}

type Data struct {
	Payload *binmsg.Payload
	Host    string
	Config  Config
	Tcp     TcpData
}

const (
	CommandName = "rmsgnogpsd"
	//DefaultPort uint16 = 10001
	DefaultHost                    = "[::1]"
	defaultClientConnectionWorkers = 16
)

var (
	data  Data
	flags Flags
)
