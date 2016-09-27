/*******************************************************************
 *  Copyright(c) 2015-2016 Company Name
 *  All rights reserved.
 *
 *  File   :
 *  Date   :
 *  Author : yangl
 *  Description:
 ******************************************************************/

package msg

//苹果消息的内部交换格式
type ApnsTransMsg struct {
	Dev        bool     `json:"dev"`
	Appid      uint16   `json:"appid"`
	ExpireTime uint32   `json:"expiretime"`
	Payload    string   `json:"payload"` // ApsPayloadTrans的json字符串，为了反序列化时不必要的过多解析，这里直接定义成了string
	Tokens     []string `json:"tokens"`
}

// payload
type ApsPayloadTrans struct {
	Aps    ApsSimple `json:"aps"`
	Extras ApsExtras `json:"extra"`
}

// aps
type ApsSimple struct {
	Alert             string `json:"alert"` //去掉更复杂的alert信息，仅定义成string(为了和其他平台统一), IOS特化的信息自行构造Alert的Json格式
	Badge             uint16 `json:"badge"`
	Sound             string `json:"sound"`
	Content_available bool   `json:"content-available"`
	Category          string `json:"category"`
}

// 消息的额外信息：消息类型， 发送者，时间等
type ApsExtras struct {
	Type   uint16
	From   uint64
	To     uint64
	Gid    uint64
	Time   uint64
	Msgid  uint64
	Extras string
}
