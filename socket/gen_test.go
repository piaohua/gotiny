package socket

import (
	"fmt"
	"gotiny/protocol"
	"testing"

	"github.com/gogo/protobuf/proto"
)

// test
func TestGen(t *testing.T) {
	Init()
	Gen()
}

// test
func TestPacket(t *testing.T) {
	msg := new(protocol.SRegist)
	ab(msg)
	c, b, err := packet(msg)
	t.Log(c, b, err)
	msg2, err := unpack(c, b)
	t.Log(msg2, err)
}

func ab(msg interface{}) {
	switch msg.(type) {
	case *protocol.SRegist:
		fmt.Println("msg 2")
	case proto.Message:
		fmt.Println("msg")
	default:
		fmt.Println("default")
	}
}
