package Cmd

import (
	"github.com/navy1125/gotcp/bw/common"
)

const BILL_USERCMD = 81

type StBillUserCmd struct {
	Cmd.StNullUserCmd
}

type StZoneBillUserCmd struct {
	StBillUserCmd
	Zoneid uint32
}

/// WebGate请求登陆WebGateServer
const REQUEST_LOGIN_BILL_USERCMD_PARA_C = 1

type StRequestLoginBillUserCmd struct {
	StBillUserCmd
	Version    int32 /**< 比宏WebGateTOOL_VERSION小的将无法登陆 */
	Servertype int32 /**< 服务器类型///ServerType */
}

func NewStRequestLoginBillUserCmd() *StRequestLoginBillUserCmd {
	cmd := &StRequestLoginBillUserCmd{}
	cmd.ByCmd = BILL_USERCMD
	cmd.ByParam = REQUEST_LOGIN_WEBGATE_USERCMD_PARA_C
	return cmd
}

/// WebGate登陆返回
const RETURN_LOGIN_BILL_USERCMD_PARA_S = 2

type StReturnLoginBillUserCmd struct {
	StBillUserCmd
}

func NewStReturnLoginBillUserCmd() *StReturnLoginBillUserCmd {
	cmd := &StReturnLoginBillUserCmd{}
	cmd.ByCmd = WEBGATE_USERCMD
	cmd.ByParam = RETURN_LOGIN_WEBGATE_USERCMD_PARA_S
	return cmd
}

const RETURN_PLATFORM_TO_GOLD_BILL_USERCMD_PARA_SC = 26

type StReturnPlatformToGoldZoneBillUserCmd struct {
	StZoneBillUserCmd
	DwNum  uint32 //本次兑换点数
	Charid uint32
	Accid  uint32
	tid    [Cmd.MAX_MACSIZE]byte //流水号
}

func NewStReturnPlatformToGoldZoneBillUserCmd() *StReturnPlatformToGoldZoneBillUserCmd {
	cmd := &StReturnPlatformToGoldZoneBillUserCmd{}
	cmd.ByCmd = WEBGATE_USERCMD
	cmd.ByParam = RETURN_PLATFORM_TO_GOLD_BILL_USERCMD_PARA_SC
	return cmd
}
