
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

local homeDir(caller) =
  local h = caller.envKV["HOME"];
  h;

local homePath(caller, p) =
  local hd = homeDir(caller);
  std.join("/", [hd, p]);

{
  mount :: mount,
  mountHostRO :: mountHostRO,
  mountHostRW :: mountHostRW,
  mountTmpfs :: mountTmpfs,
  mountsMerge :: mountsMerge,

  envMerge :: envMerge,
  envSelectKeys:: envSelectKeys,

  homeDir :: homeDir,
  homePath :: homePath,
}
