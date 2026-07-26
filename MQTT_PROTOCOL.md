# iot-go-direct MQTT 消息协议

> 本文档描述服务器与网关之间通过 MQTT 协议交换的消息格式。不涉及 LoRa 射频层、Go 代码实现等内部细节。

---

## 1. 连接信息

| 参数 | 说明 |
|------|------|
| 协议 | MQTT 3.1.1 over TCP |
| 默认端口 | 1883 |
| QoS | 0（At-most-once） |
| Retain | false |
| ClientID | 同 `gateway_node_id`，12 位十六进制字符串（如 `F12309150001`） |
| 自动重连 | 开启，最大间隔 10 秒 |
| 认证方式 | 用户名 + 密码 |

---

## 2. 主题

每个网关使用一对主题：

| 主题 | 方向 | 用途 |
|------|------|------|
| `device/{gateway_node_id}/in` | 服务器 → 网关 | 服务器向网关下发请求 |
| `device/{gateway_node_id}/out` | 网关 → 服务器 | 网关向服务器上报状态或回复 |

> 例如网关 ID 为 `F12309150001`，则主题为 `device/F12309150001/in` 和 `device/F12309150001/out`。

---

## 3. 消息约定

### 3.1 编码

所有消息 Payload 为 **UTF-8 JSON 文本**。

### 3.2 消息类型区分

| 方向 | 标识字段 | 示例 |
|------|----------|------|
| 服务器 → 网关 | `method` | `{"method": "update_glass_color_request", ...}` |
| 网关 → 服务器 | `msg_type` | `{"msg_type": "update_glass_color_reply", ...}` |

### 3.3 通用字段

| 字段 | 类型 | 说明 |
|------|------|------|
| `gateway_node_id` | string | 12 位十六进制网关 ID |
| `node_id` | string | 8 位十六进制节点 ID |

---

## 4. 服务器 → 网关（下行请求）

### 4.1 消息列表

| method | 说明 | 回复 msg_type |
|--------|------|---------------|
| `config` | 下发/更新网关配置（节点列表、射频参数、心跳周期） | `config_reply` |
| `mqtt_config` | 修改 MQTT 连接参数（IP、端口、账号密码），修改后网关重启 | `mqtt_config_reply` |
| `update_glass_color_request` | 设置一个或多个节点的玻璃颜色 | `update_glass_color_reply` |
| `group_update_glass_color_request` | 按组批量设置颜色 | `group_update_glass_color_reply` |
| `broadcast_update_glass_color_request` | 广播设置所有节点同一颜色 | `broadcast_update_glass_color_reply` |
| `get_glass_status_request` | 查询单个节点当前玻璃颜色 | `glass_status_update`（异步） |
| `node_info_request` | 查询网关自身信息 | `node_info_reply` |
| `node_list_request` | 查询网关当前管理的节点列表 | `node_list_reply` |
| `gateway_upgrade_request` | 远程升级网关固件（通过 SFTP 下载） | `gateway_upgrade_reply` |
| `node_firmware_download_request` | 远程下载节点固件（透传） | `node_firmware_download_reply` |
| `node_firmware_upgrade_request` | 远程升级节点固件（透传） | `node_firmware_download_reply` |
| `set_touch_device_node_list_request` | 设置触摸设备关联的节点列表 | 无 |
| `gateway_reboot` | 远程重启网关 | `gateway_reboot_reply` |
| `execute_cmd` | 远程执行 shell 命令 | `remote_cmd_result`（异步） |

### 4.2 config

**用途**：配置网关的节点列表、LoRa 射频参数、心跳周期等。网关收到后会向所有节点发送初始化指令。

**请求格式**：

```json
{
  "method": "config",
  "params": {
    "node_list_1": ["12345678", "abcdef01"],
    "node_list_2": ["11111111"],
    "touch_node_list_1": [],
    "touch_node_list_2": [],
    "custom": "",
    "project": "",
    "heart_beat": 20,
    "heart_beat_to_server": 60,
    "module1": {
      "freq": 4734,
      "band": 250,
      "factor": 9
    },
    "module2": {
      "freq": 4463,
      "band": 250,
      "factor": 9
    }
  }
}
```

