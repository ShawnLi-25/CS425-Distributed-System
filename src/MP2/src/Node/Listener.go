package node

// Listener is a type that implements the ListenMsg(), ListenJoinMsg() "method"
type Listener struct{}

//ListenMsg: Listen to Heartbeat or Leave Msg
func (l *Listener) ListenMsg() {

}

//ListenJoinMsg: Listen to Join Msg (Introducer-only)
func (l *Listener) ListenJoinMsg() {

}
