package actors

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/asynkron/protoactor-go/actor"
)

type VoteWriter struct{}

func (state *VoteWriter) Receive(ctx actor.Context) {
	switch msg := ctx.Message().(type) {
	case *VoteRequest:
		//在当前工作目录下创建users.txt文件, 写入票号和username
		file, err := os.OpenFile("users.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()
		if _, err := file.WriteString(fmt.Sprintf("%s\n", strings.Join(msg.Usernames, ","))); err != nil {
			log.Fatal(err)
		}

	}
}
