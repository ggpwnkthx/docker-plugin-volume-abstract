package main

import (
	"path/filepath"

	"github.com/docker/go-plugins-helpers/volume"
)

type dockerVolume struct {
	Options          []string
	Name, Mountpoint string
	connections      int
}

type volumeDriver struct {
	root string
}

func getVolumeByName(name string) (*dockerVolume, error) {
	var v = dockerVolume{
		Options:     nil,
		Name:        name,
		Mountpoint:  filepath.Join("/mnt", name),
		connections: 0,
	}
	return &v, nil
}

func listVolumes() []*volume.Volume {
	var volumes []*volume.Volume
	return volumes
}

func mountVolume(v *dockerVolume) error {
	return nil
}

func removeVolume(v *dockerVolume) error {
	return nil
}

func unmountVolume(v *dockerVolume) error {
	return nil
}

func updateVolume(v *dockerVolume) error {
	return nil
}
