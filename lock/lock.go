package lock

import (
	"fmt"
	lockfile "github.com/ipfs/go-fs-lock"
	"io"
)

const LockFileName = "lock"

func TryLockDaemon(repoPath string) (io.Closer, error) {
	locked, err := lockfile.Locked(repoPath, LockFileName)
	if err != nil {
		return nil, err
	}
	if locked {
		return nil, fmt.Errorf("daemon already running")
	}
	f, err := lockfile.Lock("", "lock")
	if err != nil {
		return nil, err
	}
	return f, nil
}

func CheckLocked(repoPath string) (bool, error) {
	return lockfile.Locked(repoPath, LockFileName)
}
