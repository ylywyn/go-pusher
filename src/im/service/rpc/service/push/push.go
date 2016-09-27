package push

import (
	"errors"
	log "im/common/log4go"
	"im/common/proto/entity/msg"
	mb "im/common/proto/entity/msg/msgbase"
	"im/common/proto/fbsgen/msg/types"
	"im/service/logic/service/id"
	. "im/service/logic/service/send"
	gen "im/service/rpc/gen"
	"strconv"

	context "golang.org/x/net/context"
)

type PushService struct{}

func (server *PushService) Push(ctx context.Context, msg *gen.PusherMessage) (*gen.PusherReply, error) {
	if msg == nil {
		return nil, errors.New("param error")
	}

	l := len(msg.Text)
	if l < 1 || l > 2048 {
		return nil, errors.New("param error: text invalid")
	}

	id := id.GetId("rpcmsg")
	if id == 0 {
		return nil, errors.New("get msg id error")
	}

	//获取结构
	m := &mb.MsgBase{
		Type:     uint16(msg.Type),
		Appid:    uint16(msg.AppId),
		From:     msg.From,
		To:       msg.To,
		Gid:      msg.Gid,
		Text:     msg.Text,
		Time:     msg.Time,
		Platform: uint8(msg.Platform),
		Msgid:    id,
	}

	//发送推送消息
	go sendAndSave(m)

	r := &gen.PusherReply{
		Result: true,
		Data:   []string{strconv.FormatUint(id, 10)},
	}
	return r, nil

}

func sendAndSave(m *mb.MsgBase) {

	var err error
	h := &msg.MsgHeader{
		Ack:      0,
		Compress: 0,
		Repeat:   0,
		Type:     m.Type,
	}

	switch m.Type {
	case types.GeneralMsgTypeMT_GENERAL_MSG,
		types.GeneralMsgTypeMT_GENERAL_NOTICE:
		err = Send.SendOne(h, m)
	case types.GeneralMsgTypeMT_GENERAL_GROUP_MSG,
		types.GeneralMsgTypeMT_GENERAL_GROUP_NOTICE:
		err = Send.SendToGroup(h, m)
	default:
		err = errors.New("不支持的消息类型")
	}

	if err != nil {
		log.Error("[http|Push|sendAndSave] error:%s", err.Error())
	}
}