**字段说明**：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `node_list_1` | string[] | 是 | Module1 管理的节点 ID 列表，每项 8 位 hex |
| `node_list_2` | string[] | 是 | Module2 管理的节点 ID 列表，每项 8 位 hex |
| `touch_node_list_1` | string[] | 是 | Module1 触摸设备节点列表 |
| `touch_node_list_2` | string[] | 是 | Module2 触摸设备节点列表 |
| `custom` | string | 否 | 自定义字段 |
| `project` | string | 否 | 项目字段 |
| `heart_beat` | int | 是 | 节点心跳周期（秒） |
| `heart_beat_to_server` | int | 是 | 网关向服务器上报周期（秒），最小 30，须为 10 的倍数 |
| `module1` | object | 是 | Module1 射频参数：freq(3700-6000)、band(125/250/500)、factor(6-12) |
| `module2` | object | 是 | Module2 射频参数，同 module1 |

**回复**：

```json
{
  "msg_type": "config_reply",
  "gateway_node_id": "F12309150001"
}
```

### 4.3 update_glass_color_request

**用途**：设置一个或多个节点的玻璃颜色。支持单次请求携带多个节点。

**请求格式**：

```json
{
  "method": "update_glass_color_request",
  "params": [
    {
      "node_id": "12345678",
      "color": "1122334455667788"
    },
    {
      "node_id": "abcdef01",
      "color": "11ff22ff33ff44ff"
    }
  ]
}
```

**字段说明**：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `node_id` | string | 是 | 8 位十六进制节点 ID |
| `color` | string | 是 | 颜色编码，每 2 字符为一组：第 1 位为区域号（1-8），第 2 位为颜色值（0/2/4/F） |

**color 编码详解**：

`color` 字符串长度为偶数，每 2 个字符控制一个玻璃区域：

```
"1122334455667788"
││││││││││││││││
└┘└┘└┘└┘└┘└┘└┘
 1  2  3  4  5  6  7  8   ← 区域编号
```

| 字符对 | 区域 | 颜色值 |
|--------|------|--------|
| `10` | 区域 1 | 透明（0） |
| `22` | 区域 2 | 浅色（2） |
| `34` | 区域 3 | 深色（4） |
| `8f` | 区域 8 | 保留不变（f） |

颜色值范围仅 `0`、`2`、`4`、`F`，共 4 种（对应串口通信协议）：

| 数值 | 含义 |
|------|------|
| `0` | 透明 |
| `2` | 浅色 |
| `4` | 深色 |
| `F` | 保留原有颜色，不做变更 |

**回复**：

```json
{
  "msg_type": "update_glass_color_reply",
  "gateway_node_id": "F12309150001",
  "status": [
    {
      "node_id": "12345678",
      "color": "1122334455667788"
    },
    {
      "node_id": "abcdef01",
      "color": "11ff22ff33ff44ff"
    }
  ]
}
```

> **注意**：此回复仅表示网关已接收并将指令下发到节点，不代表节点已执行完毕。节点实际执行结果在后续的 `heartbeat_status_update` 中通过 `completion_status` 和 `color` 字段反映。

### 4.4 group_update_glass_color_request

**用途**：按预配置的节点分组批量设置颜色。无需逐个指定节点，直接引用 `node_list_1` 和 `node_list_2` 中的节点。

**请求格式**：

```json
{
  "method": "group_update_glass_color_request",
  "params": {
    "node_list1": ["12345678", "abcdef01"],
    "node_list2": ["11111111"],
    "color": "11ff22ff33ff44ff"
  }
}
```

**字段说明**：

| 字段 | 类型 | 必填 | 说明 |
|------|------|------|------|
| `node_list1` | string[] | 是 | 要调色的 Module1 节点列表（须已在 config 的 node_list_1 中） |
| `node_list2` | string[] | 是 | 要调色的 Module2 节点列表（须已在 config 的 node_list_2 中） |
| `color` | string | 是 | 统一颜色编码，规则同 `update_glass_color_request` |

**回复**：

```json
{
  "msg_type": "group_update_glass_color_reply",
  "gateway_node_id": "F12309150001",
  "invalid_nodes": []
}
```

| 字段 | 类型 | 说明 |
|------|------|------|
| `invalid_nodes` | string[] | 不在网关管理列表中的节点 ID |

### 4.5 broadcast_update_glass_color_request

