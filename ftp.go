package goconf

import (
	"github.com/jkrecek/goftp"
	"math"
	"time"
)

type FtpConnection struct {
	Address string
	User    string
	Pass    string
	Path    string
	Timeout float64
}

func (fc *FtpConnection) getTimeout() time.Duration {
	return time.Duration(math.Max(1, fc.Timeout)) * time.Second
}

func (fc *FtpConnection) GetConnection() (ftp *goftp.FTP, err error) {
	if fc.Timeout == 0 {
		ftp, err = goftp.Connect(fc.Address)
	} else {
		ftp, err = goftp.ConnectTimeout(fc.Address, fc.getTimeout())
	}

	if err != nil {
		return
	}

	if err = ftp.Login(fc.User, fc.Pass); err != nil {
		return
	}

	err = ftp.Cwd(fc.Path)

	return
}

func (fc *FtpConnection) RuntimeTest() (err error, fatal bool) {
	conn, err := fc.GetConnection()
	if err != nil {
		fatal = true
		return
	}

	conn.Close()
	return
}
