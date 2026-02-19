# No Nonsense Containers

No Nonsense Containers (NNC) is a CLI tool for
running programs on Linux using the kernel's isolation primitives
to enforce the principle of least authority.

The example below would launch the fish shell with several files and directories passed through, and HOME set to /root inside the container.
```
nnc run \
  --mro /usr:/usr \
  --mro /bin:/bin \
  --mro /lib64:/lib64 \
  --mro /root/.config/fish:$HOME/.config/fish \
  --env HOME=/root \
  --mro "/x:$(pwd)" \
  /bin/fish
```

Unlike with other containerization tools, like Docker,
it is always assumed that resources should *not* be
passed on by default.
This means that the user must understand which resources
are required by the application to provide the desired functionality, and tell NNC to pass them through.
NNC tries to make this as easy as possible by using Presets to transfer resources from parent process to child process.

## Presets
nnc configuration is focused around command line flags.
This means that people who don't want to learn anything new can
still use all nnc features, by scripting in their shell as they would normally.
For example, look at the `etc/fish.sh` file included in this repository.

nnc also offers a more powerful way to configure a container using Presets.
A preset is a pure function:
```
(ctx: CallerContext, spec: ContainerSpec) -> ContainerSpec
```
Presets are defined using Jsonnet, and can be layered one after another.
Presets are applied in the order that they are passed to the CLI.

Presets must be defined in files with the `.jsonnet` extension, and the extension is omitted when referring to them on the command line.

The caller context `ctx` contains all the resources that the parent process has access to, represented as Jsonnet.
The `spec` is the specification for how to create the child process.
The default spec has no resources: no files, no network, nothing.
Each Preset can edit the spec by copying resources from the caller context, or by removing resources already in the spec.
The edited spec is returned, and passed as input to a subsequent Preset, if there is one.
Once all the presets have been applied, the child process is created according to the spec.

The presets directory in this repository contains presets for common applications, and users are encouraged
to contribute Presets.

## Done/Goals/Non-Goals

### Done
- Secure-by-default process isolation.
- Mount a read-only file or directory from the host
- Mount a read-write file or directory from the host
- Set environment variables in the container
- Pass arguments to the container process
- Give the container network access, as if it was a process on the host.

### Goals
- Cause the container to appear as another computer on the network.
- Force the container to communicate through WireGuard. 
- CPU and Memory constraints
- Add Presets for common applications, creating a standard library.

### Non-Goals
- Users and Groups. The child process always runs as root inside, and the parent user outside. 
- All possible mount options
- All possible network interfaces
- Images
- Daemons
- Pluggable runtimes and other nonsense

## Examples

This example uses the fish shell.
Using the `fish.jsonnet` preset available in this repository,
you can enter a shell with reduced capabilities.
```shell
nnc run /bin/fish --preset ./presets/fish
```
