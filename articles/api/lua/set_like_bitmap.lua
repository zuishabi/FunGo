-- Lua 脚本：ARGV[1]=bitmap（二进制字符串），ARGV[2]=id，ARGV[3]=uid（偏移，从0开始）
-- 返回：{createdFlag, prevBit}，createdFlag:1=新创建并写入传入bitmap，0=已存在；prevBit: 操作前该位的值（0或1）

local bitmap = ARGV[1]
local id = ARGV[2]
local uid = tonumber(ARGV[3])
if not uid then
  return redis.error_reply("invalid uid")
end

local key = "likebitmap-" .. id
local exists = redis.call("EXISTS", key)

if exists == 1 then
  -- key 存在：读取原位然后设置为 1
  local prev = redis.call("GETBIT", key, uid)
  redis.call("SETBIT", key, uid, 1)
  return {0, prev}
else
  -- key 不存在：先写入整段 bitmap（ARGV[1]），再设置对应位为 1
  -- ARGV[1] 可为二进制数据（由客户端以 raw bytes 传入）
  redis.call("SET", key, bitmap)
  redis.call("SETBIT", key, uid, 1)
  return {1, 0}
end
