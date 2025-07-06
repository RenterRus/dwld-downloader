package ftp

type FTPSenderConf struct {
	Host       string
	User       string
	Pass       string
	LocalPath  string
	RemotePath string
	Port       int
}

func NewFTPSender(conf *FTPSenderConf) Sender {
	return &FTPSender{
		Host:       conf.Host,
		User:       conf.User,
		Pass:       conf.Pass,
		LocalPath:  conf.LocalPath,
		RemotePath: conf.RemotePath,
		Port:       conf.Port,
	}
}
