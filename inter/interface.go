package inter

type Pid interface {
	Close()
	Send(interface{})
	Call(interface{}) interface{}
}
