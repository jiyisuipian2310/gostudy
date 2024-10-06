package tcpserver

type result struct {
	code int
	msg  string
	err  error
}

func (r *result) Response(handler *clientHandler) error {
	if r.code != 0 {
		return handler.writeMessage(r.code, r.msg)
	}
	return nil
}
