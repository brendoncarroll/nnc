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

## Examples

This example uses the fish shell.
Using the `fish.jsonnet` preset available in this repository,
you can enter a shell with reduced capabilities.
```shell
nnc run /bin/fish --preset ./presets/fish
```
