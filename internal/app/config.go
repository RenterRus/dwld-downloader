package app

type Server struct {
	Host   string
	Port   int
	Enable bool
}

type FTPClient struct {
	Addr            Server
	User            string
	Pass            string
	RemoteDirectory string
}

type Config struct {
	GRPC Server
	HTTP Server
	FTP  FTPClient

	PathToDB string
	Cache    Server

	EagerMode bool
}
