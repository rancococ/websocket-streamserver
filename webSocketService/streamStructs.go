package webSocketService

import (
	"encoding/json"

	"github.com/use-go/websocketStreamServer/logger"
	"github.com/use-go/websocketStreamServer/wssAPI"

	"github.com/gorilla/websocket"
)

// WS Status
const (
	WSStatusStatus = "status"
	WSStatusError  = "error"
)

//1byte type
const (
	WSPktAudio   = 8
	WSPktVideo   = 9
	WSPktControl = 18
)

// WS Cmd
const (
	WSCPlay       = 1
	WSCPlay2      = 2
	WSCResume     = 3
	WSCPause      = 4
	WSCSeek       = 5
	WSCClose      = 7
	WSCStop       = 6
	WSCPublish    = 0x10
	WSCOnMetaData = 9
)

var cmdsMap map[int]*wssAPI.Set

func init() {
	cmdsMap = make(map[int]*wssAPI.Set)
	//初始状态close，可以play,close,publish
	{
		tmp := wssAPI.NewSet()
		tmp.Add(WSCPlay)
		tmp.Add(WSCPlay2)
		tmp.Add(WSCClose)
		tmp.Add(WSCPublish)
		cmdsMap[WSCClose] = tmp
	}
	//play 可以close pause seek
	{
		tmp := wssAPI.NewSet()
		tmp.Add(WSCPause)
		tmp.Add(WSCPlay)
		tmp.Add(WSCPlay2)
		tmp.Add(WSCSeek)
		tmp.Add(WSCClose)
		cmdsMap[WSCPlay] = tmp
	}
	//play2 =play
	{
		tmp := wssAPI.NewSet()
		tmp.Add(WSCPause)
		tmp.Add(WSCPlay)
		tmp.Add(WSCPlay2)
		tmp.Add(WSCSeek)
		tmp.Add(WSCClose)
		cmdsMap[WSCPlay2] = tmp
	}
	//pause
	{
		tmp := wssAPI.NewSet()
		tmp.Add(WSCResume)
		tmp.Add(WSCPlay)
		tmp.Add(WSCPlay2)
		tmp.Add(WSCClose)
		cmdsMap[WSCPause] = tmp
	}
	//publish
	{
		tmp := wssAPI.NewSet()
		tmp.Add(WSCClose)
		cmdsMap[WSCPublish] = tmp
	}
}

func supportNewCmd(cmdOld, cmdNew int) bool {
	_, exist := cmdsMap[cmdOld]
	if false == exist {
		return false
	}
	return cmdsMap[cmdOld].Has(cmdNew)
}

// SendWsControl Control command
func SendWsControl(conn *websocket.Conn, ctrlType int, data []byte) (err error) {
	dataSend := make([]byte, len(data)+4)
	dataSend[0] = WSPktControl
	dataSend[1] = byte((ctrlType >> 16) & 0xff)
	dataSend[2] = byte((ctrlType >> 8) & 0xff)
	dataSend[3] = byte((ctrlType >> 0) & 0xff)
	copy(dataSend[4:], data)
	return conn.WriteMessage(websocket.BinaryMessage, dataSend)
}

//SendWsStatus code to client
func SendWsStatus(conn *websocket.Conn, level, code string, req int) (err error) {
	st := &stResult{Level: level, Code: code, Req: req}
	dataJSON, err := json.Marshal(st)
	if err != nil {
		logger.LOGE(err.Error())
		return
	}
	dataSend := make([]byte, len(dataJSON)+4)
	dataSend[0] = WSPktControl
	dataSend[1] = 0
	dataSend[2] = 0
	dataSend[3] = 0
	copy(dataSend[4:], dataJSON)
	err = conn.WriteMessage(websocket.BinaryMessage, dataSend)
	return
}

type stPlay struct {
	Name  string `json:"name"`
	Start int    `json:"start"`
	Len   int    `json:"len"`
	Reset int    `json:"reset"`
	Req   int    `json:"req"`
}

type stPlay2 struct {
	Name  string `json:"name"`
	Start int    `json:"start"`
	Len   int    `json:"len"`
	Reset int    `json:"reset"`
	Req   int    `json:"req"`
}

