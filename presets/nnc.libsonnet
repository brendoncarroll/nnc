
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

local mergeField(field, a, b, mergeFn) =
  if std.objectHas(a, field) then
    if std.objectHas(b, field) then mergeFn(a[field], b[field])
    else a[field]
  else if std.objectHas(b, field) then b[field]
  else null;

local merge2(a, b) =
  local mergeFn(a, b) = std.flattenArrays([x for x in [a, b] if x != null]);
  {
    mounts: mergeField("mounts", a, b, mergeFn),
    env: mergeField("env", a, b, mergeFn),
    net: mergeField("net", a, b, mergeFn),
    data: mergeField("data", a, b, mergeFn),
    main: mergeField("main", a, b, function(a, b) b),
  };
  
local merge(xs) =
  if std.length(xs) == 0 then {}
  else std.foldl(merge2, xs[1:], xs[0]);

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
  merge :: merge,
}
