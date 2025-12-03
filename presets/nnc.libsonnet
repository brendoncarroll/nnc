
local mount(dst, src) = {
  dst: dst,
  src: src,
};

local mountHostRO(dst, src) =
  mount(dst, {host_ro: src});

local mountHostRW(dst, src) =
  mount(dst, {host_rw: src});

local mountTmpfs(dst) =
  mount(dst, {tmpfs: {}});
  
local mountProcfs(dst="proc") =
  mount(dst, {procfs: {}});

local mountSysfs(dst="sys") =
  mount(dst, {sysfs: {}});

local mountDevtmpfs(dst="dev") =
  mount(dst, {devtmpfs: {}});

local mountsMerge(xs) =
  std.flattenArrays(xs);

local envSelectKeys(caller, keys) =
    std.map(function(k) k+"="+caller.envKV[k], keys);

local envMerge(xs) =
  local xs2 = std.map(
    function(x) if x == null then [] else x,
    xs
  );
  std.flattenArrays(xs2);

local netNone(name) = {
  "name": name,
  "backend": {},
};

local dataLit(path, contents, mode=420) =
  {
    "path": path,
    mode: 420,
    contents : {
      lit: contents,
    },
  };

local homeDir(caller) =
  local h = caller.envKV["HOME"];
  h;

local homePath(caller, p) =
  local hd = homeDir(caller);
  std.join("/", [hd, p]);

local mergeField(f, obj1, obj2) =
  if std.objectHas(obj1, f) && obj1[f] != null then obj1[f]
  else if std.objectHas(obj2, f) && obj2[f] != null then obj2[f]
  else null;

local merge2(a, b) =
  {
    mounts: std.flattenArrays([ x for x in [a.mounts, b.mounts] if x != null ]),
    env: std.flattenArrays([ x for x in [a.env, b.env] if x != null ]),
    net: std.flattenArrays([ x for x in [a.net, b.net] if x != null ]),
    data: std.flattenArrays([ x for x in [a.data, b.data] if x != null ]),
    main: mergeField("main", a, b),
  };

{
  mount :: mount,
  mountHostRO :: mountHostRO,
  mountHostRW :: mountHostRW,
  mountTmpfs :: mountTmpfs,
  mountDevtmpfs :: mountDevtmpfs,
  mountSysfs :: mountSysfs,
  mountProcfs :: mountProcfs,
  mountsMerge :: mountsMerge,

  netNone :: netNone,

  envMerge :: envMerge,
  envSelectKeys :: envSelectKeys,

  dataLit :: dataLit,

  homeDir :: homeDir,
  homePath :: homePath,

  merge2 :: merge2,
}
