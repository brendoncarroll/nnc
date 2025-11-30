# No Nonsense Containers

No Nonsense Containers (NNC) is a CLI tool for
running programs on Linux using the kernel's isolation primitives
to enforce the principle of least authority.

Unlike with other containerization tools, like Docker,
it is always assumed that permissions should *not* be
passed on by default.
This means that the user must understand which resources
are required to provide their desired functionality.
If they forget to or don't know how to provide the subprocess
with the needed resources, then the application will not work as intended.
That will lead them to understand which resources are necessary so they can pass them through.


## Presets
nnc configuration is focused around command line flags.
This means that people who don't want to learn anything new can
still use all nnc features, by scripting in their shell as they would normally.
For example, look at the `etc/fish.sh` file included in this repository.

nnc also offers a more powerful way to configure a container using presets.
A preset is a function `(spec: ContainerSpec, caller: CallerContext) -> ContainerSpec
Presets are defined using Jsonnet, and can be layered on top of one another.
Presets are applied in the order that they are passed to the CLI.

Presets must be defined in `.jsonnet` files, and the extension is omitted when referring to them on the command line.
However, when imported in a Jsonnet program e.g. 'local lib = import "other_preset.jsonnet"`, the extension must be included.

## Examples

This example uses the fish shell.
Using the `fish.jsonnet` preset available in this repository,
you can enter a shell with reduced capabilities.
```shell
nnc run /bin/fish --preset ./presets/fish
```
