package main

import (
	"github.com/fananchong/go-xserver/services/internal/protocol"
	"github.com/fananchong/go-xserver/services/internal/utility"
	"github.com/fananchong/gotcp"
)

// onChat : 聊天
func (accountObj *Account) onChat(data []byte) {
	Ctx.Log.Infoln("Chat, account:", accountObj.account)
	msg := &protocol.MSG_LOBBY_CHAT{}
	if gotcp.DecodeCmd(data[:len(data)-1], data[len(data)-1], msg) == nil {
		Ctx.Log.Errorln("Message parsing failed, message number is`protocol.MSG_LOBBY_CHAT`(", int(protocol.CMD_LOBBY_CHAT), "). account", accountObj.account)
		return
	}
	msg.From = accountObj.GetRole().GetName()
	if msg.GetTo() == "" {
		// 全服广播
		utility.BroadcastMsgToClient(Ctx, uint64(protocol.CMD_LOBBY_CHAT), msg)
	} else {
		// 私聊
		utility.SendMsgToClientByRoleName(Ctx, msg.GetTo(), uint64(protocol.CMD_LOBBY_CHAT), msg)
	}
}
