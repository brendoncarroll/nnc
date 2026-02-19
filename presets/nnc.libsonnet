
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

local mountDev(name) =
  mount("dev/" + name, {host_dev: 0});

local mountsMerge(xs) =
  local xs2 = std.map(
    function(x) if x == null then [] else x,
    xs
  );
  std.flattenArrays(xs2);

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
    mode: mode,
    contents : {
      lit: contents,
    },
  };

local copyHostPath(path, hostPath, mode=420) =
  {
    "path": path,
    mode: mode,
    contents : {
      host_path: hostPath,
    },
  };

local homeDir(ctx) =
  local h = ctx.envKV["HOME"];
  h;

local homePath(ctx, p) =
  local hd = homeDir(ctx);
  std.join("/", [hd, p]);

local selectEnvKeys(ctx, keys) =
    std.map(function(k) k+"="+ctx.envKV[k], keys);

local applyAll(ctx, spec, presets) =
  std.foldl(function(acc, preset) preset(ctx, acc), presets, spec);

{
  mount :: mount,
  mountHostRO :: mountHostRO,
  mountHostRW :: mountHostRW,
  mountTmpfs :: mountTmpfs,
  mountDev :: mountDev,
  mountDevtmpfs :: mountDevtmpfs,
  mountSysfs :: mountSysfs,
  mountProcfs :: mountProcfs,
  mountsMerge :: mountsMerge,

  netNone :: netNone,

  envMerge :: envMerge,

  dataLit :: dataLit,
  copyHostPath :: copyHostPath,

  homeDir :: homeDir,
  homePath :: homePath,
  selectEnvKeys :: selectEnvKeys,

  applyAll :: applyAll,
}