**用途**：向所有节点广播同一颜色，不区分节点列表。

**请求格式**：

```json
{
  "method": "broadcast_update_glass_color_request",
  "params": {
    "color": "11ff22ff33ff44ff"
  }
}
```

**回复**：

```json
{
  "msg_type": "broadcast_update_glass_color_reply",
  "gateway_node_id": "F12309150001"
}
```

### 4.6 get_glass_status_request

**用途**：查询单个节点当前的玻璃颜色状态。

**请求格式**：

```json
{
  "method": "get_glass_status_request",
  "params": {
    "node_id": "12345678"
  }
}
```

**回复**：`glass_status_update`（异步上报）

```json
{
  "msg_type": "glass_status_update",
  "gateway_node_id": "F12309150001",
  "status": [
    {
      "node_id": "12345678",
      "color": "1122334455667788"
    }
  ]
}
```

### 4.7 node_info_request

**用途**：查询网关自身信息（版本、模块参数、MQTT 配置等）。

**请求格式**：

```json
{
  "method": "node_info_request",
  "params": {
    "gateway_node_id": "F12309150001"
  }
}
```

**回复**：

```json
{
  "msg_type": "node_info_reply",
  "gateway_node_id": "F12309150001",
  "hard_version": "1.0",
  "soft_version": "2.2",
  "custom": "test",
  "project": "test",
  "node_type": 1,
  "rssi": 10,
  "ccid": "test",
  "heart_beat": 20,
  "heart_beat_to_server": 60,
  "module1": {"freq": 4734, "band": 250, "factor": 9},
  "module2": {"freq": 4463, "band": 250, "factor": 9},
  "mqtt_ip": "app.kosglass.com",
  "mqtt_port": 1883,
  "mqtt_user_name": "l8juew73i2t17wavzthg",
  "mqtt_pwd": "i0eprmhypu3r16g3wuuc"
}
```

### 4.8 node_list_request

**用途**：查询网关当前管理的所有节点列表。

**请求格式**：

```json
{
  "method": "node_list_request"
}
```

**回复**：

```json
{
  "msg_type": "node_list_reply",
  "gateway_node_id": "F12309150001",
  "total_count": 3,
  "node_list_1": ["12345678", "abcdef01"],
  "node_list_2": ["11111111"],
  "touch_node_list_1": [],
  "touch_node_list_2": []
}
```

### 4.9 mqtt_config

**用途**：远程修改 MQTT 连接参数。修改后网关会自动重启以使用新配置。

**请求格式**：

```json
{
  "method": "mqtt_config",
  "params": {
    "mqtt_ip": "new.broker.com",
    "mqtt_port": 1883,
    "mqtt_user_name": "new_user",
    "mqtt_pwd": "new_pwd"
  }
}
```

**回复**：

```json
{
  "msg_type": "mqtt_config_reply"
}
```

### 4.10 gateway_upgrade_request

**用途**：远程升级网关固件。网关通过 SFTP 从指定服务器下载固件。

**请求格式**：

```json
{
  "method": "gateway_upgrade_request",
  "params": {
    "gateway_node_id": "F12309150001",
    "target_hardware_version": "1.0",
    "target_software_version": "2.3",
    "node_type": 1,
    "crc8": 123,
    "is_upload": 0,
    "ip": "sftp.example.com",
    "port": 22,
    "user": "user",
    "pwd": "pwd",
    "path": "/firmware/main-2.3"
  }
}
```

**升级条件**（网关端校验）：

- `gateway_node_id` 必须匹配
- `target_hardware_version` 必须匹配
- `target_software_version` 必须 **大于** 本地当前版本

**回复**：

```json
{
  "msg_type": "gateway_upgrade_reply",
  "gateway_node_id": "F12309150001",
  "target_hardware_version": "1.0",
  "target_software_version": "2.3",
  "node_type": 1,
  "crc8": 123,
  "is_upload": 0,
  "ip": "sftp.example.com",
  "port": 22,
  "user": "user",
  "pwd": "pwd",
  "path": "/firmware/main-2.3",
  "run_area": 0,
  "state": "OK"
}
```

| 字段 | 说明 |
|------|------|
| `state` | `"OK"` 表示已开始下载；`"error, ..."` 表示条件不满足或下载失败 |

### 4.11 node_firmware_download_request / node_firmware_upgrade_request

