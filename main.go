/*
@Time : 2020/6/22 17:28
@Author : zhangyongyue
@File : main
@Software: GoLand
*/
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/newrelic/infra-integrations-sdk/jmx"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
	"time"
)

var filepath string

func init() {

	fptr := flag.String("fpath", "/home/inspur/icp-collect/config/collect-tomcat.cfg", "tomcat config file")
	flag.Parse()
	filepath = *fptr

}

func main() {

	var durationTime int
	var tom tomcat
	var line string
	var lines []string
	durationTime = 50000

	defer logsFile.Close()

	//采集过程
	f, err := os.Open(filepath)
	if err != nil {
		logs.Fatal(err)
	}
	//Delay closing files
	defer f.Close()

	//按行读取配置文件
	s := bufio.NewScanner(f)

	//Scan停止，或者扫描到尾部，或者出现错误，会返回false;
	//在扫描过程中，当Scan返回false时，Err方法会返回error，但是当是读取到文件尾部io.EOF, Err方法会返回nil
	for s.Scan() {

		line = s.Text()
		lines = strings.Split(line, "|")
		tom.id = lines[0]
		tom.ip = lines[1]
		tom.port = lines[2]
		tom.identification = lines[3]

		err = jmx.Open(tom.ip, tom.port, "", "")
		if err != nil {
			logs.WithFields(logrus.Fields{"tomcat_ip ": tom.ip}).Error(err)
			continue
		}

		//时间戳
		time1 := time.Unix(time.Now().Local().Unix()-int64(time.Now().Local().Minute())%5*60, 0)
		time2 := time.Unix(time.Now().Local().Unix()-int64(time.Now().Local().Minute())%5*60, 300000000000)

		//生成json文件
		jsonfile, err := os.OpenFile("../raw/"+tom.id+"-#-"+"TOMCAT-WEBAPP"+"-#-"+time1.Format("200601021504")+"-#-"+time2.Format("200601021504")+"-#-5.raw", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
		if err != nil {
			logs.Error(err)
			continue
		}

		all := make(map[string]interface{})
		all2Mysql := make(map[string]interface{})
		for _, value := range index {
			results, err := jmx.Query(value, durationTime)
			if err != nil {
				logs.Error(err)
				continue
			}

			for key, value := range results {
				all[key] = value
			}
		}
		logs.Info("getting tomcat index successes")

		/*jsontest, err := json.MarshalIndent(all, "", " ")
		if err != nil {
			fmt.Println("11")
			continue
		}
		fmt.Println(string(jsontest))
		fmt.Println(webapp)*/


		for key, value := range all {
			if value1 , ok := index2Mysql[key] ; ok {
				all2Mysql[string(value1)] = value
			}else {
				keys := strings.Split(key , "context=")
				if keys[0] == "Catalina:type=Manager," {
					appKeys := strings.Split(keys[1] , ",host=localhost,attr=")
					all2Mysql[appKeys[0]+"_"+strings.ToLower(appKeys[1])] = value
				}
			}
		}

		jsonstr, err := json.MarshalIndent(all2Mysql, "", "  ")
		if err != nil {
			logs.Error(err)
			continue
		}

		_, err = fmt.Fprintln(jsonfile, string(jsonstr))
		if err != nil {
			logs.Error(err)
			continue
		}
		logs.Info("writing json file successes")

		jmx.Close()
	}

}

func createLogsFile() (logFile *os.File, err error) {

	timeNow := time.Unix(time.Now().Local().Unix(), 0).Format("20060102")
	logsFileName := "../../../logs/tomcatWebappLog-" + string(timeNow) + ".log"
	_, err1 := os.Stat(logsFileName)
	if err1 != nil {
		if os.IsNotExist(err1) {
			f, err := os.OpenFile(logsFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			return f, err
		} else {
			fmt.Println(err1)
		}
	}
	f, err := os.OpenFile(logsFileName, os.O_APPEND|os.O_RDWR, 0666)

	return f, nil

}
