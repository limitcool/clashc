package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"gopkg.in/yaml.v2"
)

var (
	proxyName   string
	proxiesName string
	Config      Clash
	err         error
	ProxiesNum  int
	ProxyNum    int
)

type Clash struct {
	ExternalController string `yaml:"external-controller"`
	port               int    `yaml:"port"`
	socksPort          int    `yaml:"socks-port"`
	allowLan           bool   `yaml:"allow-lan"`
	mode               string `yaml:"mode"`
	logLevel           string `yaml:"log-level"`
	// proxies            []string     `yaml:"proxies"`
	ProxyGroups []ProxyGroup `yaml:"proxy-groups"`
	// Rules       []string     `yaml:"rules"`
}

type ProxyGroup struct {
	Name      string   `yaml:"name"`
	Type      string   `yaml:"type"`
	Url       string   `yaml:"url"`
	Interval  int      `yaml:"interval"`
	Tolerance int      `yaml:"tolerance"`
	Proxies   []string `yaml:"proxies"`
}

func main() {

	Config, err = UnmarshalYaml("/root/.config/clash/config.yaml")
	// Config, err = UnmarshalYaml("C:\\Users\\Andorid\\.config\\clash\\profiles\\1649774903663.yml")
	if err != nil {
		log.Println(err)
		return
	}
	err = SetProxie(GetProxiesName(), GetProxyName())
	if err != nil {
		log.Fatal(err)
	}
	log.Println("切换节点完成🆗")
}

// yaml反序列化
func UnmarshalYaml(yamlPath string) (Clash, error) {
	var clash Clash
	yamlFile, err := ioutil.ReadFile(yamlPath)
	if err != nil {
		log.Println(err)
		return clash, err
	}
	err = yaml.Unmarshal(yamlFile, &clash)
	if err != nil {
		log.Println(err)
		return clash, err
	}
	// log.Println(clash)
	// log.Panic(clash.ExternalController)
	return clash, nil
}

// 设置代理
func SetProxie(proxiesName string, proxyName string) error {
	log.Println("要切换的分组为:", proxiesName, "要切换的节点为:", proxyName)
	// Go的UrlEncode有坑点,需要将  + 转换为 %20
	url := "http://127.0.0.1" + Config.ExternalController + "/proxies/" + strings.Replace(url.QueryEscape(proxiesName), "+", "%20", -1)
	method := "PUT"
	payload := strings.NewReader(fmt.Sprintf("{\"name\":\"%s\"}", proxyName))
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println(err)
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return err
	}
	if body != nil {
		return errors.New(string(body))
	}
	return nil
}

// 获取分组名称
func GetProxiesName() string {
	for i, v := range Config.ProxyGroups {
		fmt.Println(i+1, ":", v.Name)
	}
	// 获取用户输入的分组
	fmt.Println("请输入要切换的分组序号")
	fmt.Scanln(&ProxiesNum)
	if ProxiesNum > len(Config.ProxyGroups) || ProxiesNum <= 0 {
		log.Fatal("输入分组序号错误!")
	}
	return Config.ProxyGroups[ProxiesNum-1].Name
}

// 获取节点名称

func GetProxyName() string {
	log.Println(Config.ProxyGroups[ProxiesNum-1].Proxies)
	for i, v := range Config.ProxyGroups[ProxiesNum-1].Proxies {
		fmt.Println(i+1, ":", v)
	}

	// 获取用户输入的节点名称
	fmt.Println("请输入要切换的节点序号")
	fmt.Scanln(&ProxyNum)
	if ProxyNum > len(Config.ProxyGroups[ProxiesNum-1].Proxies) || ProxyNum <= 0 {
		log.Fatal("输入分组序号错误!")
	}
	log.Println(Config.ProxyGroups[ProxiesNum-1].Proxies[ProxyNum-1])
	return Config.ProxyGroups[ProxiesNum-1].Proxies[ProxyNum-1]
}
