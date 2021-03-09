package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/docker/go-plugins-helpers/volume"
	"github.com/sirupsen/logrus"
)

func newVolumeDriver(config string) (*volumeDriver, error) {
	var d *volumeDriver
	json.Unmarshal([]byte(config), &d)
	logrus.WithField("method", "new driver").Debug(d.root)
	return d, nil
}

// Get the list of capabilities the driver supports.
// The driver is not required to implement Capabilities. If it is not implemented, the default values are used.
func (d *volumeDriver) Capabilities() *volume.CapabilitiesResponse {
	logrus.WithField("method", "capabilities").Debugf("version %s, build %s, branch: %s\n", Version, CommitHash, BranchName)
	return &volume.CapabilitiesResponse{Capabilities: volume.Capability{Scope: "global"}}
}

// Create Instructs the plugin that the user wants to create a volume,
// given a user specified volume name. The plugin does not need to actually
// manifest the volume on the filesystem yet (until Mount is called).
// Opts is a map of driver specific options passed through from the user request.
func (d *volumeDriver) Create(r *volume.CreateRequest) error {
	logrus.WithField("method", "create").Debugf("%#v", r)
	var v dockerVolume
	for key, val := range r.Options {
		switch key {
		default:
			if val != "" {
				v.Options = append(v.Options, key+"="+val)
			} else {
				v.Options = append(v.Options, key)
			}
		}
	}
	v.Mountpoint = filepath.Join(d.root, r.Name) // "/path/under/PropagatedMount"
	v.Name = r.Name
	if err := updateVolume(&v); err != nil {
		return err
	}
	return nil
}

// Get info about volume_name.
func (d *volumeDriver) Get(r *volume.GetRequest) (*volume.GetResponse, error) {
	logrus.WithField("method", "get").Debugf("%#v", r)
	v, err := getVolumeByName(r.Name)
	if err != nil {
		return &volume.GetResponse{}, logError("volume %s not found", r.Name)
	}
	logrus.WithField("get", "volumeinfo").Debugf("%#v", v)
	return &volume.GetResponse{Volume: &volume.Volume{
		Name:       r.Name,
		Mountpoint: v.Mountpoint, // "/path/under/PropogatedMount"
	}}, nil
}

// List of volumes registered with the plugin.
func (d *volumeDriver) List() (*volume.ListResponse, error) {
	logrus.WithField("method", "list").Debugf("version %s, build %s, branch: %s\n", Version, CommitHash, BranchName)
	var vols = listVolumes()
	return &volume.ListResponse{Volumes: vols}, nil
}

// Mount is called once per container start.
// If the same volume_name is requested more than once, the plugin may need to keep
// track of each new mount request and provision at the first mount request and
// deprovision at the last corresponding unmount request.
// Docker requires the plugin to provide a volume, given a user specified volume name.
// ID is a unique ID for the caller that is requesting the mount.
func (d *volumeDriver) Mount(r *volume.MountRequest) (*volume.MountResponse, error) {
	logrus.WithField("method", "mount").Debugf("%#v", r)
	v, _ := getVolumeByName(r.Name)
	mountVolume(v)
	return &volume.MountResponse{Mountpoint: v.Mountpoint}, nil
}

// Path requests the path to the volume with the given volume_name.
func (d *volumeDriver) Path(r *volume.PathRequest) (*volume.PathResponse, error) {
	logrus.WithField("method", "path").Debugf("%#v", r)
	v, err := getVolumeByName(r.Name)
	if err != nil {
		return &volume.PathResponse{}, logError("volume %s not found", r.Name)
	}
	return &volume.PathResponse{Mountpoint: v.Mountpoint}, nil
}

// Remove the specified volume from disk. This request is issued when a
// user invokes docker rm -v to remove volumes associated with a container.
func (d *volumeDriver) Remove(r *volume.RemoveRequest) error {
	logrus.WithField("method", "remove").Debugf("%#v", r)
	v, err := getVolumeByName(r.Name)
	if err != nil {
		return logError("volume %s not found", r.Name)
	}
	if v.connections != 0 {
		return logError("volume %s is currently used by a container", r.Name)
	}
	unmountVolume(v)
	if err := os.RemoveAll(v.Mountpoint); err != nil {
		logError(err.Error())
	}
	removeVolume(v)
	return nil
}

// Docker is no longer using the named volume.
// Unmount is called once per container stop.
// Plugin may deduce that it is safe to deprovision the volume at this point.
// ID is a unique ID for the caller that is requesting the mount.
func (d *volumeDriver) Unmount(r *volume.UnmountRequest) error {
	logrus.WithField("method", "unmount").Debugf("%#v", r)
	v, err := getVolumeByName(r.Name)
	if err != nil {
		return logError("volume %s not found", r.Name)
	}
	v.connections--
	err = updateVolume(v)
	if err = updateVolume(v); err != nil {
		logrus.WithField("updateVolume ERROR", err).Errorf("%#v", v)
	} else {
		logrus.WithField("updateVolume", r.Name).Debugf("%#v", v)
	}
	if v.connections <= 0 {
		v.connections = 0
		updateVolume(v)
		unmountVolume(v)
	}
	return nil
}
