package main

import (
	"github.com/docker/go-plugins-helpers/volume"
)

type dockerVolume struct {
	Options          []string
	Name, Mountpoint string
	connections      int
}

func (d *volumeDriver) listVolumes() []*volume.Volume {
	var volumes []*volume.Volume
	for _, mount := range d.volumes {
		var v volume.Volume
		v.Name = mount.Name
		v.Mountpoint = mount.Mountpoint
		volumes = append(volumes, &v)
	}
	return volumes
}

func (d *volumeDriver) mountVolume(v *dockerVolume) error {
	return nil
}

func (d *volumeDriver) removeVolume(v *dockerVolume) error {
	delete(d.volumes, v.Name)
	return nil
}

func (d *volumeDriver) unmountVolume(v *dockerVolume) error {
	return nil
}

func (d *volumeDriver) updateVolume(v *dockerVolume) error {
	d.volumes[v.Name] = v
	return nil
}
