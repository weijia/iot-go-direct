package main

// mqtt_pub_test 是一个零外部依赖(复用 paho)的 MQTT 测试发布器，用于配合
// run_emqx_test.bat 验证端到端闭环：
//   发 MQTT 指令 -> 网关经(mock)串口发给虚拟玻璃 -> MockBoard 回包 -> 状态回填 -> 心跳/回复回报
//
// 默认发 get_glass_status_request（无 schema 校验，node1 在默认节点列表），
// 可触发网关经串口发 cmd=2 并由 MockBoard 回包。
// 改色 update_glass_color_request 因 schema 要求 node_id 为 8 位十六进制，
// 默认列表里的 "node1" 会被拒；需传入符合格式的 node_id 才会触发串口。

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	broker := flag.String("broker", "broker.emqx.io:1883", "MQTT broker host:port")
	gateway := flag.String("gateway", "F12309150001", "gateway node id")
	method := flag.String("method", "get_glass_status_request", "mqtt method to send")
	node := flag.String("node", "node1", "node id")
	color := flag.String("color", "1022344052647082", "color hex (for update)")
	wait := flag.Int("wait", 8, "seconds to wait for reply")
	flag.Parse()

	server := "tcp://" + *broker
	opts := mqtt.NewClientOptions().AddBroker(server)
	opts.SetClientID(fmt.Sprintf("test-pub-%d", time.Now().UnixNano()))
	opts.SetAutoReconnect(true)
	cli := mqtt.NewClient(opts)
	if token := cli.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("connect error:", token.Error())
		os.Exit(1)
	}
	defer cli.Disconnect(1000)

	inTopic := "device/" + *gateway + "/in"
	outTopic := "device/" + *gateway + "/out"

	cli.Subscribe(outTopic, 0, func(c mqtt.Client, m mqtt.Message) {
		fmt.Printf("[REPLY] %s\n%s\n", m.Topic(), string(m.Payload()))
	})

	var payload string
	switch *method {
	case "get_glass_status_request":
		payload = fmt.Sprintf(`{"method":"get_glass_status_request","params":{"node_id":%q}}`, *node)
	case "update_glass_color_request":
		payload = fmt.Sprintf(`{"method":"update_glass_color_request","params":[{"node_id":%q,"color":%q}]}`, *node, *color)
	default:
		payload = fmt.Sprintf(`{"method":%q,"params":{"node_id":%q}}`, *method, *node)
	}

	fmt.Printf("Publishing to %s:\n%s\n", inTopic, payload)
	cli.Publish(inTopic, 0, false, payload)
	fmt.Printf("Subscribed to %s, waiting %ds for reply (Ctrl+C to exit early)...\n", outTopic, *wait)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	select {
	case <-sig:
	case <-time.After(time.Duration(*wait) * time.Second):
	}
	fmt.Println("done")
}
