package ftp

import (
	"github.com/RenterRus/dwld-downloader/internal/repo/persistent"
	"github.com/RenterRus/dwld-downloader/internal/repo/temporary"
	v1 "github.com/RenterRus/dwld-ftp-sender/docs/proto/v1"
)

func NewFTPLoader(sender v1.SenderClient, sqlRepo persistent.SQLRepo, cache temporary.CacheRepo) *FTPSender {
	return &FTPSender{
		sender:  sender,
		sqlRepo: sqlRepo,
		cache:   cache,
	}
}
