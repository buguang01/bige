package server

import (
	"os"
	"os/signal"
)

//游戏的逻辑服务器
type GameServer struct {
}

func (gs *GameServer) Run() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
}
