package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/sirupsen/logrus"
)

// mostly swiped from https://github.com/vieux/docker-volume-sshfs/blob/master/main.go

const socketAddress = "/run/docker/plugins/volume.sock"

// Version is set from the go build commandline
var Version string

// CommitHash is set from the go build commandline
var CommitHash string

// BranchName is set from the go build commandline
var BranchName string

// Error log helper
func logError(format string, args ...interface{}) error {
	logrus.Errorf(format, args...)
	return fmt.Errorf(format, args...)
}

// TODO: detect what other versions of the plugin is running (locally and on other nodes)
// TODO: make sure we can access docker socket, and that we're actually at the plugin socket we think we are
// TODO: may need to figure out if installing as a swarm plugin also gives me access to the seaweedfs_internal network:
//       https://github.com/moby/moby/blob/master/integration/service/plugin_test.go#L109
func main() {
	debug := os.Getenv("DEBUG")
	if ok, _ := strconv.ParseBool(debug); ok {
		logrus.SetLevel(logrus.DebugLevel)
	}
	logrus.Infof("Version %s, build %s, branch: %s\n", Version, CommitHash, BranchName)
	d, err := newVolumeDriver(`{"root":"/mnt"}`)
	if err != nil {
		log.Fatal(err)
	}
	h := volume.NewHandler(d)
	logrus.Infof("listening on %s", socketAddress)

	logrus.Error(h.ServeUnix(socketAddress, 0))
}
