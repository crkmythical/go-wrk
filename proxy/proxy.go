package proxy

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync/atomic"
)

type ProxyConfig struct {
	index     uint64
	proxyList []string
	Activeed  bool
}

type Entry struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
	URL  string `json:"url"`
}

var gProxyConfig ProxyConfig

func InitProxyCfg(dir string) error {
	urlMap := make(map[string]struct{})
	// 遍历目录中的所有文件
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" {
			// 读取 JSON 文件
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			// 解析 JSON 数据
			var entries []Entry
			if err := json.Unmarshal(data, &entries); err != nil {
				return err
			}

			// 提取 URL 并去重
			for _, entry := range entries {
				urlMap[entry.URL] = struct{}{}
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	// 将去重后的 URL 存入列表
	uniqueURLs := make([]string, 0, len(urlMap))
	for url := range urlMap {
		uniqueURLs = append(uniqueURLs, url)
	}
	gProxyConfig.proxyList = uniqueURLs
	gProxyConfig.Activeed = true
	return nil
}

func IsProxyNeed() bool {
	return gProxyConfig.Activeed
}

func GetProxy() string {
	Len := len(gProxyConfig.proxyList)
	if Len == 0 {
		return ""
	}
	index := atomic.SwapUint64(&gProxyConfig.index, gProxyConfig.index+1)
	return gProxyConfig.proxyList[index%uint64(Len)]
}