**用途**：节点固件下载/升级。参数透传给节点。

**请求格式**：

```json
{
  "method": "node_firmware_download_request",
  "params": {
    "target_hardware_version": "1.0",
    "target_software_version": "2.1",
    "node_type": 1,
    "target_run_area": 1,
    "crc8": 123,
    "ip": "sftp.example.com",
    "port": 22,
    "user": "user",
    "pwd": "pwd",
    "path": "/firmware/node-2.1.bin"
  }
}
```

**回复**：

```json
{
  "msg_type": "node_firmware_download_reply",
  "target_hardware_version": "1.0",
  "target_software_version": "2.1",
  "node_type": 1,
  "target_run_area": 1,
  "crc8": 123,
  "ip": "sftp.example.com",
  "port": 22,
  "user": "user",
  "pwd": "pwd",
  "path": "/firmware/node-2.1.bin",
  "state": "OK"
}
```

### 4.12 gateway_reboot

**用途**：远程重启网关。

**请求格式**：

```json
{
  "method": "gateway_reboot"
}
```

**回复**：

```json
{
  "msg_type": "gateway_reboot_reply",
  "gateway_node_id": "F12309150001"
}
```

### 4.13 execute_cmd

**用途**：远程执行 shell 命令。

**请求格式**：

```json
{
  "method": "execute_cmd",
  "params": {
    "cmd": "ls -l /tmp"
  }
}
```

**回复**（异步）：

```json
{
  "msg_type": "remote_cmd_result",
  "gateway_node_id": "F12309150001",
  "stdout": "total 8\ndrwxr-xr-x 2 root root 4096 ...",
  "stderr": ""
}
```

### 4.14 set_touch_device_node_list_request

**用途**：设置触摸设备关联的节点列表。（当前实现仅打印，无实际回复）

**请求格式**：

```json
{
  "method": "set_touch_device_node_list_request",
  "params": [
    {
      "node_id": "touch001",
      "node_ids": ["12345678", "abcdef01"]
    }
  ]
}
```

---

## 5. 网关 → 服务器（上行消息）

### 5.1 消息列表

| msg_type | 触发时机 | 说明 |
|----------|----------|------|
| `init` | 网关启动，MQTT 连接建立后 | 网关上线通知，携带完整配置信息 |
| `heartbeat_status_update` | 周期性（默认 60 秒） | 上报所有节点实时状态 |
| `config_reply` | 收到 `config` 后 | 配置已接收并开始执行 |
| `mqtt_config_reply` | 收到 `mqtt_config` 后 | MQTT 参数已更新，即将重启 |
| `update_glass_color_reply` | 收到 `update_glass_color_request` 后 | 指令已下发到节点 |
| `group_update_glass_color_reply` | 收到 `group_update_glass_color_request` 后 | 组调色指令已下发 |
| `broadcast_update_glass_color_reply` | 收到 `broadcast_update_glass_color_request` 后 | 广播调色指令已下发 |
| `glass_status_update` | 收到 `get_glass_status_request` 后 | 节点当前颜色状态 |
| `node_info_reply` | 收到 `node_info_request` 后 | 网关完整信息 |
| `node_list_reply` | 收到 `node_list_request` 后 | 节点列表 |
| `gateway_upgrade_reply` | 收到 `gateway_upgrade_request` 后 | 升级结果 |
| `node_firmware_download_reply` | 收到节点固件请求后 | 固件下载/升级结果 |
| `gateway_reboot_reply` | 收到 `gateway_reboot` 后 | 即将重启 |
| `remote_cmd_result` | 收到 `execute_cmd` 后 | 命令执行结果（异步） |
| `unknown_node` | 收到未配置节点的回复时 | 主动上报未知节点 |

### 5.2 init

**触发时机**：网关启动后，MQTT 连接建立并订阅成功时立即发送。

```json
{
  "msg_type": "init",
  "gateway_node_id": "F12309150001",
  "hard_version": "1.0",
  "soft_version": "2.2",
  "custom": "test",
  "project": "test",
  "node_type": 1,
  "rssi": 10,
  "ccid": "test",
  "heart_beat": 20,
  "heart_beat_to_server": 60,
  "module0": {"freq": 4723, "band": 250, "factor": 9},
  "module1": {"freq": 4734, "band": 250, "factor": 9},
  "module2": {"freq": 4463, "band": 250, "factor": 9}
}
```

