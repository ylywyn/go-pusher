// automatically generated, do not modify

package types

///状态消息
const (
	StatusMsgTypeMT_STATUS_NONE = 400
	StatusMsgTypeMT_STATUS_HEARTBEAT = 401
	StatusMsgTypeMT_STATUS_ACK_FROMSERVER = 402
	StatusMsgTypeMT_STATUS_ACK_FROMCLIENT = 403
	StatusMsgTypeMT_STATUS_SESSIONINVALIDATE = 404
	///MT_STATUS_CRITICAL_ERROR以后的为严重错误，客户端接到消息后必须关闭连接
	StatusMsgTypeMT_STATUS_CRITICAL_ERROR = 450
	///消息头接收错误
	StatusMsgTypeMT_STATUS_HEADERREAD_ERROR = 451
	///消息头解析错误
	StatusMsgTypeMT_STATUS_HEADERPARSE_ERROR = 452
	///消息体接收错误
	StatusMsgTypeMT_STATUS_BODYREAD_ERROR = 453
	///消息体解析错误
	StatusMsgTypeMT_STATUS_BODYPARSE_ERROR = 454
	StatusMsgTypeMT_STATUS_END = 499
)
