package Cmd

import (
	"github.com/navy1125/gotcp/bw/common"
)

const WEBGATE_USERCMD = 94

type StWebGateUserCmd struct {
	Cmd.StNullUserCmd
}

type StRequestLoginLoginCmd struct {
	StWebGateUserCmd
}

func NewStRequestLoginLoginCmd() *StRequestLoginLoginCmd {
	cmd := &StRequestLoginLoginCmd{}
	cmd.ByCmd = 1
	cmd.ByParam = 1
	return cmd
}

/// WebGate请求登陆WebGateServer
const REQUEST_LOGIN_WEBGATE_USERCMD_PARA_C = 1

type StRequestLoginWebGateUserCmd struct {
	StWebGateUserCmd
	Version    int32 /**< 比宏WebGateTOOL_VERSION小的将无法登陆 */
	Servertype int32 /**< 服务器类型///ServerType */
}

func NewStRequestLoginWebGateUserCmd() *StRequestLoginWebGateUserCmd {
	cmd := &StRequestLoginWebGateUserCmd{}
	cmd.ByCmd = WEBGATE_USERCMD
	cmd.ByParam = REQUEST_LOGIN_WEBGATE_USERCMD_PARA_C
	return cmd
}

/// WebGate登陆返回
const RETURN_LOGIN_WEBGATE_USERCMD_PARA_S = 2

type StReturnLoginWebGateUserCmd struct {
	StWebGateUserCmd
}

func NewStReturnLoginWebGateUserCmd() *StReturnLoginWebGateUserCmd {
	cmd := &StReturnLoginWebGateUserCmd{}
	cmd.ByCmd = WEBGATE_USERCMD
	cmd.ByParam = RETURN_LOGIN_WEBGATE_USERCMD_PARA_S
	return cmd
}

type stZoneWebGateUserCmd struct {
	StWebGateUserCmd
	Zoneid   uint32
	Accid    uint32
	Reserved uint32
}

/// WEBGATE客户端登陆登陆服务器
const WEB_LOGIN_USER_TOKEN_WEBGATE_USERCMD_PARA_S = 33

type StWebLoginUserTokenWebGateUserCmd struct {
	stZoneWebGateUserCmd
	Account   [Cmd.MAX_ACCNAMESIZE]byte
	Token     [Cmd.MAX_TOKENSIZE]byte
	Lifetime  uint32 //token过期时间,0表示只登陆一次
	UserType  uint32 ///ChannelType
	LoginType uint32 ///LoginType
	Password  [Cmd.MAX_PASSWORD]byte
	Type      uint32 ///登陆类型，1表示成年人
}

func NewStWebLoginUserTokenWebGateUserCmd() *StWebLoginUserTokenWebGateUserCmd {
	cmd := &StWebLoginUserTokenWebGateUserCmd{}
	cmd.ByCmd = WEBGATE_USERCMD
	cmd.ByParam = WEB_LOGIN_USER_TOKEN_WEBGATE_USERCMD_PARA_S
	return cmd
}