### 5.3 heartbeat_status_update

**触发时机**：周期性上报（默认 60 秒）。这是服务器获取节点实时状态的主要途径。

```json
{
  "msg_type": "heartbeat_status_update",
  "gateway_node_id": "F12309150001",
  "status": [
    {
      "node_id": "12345678",
      "color": "11ff223344556677",
      "hard_version": "100",
      "soft_version": "101",
      "run_area": 1,
      "rssi": -45.5,
      "snr": 8.2,
      "is_offline": false,
      "completion_status": 2,
      "node_requesting_color": "11ff223344556677"
    }
  ]
}
```

**字段说明**：

| 字段 | 类型 | 说明 |
|------|------|------|
| `node_id` | string | 节点 ID |
| `color` | string | 节点**实际报告**的玻璃颜色（8 区域编码） |
| `node_requesting_color` | string | 服务器**最近一次请求**的颜色 |
| `hard_version` | string | 硬件版本（十六进制字符串） |
| `soft_version` | string | 软件版本（十六进制字符串） |
| `run_area` | int | 运行区域 |
| `rssi` | float | 信号强度（dBm） |
| `snr` | float | 信噪比 |
| `is_offline` | bool | 是否离线（心跳超时未响应） |
| `completion_status` | int | 命令执行状态：`0`=失败，`1`=执行中，`2`=已完成 |

> **颜色对比**：通过比较 `color`（实际颜色）与 `node_requesting_color`（请求颜色），服务器可判断节点是否已成功执行调色指令。

### 5.4 unknown_node

**触发时机**：网关收到未在 `config` 节点列表中的设备的回复时，主动上报。

```json
{
  "msg_type": "unknown_node",
  "unknown_node_list": ["99999999"]
}
```

---

## 6. 典型交互时序

### 6.1 网关上线

```
网关                                    服务器
 │ ──── init ───────────────────────────> │
 │ <────────────────────────── config ─── │
 │ ──── config_reply ───────────────────> │
 │                                        │
 │ ◄──────── 周期性 heartbeat_status_update ───────►
```

### 6.2 单节点调色

```
网关                                    服务器
 │ <──────── update_glass_color_request ─ │
 │ ──── update_glass_color_reply ───────> │
 │                                        │
 │ ◄──────── 后续 heartbeat_status_update ───────►
 │              （color / completion_status 反映执行结果）
```

### 6.3 远程重启

```
网关                                    服务器
 │ <──────── gateway_reboot ──────────── │
 │ ──── gateway_reboot_reply ───────────> │
 │ ──── [断开连接]                        │
 │ ──── [重新启动]                        │
 │ ──── init ───────────────────────────> │
```

---

## 7. 错误与边界行为

| 场景 | 行为 |
|------|------|
| JSON 解析失败 | 网关打印错误日志，丢弃该消息 |
| `config` / `update_glass_color_request` Schema 校验失败 | 丢弃消息，不回复 |
| `update_glass_color_request` 中 color 包含非法颜色值（非 0/2/4/F） | 该节点颜色被标记为无效（`SetColorForNodeAsInvalid`），直接加入回复 |
| `update_glass_color_request` 中 node_id 不在任何节点列表 | 该节点颜色被标记为无效，直接加入回复 |
| 节点心跳超时 | 在 `heartbeat_status_update` 中 `is_offline=true`，`color` 被置为无效值 |
| `group_update_glass_color` 节点未完成 | 心跳上报时发现 `color != node_requesting_color`，自动重发最多 5 次 |
| `gateway_upgrade` 版本条件不满足 | 回复 `state="error, upgrade condition not met..."`，不执行下载 |
| 同一 gateway_node_id 重复上线 | 后上线的实例会踢掉先上线的（MQTT ClientID 冲突） |

---

## 8. 安全提示

1. **明文传输**：当前使用 TCP 明文（端口 1883），未启用 TLS。生产环境建议配置 TLS。
2. **远程命令执行**：`execute_cmd` 可执行任意 shell 命令，建议服务端对命令做白名单限制。
3. **认证信息**：MQTT 用户名密码存储在本地配置文件 `iot_go.json` 中，注意文件权限保护。
