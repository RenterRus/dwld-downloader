package ftp

import (
	"github.com/RenterRus/dwld-downloader/internal/repo/persistent"
	"github.com/RenterRus/dwld-downloader/internal/repo/temporary"
)

type FTPSenderConf struct {
	Host       string
	User       string
	Pass       string
	LocalPath  string
	RemotePath string
	Port       int
	Enable     bool
	SqlRepo    persistent.SQLRepo
	Cache      temporary.CacheRepo
}

func NewFTPSender(conf *FTPSenderConf) Sender {
	return &FTPSender{
		Host:       conf.Host,
		User:       conf.User,
		Pass:       conf.Pass,
		LocalPath:  conf.LocalPath,
		RemotePath: conf.RemotePath,
		Port:       conf.Port,
		Enable:     conf.Enable,
		sqlRepo:    conf.SqlRepo,
		cache:      conf.Cache,
	}
}
