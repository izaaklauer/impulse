package engine

import (
    "github.com/opencontainers/runtime-spec/specs-go"
)


func getContainerSpec(id string, rootfsPath string) specs.Spec {
    
    capabilities := []string{
        "CAP_AUDIT_WRITE",
        "CAP_KILL",
        "CAP_NET_BIND_SERVICE",
    }
    
    spec := specs.Spec {
        Hostname: id,
        Root: &specs.Root{
            Path: rootfsPath,
        },
        Process: &specs.Process{
            Terminal: true,
            User: specs.User{
                UID: 0,
                GID: 0,
                AdditionalGids: []uint32{0, 1, 2, 3, 4, 6, 10, 11, 20, 26, 27},
            },
            Cwd: "/",
            Capabilities: &specs.LinuxCapabilities{
                Bounding: capabilities,
                Effective: capabilities,
                Inheritable: capabilities,
                Permitted: capabilities,
                Ambient: capabilities,
            },
            Rlimits: []specs.POSIXRlimit{
                {
                    Type: "RLIMIT_NOFILE",
                    Hard: 1024,
                    Soft: 1024,
                },
            },
            NoNewPrivileges: true,
        },
        Linux: &specs.Linux{
            Namespaces: []specs.LinuxNamespace{
                {Type: specs.PIDNamespace},
                {Type: specs.IPCNamespace},
                {Type: specs.UTSNamespace},
                {Type: specs.MountNamespace},
                //{Type: specs.NetworkNamespace},
            },
            Resources: &specs.LinuxResources{
                Devices: []specs.LinuxDeviceCgroup{
                    {
                        Allow: false,
                        Access: "rwm",
                    },
                },
            },
            MaskedPaths: []string {
                "/proc/kcore",
                "/proc/latency_stats",
                "/proc/timer_list",
                "/proc/timer_stats",
                "/proc/sched_debug",
                "/sys/firmware",
                "/proc/scsi",
            },
            ReadonlyPaths: []string{
                "/proc/asound",
                "/proc/bus",
                "/proc/fs",
                "/proc/irq",
                "/proc/sys",
                "/proc/sysrq-trigger",
            },
        },
        Mounts: []specs.Mount{
            {
                Destination: "/proc",
                Type: "proc",
                Source: "proc",
            },
            {
                Destination: "/dev",
                Type: "tmpfs",
                Source: "tmpfs",
                Options: []string {
                    "nosuid", "strictatime", "mode=755", "size=65536k", 
                },
            },
            {
                Destination: "/dev/pts",
                Type: "devpts",
                Source: "devpts",
                Options: []string {
                    "nosuid", "noexec", "newinstance", "ptmxmode=0666", "mode=0620", "gid=5",
                },
            },
            {
                Destination: "/dev/shm",
                Type: "tmpfs",
                Source: "shm",
                Options: []string {
                    "nosuid", "noexec", "nodev", "mode=1777", "size=65536k",
                },
            },
            {
                Destination: "/dev/mqueue",
                Type: "mqueue",
                Source: "mqueue",
                Options: []string {
                    "nosuid", "noexec", "nodev",
                },
            },
            {
                Destination: "/sys",
                Type: "sysfs",
                Source: "sysfs",
                Options: []string {
                    "nosuid", "noexec", "nodev", "ro",
                },
            },
            {
                Destination: "/sys/fs/cgroup",
                Type: "cgroup",
                Source: "cgroup",
                Options: []string {
                    "nosuid", "noexec", "nodev","relatime", "ro",
                },
            },
        },
    }
    
    // Python specific
    spec.Process.Env = []string {
        "PATH=/usr/local/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin",
        "TERM=xterm",
        "LANG=C.UTF-8",
        "GPG_KEY=E3FF2839C048B25C084DEBE9B26995E310250568",
        "PYTHON_VERSION=3.9.1",
        "PYTHON_PIP_VERSION=21.0.1",
        "PYTHON_GET_PIP_URL=https://github.com/pypa/get-pip/raw/4be3fe44ad9dedc028629ed1497052d65d281b8e/get-pip.py",
        "PYTHON_GET_PIP_SHA256=8006625804f55e1bd99ad4214fd07082fee27a1c35945648a58f9087a714e9d4",
        "HOME=/root",
    }
    
    spec.Process.Args = []string{"env"}
    
    return spec
}
