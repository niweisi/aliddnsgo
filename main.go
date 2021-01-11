package main

import (
	"encoding/json"
	"flag"
	"github.com/ahmetb/go-linq/v3"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"strings"
	"time"
)

var commandModel CommandModel
var configModel ConfigurationModel

func main() {
	initCommandModel()
	loadConfig()

	if commandModel.Interval == nil || *commandModel.Interval == 0 {
		update()
		return
	}

	intervalFunction()
}

func update() {

	subDomains := getSubDomains()

	for _, sub := range subDomains {
		mac := linq.From(*configModel.SubDomains).FirstWith(func(subDomain interface{}) bool {
			return subDomain.(SubDomainModel).Name == sub.RR
		}).(SubDomainModel).Mac

		publicIp := ""
		if "public" == mac {
			publicIp = getPublicIp()
		} else {
			publicIp = macToIp(mac)
		}

		log.Println(publicIp)

		if sub.Value != publicIp {
			// 更新域名绑定的 IP 地址。
			sub.Value = publicIp
			sub.TTL = linq.From(*configModel.SubDomains).FirstWith(func(subDomain interface{}) bool {
				return subDomain.(SubDomainModel).Name == sub.RR
			}).(SubDomainModel).Interval
			updateSubDomain(&sub)
		}
	}

	log.Printf("域名记录更新成功...")
}

func intervalFunction() {
	tick := time.Tick(time.Second * time.Duration(*commandModel.Interval))
	for {
		select {
		case <-tick:
			update()
		}
	}
}

func initCommandModel() {
	commandModel.FilePath = flag.String("f", "", "指定自定义的配置文件，请传入配置文件的路径。")
	commandModel.Interval = flag.Int("i", 0, "指定程序的自动检测周期，单位是秒。")

	flag.Parse()
}

func loadConfig() {
	var configFile string
	if *commandModel.FilePath == "" {
		dir, _ := os.Getwd()
		configFile = path.Join(dir, "settings.json")
	} else {
		configFile = *commandModel.FilePath
	}

	// 打开配置文件，并进行反序列化。
	f, err := os.Open(configFile)
	if err != nil {
		log.Fatalf("无法打开文件：%s", err)
		os.Exit(-1)
	}
	defer f.Close()
	data, _ := ioutil.ReadAll(f)

	if err := json.Unmarshal(data, &configModel); err != nil {
		log.Fatalf("数据反序列化失败：%s", err)
		os.Exit(-1)
	}
}

func getPublicIp() string {
	resp, err := http.Get(GetPublicIpUrl)
	if err != nil {
		log.Printf("获取公网 IP 出现错误，错误信息：%s", err)
		os.Exit(-1)
	}
	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	return strings.Replace(string(bytes), "\n", "", -1)
}

func getSubDomains() []alidns.Record {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", configModel.AccessId, configModel.AccessKey)

	request := alidns.CreateDescribeDomainRecordsRequest()
	request.Scheme = "https"

	request.DomainName = configModel.MainDomain

	response, err := client.DescribeDomainRecords(request)
	if err != nil {
		log.Println(err.Error())
	}

	// 过滤符合条件的子域名信息。
	var queryResult []alidns.Record
	linq.From(response.DomainRecords.Record).Where(func(c interface{}) bool {
		return linq.From(*configModel.SubDomains).Select(func(x interface{}) interface{} {
			return x.(SubDomainModel).Name
		}).Contains(c.(alidns.Record).RR)
	}).ToSlice(&queryResult)

	return queryResult
}

func updateSubDomain(subDomain *alidns.Record) {
	client, err := alidns.NewClientWithAccessKey("cn-hangzhou", configModel.AccessId, configModel.AccessKey)

	request := alidns.CreateUpdateDomainRecordRequest()
	request.Scheme = "https"
	request.RecordId = subDomain.RecordId
	request.RR = subDomain.RR
	request.Type = subDomain.Type
	request.Value = subDomain.Value
	request.TTL = requests.NewInteger64(subDomain.TTL)

	_, err = client.UpdateDomainRecord(request)
	if err != nil {
		log.Print(err.Error())
	}
}

// 根据mac地址获取本地ip
func macToIp(mac string) string {
	var ips = make(map[string]string)
	ips, _ = Ips()
	//key是网卡名称，value是网卡IP
	for k, v := range ips {
		if mac == k {
			return v
		}
	}
	return ""
}

func Ips() (map[string]string, error) {
	ips := make(map[string]string)
	//返回 interface 结构体对象的列表，包含了全部网卡信息
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	//遍历全部网卡
	for _, inter := range interfaces {

		mac := inter.HardwareAddr
		if mac.String() != "" {
		}

		// Addrs() 方法返回一个网卡上全部的IP列表
		address, err := inter.Addrs()
		if err != nil {
			return nil, err
		}

		//遍历一个网卡上全部的IP列表，组合为一个字符串，放入对应网卡名称的map中
		i := 0
		for _, v := range address {
			if 1 == i {
				str := v.String()
				content := str[0 : len(str)-3]
				ips[mac.String()] += content
			}
			i++
		}
	}
	return ips, nil
}
