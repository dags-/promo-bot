package server

import (
	"fmt"

	"github.com/qiangxue/fasthttp-routing"
)

const addUrl = "https://discordapp.com/api/oauth2/authorize?client_id=%s&scope=bot&permissions=0x00002000"

func (s *Server) handleAdd(context *routing.Context) error {
	context.Redirect(fmt.Sprintf(addUrl, s.clientId), redirectCode)
	return nil
}
