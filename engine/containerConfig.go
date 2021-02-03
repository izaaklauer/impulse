package engine

import (
    "github.com/opencontainers/runc/libcontainer/configs"
    "github.com/opencontainers/runc/libcontainer/specconv"
    "golang.org/x/sys/unix"
)

func containerConfig(rootfsPath string) *configs.Config {
    defaultMountFlags := unix.MS_NOEXEC | unix.MS_NOSUID | unix.MS_NODEV
    config := &configs.Config{
        Rootfs: rootfsPath,
        Capabilities: &configs.Capabilities{
            Bounding: []string{
                "CAP_AUDIT_WRITE",
                "CAP_KILL",
                "CAP_NET_BIND_SERVICE",
            },
            Effective: []string{
                "CAP_AUDIT_WRITE",
                "CAP_KILL",
                "CAP_NET_BIND_SERVICE",
            },
            Inheritable: []string{
                "CAP_AUDIT_WRITE",
                "CAP_KILL",
                "CAP_NET_BIND_SERVICE",
            },
            Permitted: []string{
                "CAP_AUDIT_WRITE",
                "CAP_KILL",
                "CAP_NET_BIND_SERVICE",
            },
            Ambient: []string{
                "CAP_AUDIT_WRITE",
                "CAP_KILL",
                "CAP_NET_BIND_SERVICE",
            },
        },
        Namespaces: configs.Namespaces([]configs.Namespace{
            {Type: configs.NEWNS},
            {Type: configs.NEWUTS},
            {Type: configs.NEWIPC},
            {Type: configs.NEWPID},
            {Type: configs.NEWUSER},
            {Type: configs.NEWNET},
            {Type: configs.NEWCGROUP},
        }),
        //Cgroups: &configs.Cgroup{
        //    Name:   "test-container",
        //    Parent: "system",
        //    Resources: &configs.Resources{
        //        MemorySwappiness: nil,
        //        Devices:          specconv.AllowedDevices,
        //    },
        //},
        MaskPaths: []string{
            "/proc/kcore",
            "/sys/firmware",
        },
        ReadonlyPaths: []string{
            "/proc/sys", "/proc/sysrq-trigger", "/proc/irq", "/proc/bus",
        },
        Devices:  specconv.AllowedDevices,
        Hostname: "testing",
        Mounts: []*configs.Mount{
            {
                Source:      "proc",
                Destination: "/proc",
                Device:      "proc",
                Flags:       defaultMountFlags,
            },
            {
                Source:      "tmpfs",
                Destination: "/dev",
                Device:      "tmpfs",
                Flags:       unix.MS_NOSUID | unix.MS_STRICTATIME,
                Data:        "mode=755",
            },
            {
                Source:      "devpts",
                Destination: "/dev/pts",
                Device:      "devpts",
                Flags:       unix.MS_NOSUID | unix.MS_NOEXEC,
                Data:        "newinstance,ptmxmode=0666,mode=0620,gid=5",
            },
            {
                Device:      "tmpfs",
                Source:      "shm",
                Destination: "/dev/shm",
                Data:        "mode=1777,size=65536k",
                Flags:       defaultMountFlags,
            },
            {
                Source:      "mqueue",
                Destination: "/dev/mqueue",
                Device:      "mqueue",
                Flags:       defaultMountFlags,
            },
            {
                Source:      "sysfs",
                Destination: "/sys",
                Device:      "sysfs",
                Flags:       defaultMountFlags | unix.MS_RDONLY,
            },
        },
        //UidMappings: []configs.IDMap{
        //    {
        //        ContainerID: 0,
        //        HostID: 1000,
        //        Size: 65536,
        //    },
        //},
        //GidMappings: []configs.IDMap{
        //    {
        //        ContainerID: 0,
        //        HostID: 1000,
        //        Size: 65536,
        //    },
        //},
        Networks: []*configs.Network{
            {
                Type:    "loopback",
                Address: "127.0.0.1/0",
                Gateway: "localhost",
            },
        },
        Rlimits: []configs.Rlimit{
            {
                Type: unix.RLIMIT_NOFILE,
                Hard: uint64(1025),
                Soft: uint64(1025),
            },
        },
    }
    return config
}
