package RTSPService

import (
	"encoding/json"
	"errors"
	"net"
	"strconv"

	"github.com/use-go/websocketStreamServer/logger"
	"github.com/use-go/websocketStreamServer/wssAPI"
)

type RTSPService struct {
}

type RTSPConfig struct {
	Port       int `json:"port"`
	TimeoutSec int `json:"timeoutSec"`
}

var service *RTSPService
var serviceConfig RTSPConfig

func (rtspService *RTSPService) Init(msg *wssAPI.Msg) (err error) {
	if nil == msg || nil == msg.Param1 {
		logger.LOGE("init rtsp server failed")
		return errors.New("invalid param init rtsp server")
	}
	fileName, ok := msg.Param1.(string)
	if false == ok {
		logger.LOGE("bad param init rtsp server")
		return errors.New("invalid param init rtsp server")
	}
	err = rtspService.loadConfigFile(fileName)
	if err != nil {
		logger.LOGE("load rtsp config failed:" + err.Error())
		return
	}
	return
}

func (rtspService *RTSPService) loadConfigFile(fileName string) (err error) {
	data, err := wssAPI.ReadFileAll(fileName)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &serviceConfig)
	if err != nil {
		return
	}
	if serviceConfig.TimeoutSec <= 0 {
		serviceConfig.TimeoutSec = 60
	}

	if serviceConfig.Port == 0 {
		serviceConfig.Port = 554
	}
	return
}

func (rtspService *RTSPService) Start(msg *wssAPI.Msg) (err error) {
	logger.LOGT("start RTSP server")
	strPort := ":" + strconv.Itoa(serviceConfig.Port)
	tcp, err := net.ResolveTCPAddr("tcp4", strPort)
	if err != nil {
		logger.LOGE(err.Error())
		return
	}
	listener, err := net.ListenTCP("tcp4", tcp)
	if err != nil {
		logger.LOGE(err.Error())
		return
	}
	go rtspService.rtspLoop(listener)
	return
}

func (rtspService *RTSPService) rtspLoop(listener *net.TCPListener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			logger.LOGE(err.Error())
			continue
		}
		go rtspService.handleConnect(conn)
	}
}

func (rtspService *RTSPService) handleConnect(conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	handler := &RTSPHandler{}
	handler.conn = conn
	handler.Init(nil)
	for {
		data, err := ReadPacket(conn, handler.tcpTimeout)
		if err != nil {
			logger.LOGE("read rtsp failed")
			logger.LOGE(err.Error())
			handler.handlePacket(nil)
			return
		}
		logger.LOGT(string(data))
		err = handler.handlePacket(data)
		if err != nil {
			logger.LOGE(err.Error())
			return
		}
	}
}

func (rtspService *RTSPService) Stop(msg *wssAPI.Msg) (err error) {
	return
}

func (rtspService *RTSPService) GetType() string {
	return wssAPI.OBJ_RTSPServer
}

func (rtspService *RTSPService) HandleTask(task wssAPI.Task) (err error) {
	return
}

func (rtspService *RTSPService) ProcessMessage(msg *wssAPI.Msg) (err error) {
	return
}
