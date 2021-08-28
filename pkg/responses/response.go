package responses

type Response struct {
	Code uint32      `json:"code"` // 状态码 0 成功 其它失败
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"` // 消息
}

func New(code uint32, data interface{}, msg string) Response {
	return Response{
		code,
		data,
		msg,
	}
}
