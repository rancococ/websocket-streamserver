// Copyright 2018 The use-go websocket-streamserver Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wssAPI

import (
	"regexp"
	"strings"
	"testing"

	"github.com/use-go/websocket-streamserver/logger"
)

func init() {

}

func TestDecodeWSPCtrlMsg(t *testing.T) {
	data := []byte("WSP/1.1 INIT\r\nproto: rtsp\r\nhost: 127.0.0.1\r\nport: 8554\r\nseq: 1\r\n\r\n")

	// Check the Data beganing
	protocalDataFormatCheck := regexp.MustCompile(`WSP/1\.1\s+\w+`)
	if protocalDataFormatCheck.Match(data) {

		contentString := string(data)
		logger.LOGD(contentString)
		strlines := strings.Split(contentString, "\r\n")
		arryLenth := len(strlines)
		// \r\n\r\n in the end ,2 empty string occured, Must more 3
		if arryLenth > 2 {

			headLineFormatReg := regexp.MustCompile(`[\w\.\\/]+`)
			targetArray := headLineFormatReg.FindAllString(strlines[0], -1)
			logger.LOGD("protocol version : " + strlines[0])
			parseredMsg := WSPMessage{Headers: map[string]string{}}
			parseredMsg.MsgType = targetArray[1]

			for i := 1; i < arryLenth-2; i++ {
				//get the version
				tmpStr := strlines[i]
				if len(tmpStr) < 1 {
					continue
				}
				//tmpStr = strings.TrimSpace(tmpStr)
				tmpStr = strings.Replace(tmpStr, " ", "", -1)
				tmpStrArry := strings.SplitN(tmpStr, ":", 2)
				parseredMsg.Headers[tmpStrArry[0]] = tmpStrArry[1]
			}
			return
		}
		t.Fatalf("message body: got wrong, can not parsing the fileld")
		//controlProcDes := data[versionIndex+2:]
	}
}

func TestEncodeWSPCtrlMsg(t *testing.T) {

	protocalDataFormatCheck := regexp.MustCompile(`WSP/1\.1\s+\w+`)

	header := map[string]string{"channel": "channel", "channel2": "channel2"}

	protocolStr := EncodeWSPCtrlMsg("200 OK", 1234, header, "")

	if !protocalDataFormatCheck.MatchString(protocolStr) {

		t.Fatalf("message format Wrong")

	}

	t.Log(protocolStr)

}