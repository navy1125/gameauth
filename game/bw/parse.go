package bw

import (
	"github.com/navy1125/gotcp/bw/common"
	"github.com/navy1125/gotcp/gotcp"
)

var (
	handleMessageMap gotcp.HanldeMessageMap
)

func init() {
	RegisterMessage(Cmd.TIME_USERCMD, Cmd.GAMETIME_TIMER_USERCMD_PARA, parseStGameTimeTimerUserCmd)
	RegisterMessage(Cmd.TIME_USERCMD, Cmd.REQUESTUSERGAMETIME_TIMER_USERCMD_PARA, parseStRequestUserGameTimeTimerUserCmd)
	RegisterMessage(Cmd.TIME_USERCMD, Cmd.USERGAMETIME_TIMER_USERCMD_PARA, parseStUserGameTimeTimerUserCmd)

}

func RegisterMessage(byCmd, byParam byte, fun gotcp.HandleMessageFunc) {
	handleMessageMap[byCmd][byParam] = fun
}

func parseStGameTimeTimerUserCmd(task *gotcp.Task, data []byte) {
	task.Debug("heartBeat")
	cmd := Cmd.NewStGameTimeTimerUserCmd()
	task.SendCmd(cmd)
}

func parseStRequestUserGameTimeTimerUserCmd(task *gotcp.Task, data []byte) {
	cmd := Cmd.NewStRequestUserGameTimeTimerUserCmd()
	err := task.GetCmd(data, cmd)
	if err != nil {
		task.Error(err.Error())
		return
	}
	//task.Debug("parseStRequestUserGameTimeTimerUserCmd:%d,%d", cmd.ByCmd, cmd.ByParam)
	task.SendCmd(cmd)
}

func parseStUserGameTimeTimerUserCmd(task *gotcp.Task, data []byte) {
	cmd := Cmd.NewStUserGameTimeTimerUserCmd()
	err := task.GetCmd(data, cmd)
	if err != nil {
		task.Error(err.Error())
		return
	}
	cmd.QwGameTime = 0
	cmd.Mac = 0
	//task.Debug("parseStUserGameTimeTimerUserCmd:%d,%d", cmd.ByCmd, cmd.ByParam)
}
