# iot-go-direct MQTT 协议总结

> 基于 Go 实现的物联网玻璃调色网关，通过 MQTT 协议与服务器通信，通过 LoRa 协议与终端节点通信。

- **Repo**: `weijia/iot-go-direct`
- **App Version**: `2.2`
- **Go**: `1.22.2`
- **MQTT Lib**: `github.com/eclipse/paho.mqtt.golang v1.4.3`

---

## 目录

1. [协议总览](#1-协议总览)
2. [连接配置](#2-连接配置)
3. [主题结构](#3-主题结构)
4. [消息格式](#4-消息格式)
5. [收发流程](#5-收发流程)
6. [服务器 → 网关（下行指令）](#6-服务器--网关下行指令)
7. [网关 → 服务器（上行消息）](#7-网关--服务器上行消息)
8. [与 LoRa 协议的桥接](#8-与-lora-协议的桥接)
9. [完整交互示例](#9-完整交互示例)
10. [关键常量](#10-关键常量)

---

## 1. 协议总览

该网关是一个 **LoRa-MQTT 桥接设备**：服务器通过 **MQTT** 下发控制指令到网关，网关将其转译为 **LoRa 二进制帧** 发送给终端节点（玻璃调色器），节点回复经网关汇总后通过 MQTT 上报服务器。整个系统围绕"玻璃颜色控制"这一核心业务展开。

```
服务器 (MQTT Broker)  <--->  网关 (iot-go-direct)  <--->  终端节点 (LoRa)
      MQTT json                    LoRa 二进制帧
```

| 参数 | 值 | 说明 |
|------|-----|------|
| 传输层 | TCP / MQTT 3.1.1 | 明文 TCP，默认端口 1883，未启用 TLS |
| QoS | QoS 0 | 订阅与发布均为 At-most-once |
| Retain | false | 所有发布消息 retain = false |
| ClientID | GatewayNodeId | 12 位十六进制字符串，如 `F12309150001` |
| 自动重连 | 开启 | 最大重连间隔 10 秒 |
| 消息编码 | JSON | UTF-8 文本 payload |

---

## 2. 连接配置

### MqttParams 结构体

```go
type MqttParams struct {
    MqttIP       string `json:"mqtt_ip"`
    MqttPort     int    `json:"mqtt_port"`
    MqttUserName string `json:"mqtt_user_name"`
    MqttPwd      string `json:"mqtt_pwd"`
}
```

### 默认配置值

| 参数 | 默认值 |
|------|--------|
| `mqtt_ip` | `app.kosglass.com`（生产）/ `120.79.55.61`（fallback） |
| `mqtt_port` | `1883` |
| `mqtt_user_name` | `l8juew73i2t17wavzthg` |
| `mqtt_pwd` | `i0eprmhypu3r16g3wuuc` |

### 连接生命周期关键点

- **broker 地址**：`tcp://{MqttIP}:{MqttPort}`
- **断线回调**：关闭 `sys-led-net` 网络指示灯
- **连上回调**：点亮 `sys-led-net`，自动重新订阅主题
- **重连间隔**：最大 10 秒（`SetMaxReconnectInterval`）
- **首次连接**：最多重试 100 次，超过则 `log.Fatal`

> **注意**：连接参数可由服务器通过 `mqtt_config` 指令远程修改，修改后网关会自动重启生效。

---

## 3. 主题结构

主题采用 `device/{gatewayNodeId}/{direction}` 的两级命名约定：

| 主题 | 方向 | 说明 |
|------|------|------|
| `device/{gatewayNodeId}/in` | ↓ 下行 | 网关订阅此主题，接收服务器下发的请求 |
| `device/{gatewayNodeId}/out` | ↑ 上行 | 网关向此主题发布状态与回复 |

其中 `{gatewayNodeId}` 是 12 位十六进制字符串（如 `F12309150001`），来源于 `gateway_id.txt` 文件或配置文件。订阅时 QoS = 0，发布时 QoS = 0、retain = false。

> ⚠️ ClientID 也使用 `gatewayNodeId`，因此同一网关 ID 不应同时启动两个实例，否则会触发 broker 的互踢。

---

## 4. 消息格式

所有 MQTT payload 均为 JSON 文本。消息分为 **请求（服务器→网关）** 与 **回复/上报（网关→服务器）** 两类：

| 类型 | 标识字段 | 结构示例 |
|------|----------|----------|
| 请求（下行） | `method` | `{"method":"xxx", "params":{...}}` |
| 回复/上报（上行） | `msg_type` | `{"msg_type":"xxx_reply", ...}` |

### 基础请求结构

```go
// BaseRequest 仅用于路由分发
type BaseRequest struct {
    Method string `json:"method"`
}
```

### 基础回复结构

```go
// BaseReply 所有回复的公共字段
type BaseReply struct {
    MsgType string `json:"msg_type"`
}

// 多数回复会附带 GatewayNodeId 用于服务器识别来源
type GatewayNodeIdReply struct {
    MsgType       string `json:"msg_type"`
    GatewayNodeId string `json:"gateway_node_id"`
}
```

### JSON Schema 校验

对 `config` 与 `update_glass_color_request` 两类消息，网关在反序列化前会使用 `gojsonschema` 进行 schema 校验，校验失败则丢弃消息。Schema 文件位于 `pkg/msg/json_schema/`，通过 `//go:embed` 编译进二进制。

---

## 5. 收发流程

网关主循环 `TopLevelMsgLoop` 使用 Go 的 `select` 在多个 channel 间轮转，MQTT 收发通过带缓冲 channel（容量 10）解耦。

### 主循环的 5 个 select 分支

1. **MqttToServerCh**：将待发消息通过 `SendToServer` 发布到 `/out` 主题
2. **MqttFromServerCh**：调用 `HandleMqttMsg` 解析并分发下行请求
3. **NodeMsgCh**：处理来自 LoRa 模块的节点回复（`HandleNodeMsg`）
4. **TimeoutNodeIdCh**：节点心跳超时，标记为离线
5. **ticker.C**：每 10 秒触发一次，达到上报周期后发送 `heartbeat_status_update`

### 上报周期计算

每 10 秒一个 tick，`PeriodNumberForReportingToServer = HeartbeatToServer / 10`。默认 `HeartbeatToServer = 60` 秒，即每 6 个 tick 上报一次。最小允许值 30 秒。

---

## 6. 服务器 → 网关（下行指令）

所有下行指令通过 `method` 字段路由，在 `HandleMqttMsg` 的 `switch` 中分发。

### 6.1 下行指令总表

| method | 用途 | params 结构 | 回复 msg_type |
|--------|------|-------------|---------------|
| `config` | 下发节点列表与 LoRa 模块参数配置（触发节点初始化） | ConfigParams：node_list_1/2、module1/2、heart_beat 等 | `config_reply` |
| `mqtt_config` | 远程修改 MQTT 连接参数（改后重启） | MqttParams：ip/port/user/pwd | `mqtt_config_reply` |
| `update_glass_color_request` | **逐个节点设置玻璃颜色** | []UpdateGlassColorParams：node_id + color(hex) | `update_glass_color_reply` |
| `group_update_glass_color_request` | 按组批量设置颜色（一次下发多节点同色） | node_list1/2 + color | `group_update_glass_color_reply` |
| `broadcast_update_glass_color_request` | 广播设置所有节点同一颜色 | ColorParams：color | `broadcast_update_glass_color_reply` |
| `get_glass_status_request` | 查询单个节点当前玻璃颜色 | GlassStatusParams：node_id | `glass_status_update`（异步） |
| `node_info_request` | 查询网关自身信息（版本、模块、MQTT 参数） | GatewayParams：gateway_node_id | `node_info_reply` |
| `node_list_request` | 查询网关当前管理的所有节点列表 | 无 | `node_list_reply` |
| `gateway_upgrade_request` | 网关固件升级（通过 SFTP 下载） | GatewayUpgradeParams + SftpInfo | `gateway_upgrade_reply` |
| `node_firmware_download_request` | 节点固件下载（透传给节点） | NodeFirmwareDownloadParam + SftpInfo | `node_firmware_download_reply` |
| `node_firmware_upgrade_request` | 节点固件升级（透传给节点） | NodeFirmwareUpgradeParams | `node_firmware_download_reply` |
| `set_touch_device_node_list_request` | 设置触摸设备关联的节点列表 | []TouchDeviceNodeList | 无（处理函数仅打印） |
| `gateway_reboot` | 远程重启网关 | 无 | `gateway_reboot_reply` |
| `execute_cmd` | 远程执行 shell 命令 | RemoteCmdParams：cmd | `remote_cmd_result`（异步） |

### 6.2 config（初始化配置）

这是网关启动后最重要的指令，服务器下发节点列表与 LoRa 射频参数，网关收到后会对所有节点执行 LoRa 初始化流程。

```json
{
  "method": "config",
  "params": {
    "node_list_1": ["12345678", "abcdef01"],
    "node_list_2": ["11111111"],
    "touch_node_list_1": [],
    "touch_node_list_2": [],
    "heart_beat": 20,
    "heart_beat_to_server": 60,
    "module1": {"freq": 4734, "band": 250, "factor": 9},
    "module2": {"freq": 4463, "band": 250, "factor": 9}
  }
}
```

**Schema 约束（`init_msg_schema.json`）**：

- `node_id` 必须匹配 `[0-9a-fA-F]{8}`（8 位十六进制）
- `module1.freq` 范围 3700–6000
- `band` 枚举：125 / 250 / 500
- `factor` 枚举：6–12
- 必填：node_list_1/2、touch_node_list_1/2、heart_beat、module1/2

### 6.3 update_glass_color_request（单节点调色）⭐

**核心调色指令**，支持一次请求中携带多个节点的颜色参数。

#### 请求结构体（Go）

```go
// 请求体
type UpdateGlassColorRequest struct {
    Method string                          `json:"method"`
    Params []shared.UpdateGlassColorParams `json:"params"`
}

// 单个节点颜色参数
type UpdateGlassColorParams struct {
    NodeId string `json:"node_id"`  // 8 位十六进制节点 ID
    Color  string `json:"color"`   // 颜色十六进制字符串
}
```

#### 请求示例

```json
{
  "method": "update_glass_color_request",
  "params": [
    {"node_id": "12345678", "color": "1122334455667788"},
    {"node_id": "abcdef01", "color": "11ff22ff33ff44ff"}
  ]
}
```

#### color 字段编码规则

`color` 为十六进制字符串，**每 2 个字符表示一个玻璃区域**：

- **第 1 位**：区域编号（`1`–`8`），对应玻璃的 8 个分区
- **第 2 位**：该区域的颜色值（`0`–`F`），共 16 种颜色

例如 `1122334455667788` 表示：

| 字符对 | 区域 | 颜色值 |
|--------|------|--------|
| `11` | 区域 1 | 颜色 1 |
| `22` | 区域 2 | 颜色 2 |
| `33` | 区域 3 | 颜色 3 |
| `44` | 区域 4 | 颜色 4 |
| `55` | 区域 5 | 颜色 5 |
| `66` | 区域 6 | 颜色 6 |
| `77` | 区域 7 | 颜色 7 |
| `88` | 区域 8 | 颜色 8 |

> 颜色字符串总长度必须为 **偶数**，由 `hex.DecodeString` 验证，非法值会导致该节点被标记为无效颜色。

#### JSON Schema 约束（`update_color_schema.json`）

```json
{
  "type": "object",
  "properties": {
    "method": {"type": "string"},
    "params": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "node_id": {
            "type": "string",
            "pattern": "[0-9a-fA-F]{8}"
          },
          "color": {
            "type": "string",
            "pattern": "([0-9a-fA-F][0-9a-fA-F])+"
          }
        },
        "required": ["node_id", "color"]
      }
    }
  },
  "required": ["method", "params"]
}
```

#### 网关处理流程

1. **Schema 校验**：先通过 `gojsonschema` 校验消息格式，不通过则丢弃
2. **逐节点处理**：遍历 `Params` 数组中的每个节点
3. **颜色验证**：对每个 `Color` 执行 `hex.DecodeString`，验证是否为合法十六进制
4. **设置请求颜色**：调用 `bsp.SetRequestingColor(nodeId, color)` 记录该节点的目标颜色到状态机
5. **模块分发**：
   - 若节点在 `node_list_1` 中 → 通过 **Module1** 发送 LoRa 指令
   - 若节点在 `node_list_2` 中 → 通过 **Module2** 发送 LoRa 指令
   - 若节点不在任何列表中 → 将该节点颜色标记为无效（`SetColorForNodeAsInvalid`），直接加入回复
6. **并行发送**：Module1 和 Module2 的 LoRa 发送在两个 goroutine 中并行执行，通过 `sync.WaitGroup` 等待完成
7. **组装回复**：将所有处理结果（含无效节点）组装为 `update_glass_color_reply`
8. **上报回复**：将回复放入 `MqttToServerCh`，由主循环发布到 `/out` 主题

#### 回复结构体（Go）

```go
// 调色回复
type UpdateGlassColorReply struct {
    MsgType       string                   `json:"msg_type"`
    GatewayNodeId string                   `json:"gateway_node_id"`
    Status        []UpdateGlassColorParams `json:"status"`
}
```

#### 回复示例

```json
{
  "msg_type": "update_glass_color_reply",
  "gateway_node_id": "F12309150001",
  "status": [
    {"node_id": "12345678", "color": "1122334455667788"},
    {"node_id": "abcdef01", "color": "11ff22ff33ff44ff"}
  ]
}
```

> **注意**：回复中的 `color` 是请求中传入的原始值。网关**不会**在此时等待节点实际执行完毕，实际执行结果会在后续 `heartbeat_status_update` 中通过 `completion_status` 和 `node_reported_color` 反映。

### 6.4 group_update_glass_color_request（分组调色）

与 `update_glass_color_request` 类似，但按组下发，节点列表使用网关已知的 `node_list_1` 和 `node_list_2`：

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

**回复**：`group_update_glass_color_reply`，含 `invalid_nodes[]`（不在任何列表中的节点）。

**重发机制**：该请求会被缓存到 `KeptGroupUpdateGlassColorRequest`，若心跳上报发现 `NodeReportedColor != NodeRequestingColor`，则自动重发最多 5 次。

### 6.5 broadcast_update_glass_color_request（广播调色）

向所有节点广播同一颜色：

```json
{
  "method": "broadcast_update_glass_color_request",
  "params": {
    "color": "11ff22ff33ff44ff"
  }
}
```

广播帧不指定节点 ID，通过 LoRa 公共频率（Module1 和 Module2 同时发送）下发。

### 6.6 gateway_upgrade_request（网关升级）

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

**升级条件校验**：

- `gateway_node_id` 必须匹配
- `target_hardware_version` 必须匹配
- `target_software_version` 必须 **大于** 本地版本

下载成功后设置 `IsRebootNeeded`，在下一次消息发送后执行 `reboot`。

---

## 7. 网关 → 服务器（上行消息）

上行消息由三类触发：①下行指令的同步回复；②LoRa 节点回复触发的异步上报；③周期性心跳上报。

### 7.1 上行消息总表

| msg_type | 触发方式 | 关键字段 |
|----------|----------|----------|
| `init` | 启动后自动发送（连接成功即发） | gateway_node_id、hard/soft_version、module0/1/2、heart_beat 等 |
| `heartbeat_status_update` | 周期性（默认 60s） | gateway_node_id、status[]：每节点颜色/RSSI/SNR/离线状态/完成状态 |
| `config_reply` | config 指令回复 | gateway_node_id |
| `mqtt_config_reply` | mqtt_config 指令回复 | msg_type |
| `update_glass_color_reply` | update_glass_color_request 回复 | gateway_node_id、status[]：node_id + color |
| `group_update_glass_color_reply` | group_update 回复 | gateway_node_id、invalid_nodes[] |
| `broadcast_update_glass_color_reply` | broadcast 回复 | gateway_node_id |
| `glass_status_update` | get_glass_status 异步回复 | gateway_node_id、status[]：node_id + color |
| `node_info_reply` | node_info_request 回复 | 网关全部信息 + mqtt_params |
| `node_list_reply` | node_list_request 回复 | node_list_1/2、touch_node_list_1/2、total_count |
| `gateway_upgrade_reply` | gateway_upgrade 回复 | 原参数 + state(OK/error) |
| `node_firmware_download_reply` | 节点固件下载/升级回复 | 原参数 + state |
| `gateway_reboot_reply` | gateway_reboot 回复 | gateway_node_id |
| `remote_cmd_result` | execute_cmd 异步回复 | gateway_node_id、stdout、stderr |
| `unknown_node` | 收到未知节点回复时主动上报 | unknown_node_list[] |

### 7.2 init（启动上报）

网关 MQTT 连接建立并订阅成功后，立即发送 `init` 消息：

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

服务器据此识别网关上线并获取其配置信息。

### 7.3 heartbeat_status_update（周期心跳）

核心上报消息，包含网关下所有节点的实时状态。每 10 秒检查一次，达到 `HeartbeatToServer` 周期后发送。

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
| `node_id` | string | 8 位十六进制节点 ID |
| `color` | string | 节点**实际报告**的玻璃颜色（8 区域×2字符） |
| `node_requesting_color` | string | 服务器**最近一次请求**的颜色 |
| `hard_version` | string | 节点硬件版本（十六进制） |
| `soft_version` | string | 节点软件版本（十六进制） |
| `run_area` | int | 节点运行区域 |
| `rssi` | float64 | 信号强度（dBm） |
| `snr` | float64 | 信噪比 |
| `is_offline` | bool | 是否离线（心跳超时） |
| `completion_status` | int | 命令完成状态：`0`=错误，`1`=进行中，`2`=已完成 |

**color 编码**：由 `GetColorStrFromSlice` 生成，8 个区域，每区域 2 个字符，第 1 位为区域序号（1-8），第 2 位为颜色值（0-F）。离线节点的颜色会被 `SetColorForNodeAsInvalid` 处理（每隔一位插入 f）。

### 7.4 unknown_node（未知节点上报）

当网关收到未在 `node_list` 中的节点的 LoRa 回复（`UNKNOWN_NODE_REPLY`）时，主动上报该节点 ID：

```json
{
  "msg_type": "unknown_node",
  "unknown_node_list": ["99999999"]
}
```

---

## 8. 与 LoRa 协议的桥接

MQTT 消息本身不直接到达节点，需经网关转译为 LoRa 二进制帧。

### 8.1 LoRa 帧结构

```
| len(1B) | boardType(1B) | cmdType(1B) | gatewayId(6B) | nodeId(4B) | payload... | CRC8(1B) |
```

所有帧以 **CRC8** 校验结尾（查表法）。网关收到节点回复时先校验 CRC 与 gatewayId，不匹配则丢弃。

### 8.2 LoRa 命令类型与 MQTT 消息映射

| cmdType | 名称 | 方向 | 对应 MQTT 消息 |
|---------|------|------|----------------|
| 1 | HEARTBEAT_REQ | 网关→节点 | 由 config 触发的周期心跳 |
| 2 | HEARTBEAT_REPLY | 节点→网关 | 上报为 heartbeat_status_update |
| 3 | CONFIG_NODE_REQ | 网关→节点 | 由 config 触发的节点初始化 |
| 4 | CONFIG_NODE_REPLY | 节点→网关 | 更新节点状态（版本/模块参数） |
| 5 | UPDATE_GLASS_COLOR_REQ | 网关→节点 | update_glass_color_request |
| 6 | UPDATE_GLASS_COLOR_REPLY | 节点→网关 | 更新节点颜色状态 |
| 7 | GROUP_UPDATE_COLOR_REQ | 网关→节点 | group_update_glass_color_request |
| 8 | BROADCAST_UPDATE_COLOR_REQ | 网关→节点 | broadcast_update_glass_color_request |
| 9 | GET_GLASS_STATE_REQ | 网关→节点 | get_glass_status_request |
| 10 | GET_GLASS_STATE_REPLY | 节点→网关 | glass_status_update |
| 40 | UNKNOWN_NODE_REPLY | 节点→网关 | unknown_node |

### 8.3 调色对应的 LoRa 帧（UPDATE_GLASS_COLOR_REQ，cmdType=5）

```go
// 帧内容：
// byte[0]  = 包长度（含 CRC）
// byte[1]  = boardType = 1（网关）
// byte[2]  = cmdType = 5
// byte[3:9] = gatewayId（6 字节）
// byte[9:13] = nodeId（4 字节）
// byte[13]  = color 参数长度（字节数 = len(color)/2）
// byte[14:] = color 数据（hex 解码后的字节）
// last byte = CRC8
```

例如 `color = "11223344"`（4 个区域），hex 解码后为 4 字节，参数长度为 4。

---

## 9. 完整交互示例

### 示例 1：网关启动上线

```
// 1. 网关 → 服务器：init
topic: device/F12309150001/out
payload: {"msg_type":"init","gateway_node_id":"F12309150001","soft_version":"2.2",...}

// 2. 服务器 → 网关：config
topic: device/F12309150001/in
payload: {"method":"config","params":{"node_list_1":["12345678"],...}}

// 3. 网关 → 服务器：config_reply
topic: device/F12309150001/out
payload: {"msg_type":"config_reply","gateway_node_id":"F12309150001"}
```

### 示例 2：单节点调色（update_glass_color_request）

```
// 1. 服务器 → 网关
topic: device/F12309150001/in
payload: {
  "method": "update_glass_color_request",
  "params": [{"node_id":"12345678","color":"11ff223344556677"}]
}

// 2. 网关 → 节点（LoRa 二进制帧）
//   cmdType=5, gatewayId=F12309150001, nodeId=12345678
//   colorLen=8, colorData=[0x11, 0xff, 0x22, 0x33, 0x44, 0x55, 0x66, 0x77]

// 3. 节点 → 网关（LoRa 回复）
//   cmdType=6, result=1(成功)

// 4. 网关 → 服务器：立即回复
topic: device/F12309150001/out
payload: {
  "msg_type":"update_glass_color_reply",
  "gateway_node_id":"F12309150001",
  "status":[{"node_id":"12345678","color":"11ff223344556677"}]
}

// 5. 网关 → 服务器：下次心跳上报实际颜色
topic: device/F12309150001/out
payload: {
  "msg_type":"heartbeat_status_update",
  "gateway_node_id":"F12309150001",
  "status":[{
    "node_id":"12345678",
    "color":"11ff223344556677",
    "completion_status":2,
    "rssi":-45.5,
    "snr":8.2,
    "is_offline":false
  }]
}
```

### 示例 3：分组调色

```
// 服务器 → 网关
topic: device/F12309150001/in
payload: {
  "method": "group_update_glass_color_request",
  "params": {
    "node_list1": ["12345678", "abcdef01"],
    "node_list2": ["11111111"],
    "color": "11ff22ff33ff44ff"
  }
}

// 网关 → 服务器
topic: device/F12309150001/out
payload: {
  "msg_type":"group_update_glass_color_reply",
  "gateway_node_id":"F12309150001",
  "invalid_nodes":[]
}
```

### 示例 4：远程命令执行

```
// 服务器 → 网关
topic: device/F12309150001/in
payload: {"method":"execute_cmd","params":{"cmd":"ls -l /tmp"}}

// 网关 → 服务器（异步）
topic: device/F12309150001/out
payload: {
  "msg_type":"remote_cmd_result",
  "gateway_node_id":"F12309150001",
  "stdout":"total 8\ndrwxr-xr-x 2 root root 4096 ...",
  "stderr":""
}
```

---

## 10. 关键常量

| 常量 | 值 | 说明 |
|------|-----|------|
| `APP_VERSION` | 2.2 | 网关软件版本 |
| `TO_SERVER_HEARTBEAT_SECONDS` | 60 | 默认服务器心跳周期（秒） |
| `HEARTBEAT_REPLY_TIMEOUT` | 10 | 节点心跳回复超时（秒） |
| `HEARTBEAT_RETRY_CNT` | 3 | 节点心跳重试次数 |
| `NODE_INIT_RETRY_CNT` | 3 | 节点初始化重试次数 |
| `UPDATE_GLASS_COLOR_RETRY_CNT` | 2 | 调色重试次数 |
| `GET_GLASS_STATUS_MAX_RETRY` | 3 | 查询玻璃状态重试次数 |
| `NODE_INIT_REPLY_TIMEOUT_SECONDS` | 6 | 节点初始化回复超时（秒） |
| `GATEWAY_ID_LEN` | 6 | LoRa 帧中 gatewayId 字节数 |
| `NODE_ID_LEN` | 4 | LoRa 帧中 nodeId 字节数 |
| 默认 Module0 | freq=4723, band=250, factor=9 | 公共频率（用于未入网节点） |
| 默认 Module1 | freq=4734, band=250, factor=9 | 节点列表 1 的射频参数 |
| 默认 Module2 | freq=4463, band=250, factor=9 | 节点列表 2 的射频参数 |

---

## ⚠️ 安全提示

1. **明文传输**：该实现使用明文 MQTT（无 TLS），用户名密码硬编码在 `bsp_config.go` 中作为默认值。生产环境建议启用 TLS。
2. **命令注入风险**：`execute_cmd` 指令可远程执行任意 shell 命令，建议对 `cmd` 做白名单校验。
3. **ClientID 冲突**：同一 `gatewayNodeId` 不应同时启动两个实例，否则会触发 broker 互踢。
