local nnc = import "./nnc.libsonnet";

function(ctx, spec)
  spec + {
    data: nnc.mountsMerge([spec.data,
      [
        nnc.dataLit("/etc/passwd", "root:x:0:0:root:/root:/bin/sh\n", 420),
        nnc.dataLit("/etc/group", "root:x:0:\n", 420),
      ],
    ]),
  }