type stResume struct {
	Req int `json:"req"`
}

type stPause struct {
	Req int `json:"req"`
}

type stSeek struct {
	Offset int `json:"offset"`
	Req    int `json:"req"`
}

type stClose struct {
	Req int `json:"req"`
}

type stStop struct {
	Req int `json:"req"`
}

type stPublish struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Req  int    `json:"req"`
}

type stResult struct {
	Level string `json:"level"`
	Code  string `json:"code"`
	Req   int    `json:"req"`
}

// const string resource
const (
	NETCONNECTION_CALL_FAILED           = "NetConnection.Call.Failed"
	NETCONNECTION_CONNECT_APPSHUTDOWN   = "NetConnection.Connect.AppShutdown"
	NETCONNECTION_CONNECT_CLOSED        = "NetConnection.Connect.Closed"
	NETCONNECTION_CONNECT_FAILED        = "NetConnection.Connect.Failed"
	NETCONNECTION_CONNECT_IDLETIMEOUT   = "NetConnection.Connect.IdleTimeout"
	NETCONNECTION_CONNECT_INVALIDAPP    = "NetConnection.Connect.InvalidApp"
	NETCONNECTION_CONNECT_REJECTED      = "NetConnection.Connect.Rejected"
	NETCONNECTION_CONNECT_SUCCESS       = "NetConnection.Connect.Success"
	NETSTREAM_BUFFER_EMPTY              = "NetStream.Buffer.Empty"
	NETSTREAM_BUFFER_FLUSH              = "NetStream.Buffer.Flush"
	NETSTREAM_BUFFER_FULL               = "NetStream.Buffer.Full"
	NETSTREAM_FAILED                    = "NetStream.Failed"
	NETSTREAM_PAUSE_NOTIFY              = "NetStream.Pause.Notify"
	NETSTREAM_PLAY_FAILED               = "NetStream.Play.Failed"
	NETSTREAM_PLAY_FILESTRUCTUREINVALID = "NetStream.Play.FileStructureInvalid"
	NETSTREAM_PLAY_PUBLISHNOTIFY        = "NetStream.Play.PublishNotify"
	NETSTREAM_PLAY_RESET                = "NetStream.Play.Reset"
	NETSTREAM_PLAY_START                = "NetStream.Play.Start"
	NETSTREAM_PLAY_STOP                 = "NetStream.Play.Stop"
	NETSTREAM_PLAY_STREAMNOTFOUND       = "NetStream.Play.StreamNotFound"
	NETSTREAM_PLAY_UNPUBLISHNOTIFY      = "NetStream.Play.UnpublishNotify"
	NETSTREAM_PUBLISH_BADNAME           = "NetStream.Publish.BadName"
	NETSTREAM_PUBLISH_IDLE              = "NetStream.Publish.Idle"
	NETSTREAM_PUBLISH_START             = "NetStream.Publish.Start"
	NETSTREAM_RECORD_ALREADYEXISTS      = "NetStream.Record.AlreadyExists"
	NETSTREAM_RECORD_FAILED             = "NetStream.Record.Failed"
	NETSTREAM_RECORD_NOACCESS           = "NetStream.Record.NoAccess"
	NETSTREAM_RECORD_START              = "NetStream.Record.Start"
	NETSTREAM_RECORD_STOP               = "NetStream.Record.Stop"
	NETSTREAM_SEEK_FAILED               = "NetStream.Seek.Failed"
	NETSTREAM_SEEK_INVALIDTIME          = "NetStream.Seek.InvalidTime"
	NETSTREAM_SEEK_NOTIFY               = "NetStream.Seek.Notify"
	NETSTREAM_STEP_NOTIFY               = "NetStream.Step.Notify"
	NETSTREAM_UNPAUSE_NOTIFY            = "NetStream.Unpause.Notify"
	NETSTREAM_UNPUBLISH_SUCCESS         = "NetStream.Unpublish.Success"
	NETSTREAM_VIDEO_DIMENSIONCHANGE     = "NetStream.Video.DimensionChange"
)
