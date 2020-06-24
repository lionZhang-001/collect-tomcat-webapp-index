/*
@Time : 2020/6/23 15:37
@Author : zhangyongyue
@File : basicStructure
@Software: GoLand
*/

package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strings"
)

var port1, port2 string
var catalina1, name1, name2 string
var catalina2 string
var attr1 []string = make([]string, 5)
var attr2 []string = make([]string, 30)
var logs *logrus.Logger = logrus.New()
var logsFile *os.File

var index2Mysql map[string]string = make(map[string]string)

func init() {
	//日志
	//logs = logrus.New()
	logs.SetFormatter(&logrus.TextFormatter{})
	logsFile, err := createLogsFile()
	if err != nil {
		fmt.Printf("creating logs file ", logsFile.Name(), " went wrong :", err)
		return
	}
	logs.SetOutput(logsFile)

	//读取connector port
	data, err := ioutil.ReadFile("./config/connector_port")
	if err != nil {
		logs.Info(err, " , unable to get connector information")
	}
	ports := strings.Split(string(data), "|")
	port2 := ports[1]

	catalina1 = "Catalina:type=GlobalRequestProcessor,name="

	name2 = "\"http-bio-" + port2 + "\""
	attr1 = []string{"bytesReceived", "bytesSent", "errorCount", "processingTime", "requestCount", "maxTime", "modelerType"}

	index2Mysql = map[string]string{
		catalina1 + name2 + ",attr=" + attr1[0]: "http_bytes_received",
		catalina1 + name2 + ",attr=" + attr1[1]: "http_bytes_sent",
		catalina1 + name2 + ",attr=" + attr1[2]: "http_error_count",
		catalina1 + name2 + ",attr=" + attr1[3]: "http_processing_time",
		catalina1 + name2 + ",attr=" + attr1[4]: "http_request_count",
	}

}

type tomcat struct {
	id, ip, port, identification string
}

var index []string = []string{"Catalina:type=Manager,context=*,host=localhost", "Catalina:type=GlobalRequestProcessor,name=*"}
