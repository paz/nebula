package nebula

import (
	"fmt"
	"io"
	"net"
	"os"
)

type Tun struct {
	Inside
}

func newTunFromFd(l *logrus.Logger, deviceFd int, cidr *net.IPNet, defaultMTU int, routes []route, unsafeRoutes []route, txQueueLen int) (ifce *Tun, err error) {
	return nil, fmt.Errorf("newTunFromFd not supported in Windows")
}

func newTun(l *logrus.Logger, deviceName string, cidr *net.IPNet, defaultMTU int, routes []route, unsafeRoutes []route, txQueueLen int, multiqueue bool) (ifce *Tun, err error) {
	if len(routes) > 0 {
		return nil, fmt.Errorf("route MTU not supported in Windows")
	}

	useWintun := true
	if err = checkWinTunExists(); err != nil {
		l.WithError(err).Warn("Check Wintun driver failed, fallback to wintap driver")
		useWintun = false
	}

	var inside Inside
	if useWintun {
		inside, err = newWinTun(deviceName, cidr, defaultMTU, unsafeRoutes, txQueueLen)
		if err != nil {
			return nil, fmt.Errorf("Create Wintun interface failed, %w", err)
		}
	} else {
		inside, err = newWindowsWaterTun(deviceName, cidr, defaultMTU, unsafeRoutes, txQueueLen)
		if err != nil {
			return nil, fmt.Errorf("Create wintap driver failed, %w", err)
		}
	}

	return &Tun{
		Inside: inside,
	}, nil
}

func checkWinTunExists() error {
	_, err := os.Stat("wintun.dll")
	return err
}

func (t *Tun) NewMultiQueueReader() (io.ReadWriteCloser, error) {
	return nil, fmt.Errorf("TODO: multiqueue not implemented for windows")
}
