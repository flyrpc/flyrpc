package fly

import (
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	server := NewServer(&ServerOpts{
		Serializer: Protobuf,
	})
	err := server.Listen("127.0.0.1:15555")
	assert.Nil(t, err)
	assert.NotNil(t, server.listener)

	server.OnMessage(1, func(ctx *Context, in *TestUser) *TestUser {
		reply := new(TestUser)
		log.Println("server on", in)
		// call client cmd 2 and response another user
		log.Println("server call 2")
		err := ctx.Call(2, reply, in)
		assert.Nil(t, err)
		log.Println("client response", reply)
		reply.Id += 10
		return reply
	})

	log.Println("test l", server.listener)
	go server.HandleConnections()

	conn, err := net.Dial("tcp", "127.0.0.1:15555")
	assert.Nil(t, err)
	log.Println("connected")
	client := NewClient(conn, Protobuf)
	client.OnMessage(2, func(ctx *Context, in *TestUser) *TestUser {
		log.Println("client on", in)
		return &TestUser{Id: in.Id + 1}
	})
	reply := new(TestUser)
	log.Println("client call 1")
	err = client.Call(1, reply, &TestUser{Id: 100})
	assert.Nil(t, err)
	log.Println("server response", reply)

	assert.Equal(t, 111, reply.Id)
}
