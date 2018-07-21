package registry

import (
	"strconv"
	"strings"

	motan "github.com/weibocom/motan-go/core"
	"github.com/weibocom/motan-go/log"
)

// 所有的服务通过URL来引用
// URL没有特别的Schema，可以是普通字符串，也可以是ip:port
// 直接注册模式下，URL就是ip:port
// 注册基本没任何事情，Discover直接将url变成[]*motan.URL
type DirectRegistry struct {
	url  *motan.URL
	urls []*motan.URL
}

func (d *DirectRegistry) GetURL() *motan.URL {
	return d.url
}
func (d *DirectRegistry) SetURL(url *motan.URL) {
	d.url = url
	d.urls = parseURLs(url)
}
func (d *DirectRegistry) GetName() string {
	return "direct"
}

func (d *DirectRegistry) InitRegistry() {
}
func (d *DirectRegistry) Subscribe(url *motan.URL, listener motan.NotifyListener) {
}

func (d *DirectRegistry) Unsubscribe(url *motan.URL, listener motan.NotifyListener) {
}

func (d *DirectRegistry) Discover(url *motan.URL) []*motan.URL {
	if d.urls == nil {
		d.urls = parseURLs(d.url)
	}
	result := make([]*motan.URL, 0, len(d.urls))
	for _, u := range d.urls {
		newURL := *url
		newURL.Host = u.Host
		newURL.Port = u.Port
		result = append(result, &newURL)
	}
	return result
}
func (d *DirectRegistry) Register(serverURL *motan.URL) {
	vlog.Infof("direct registry:register url :%+v\n", serverURL)
}
func (d *DirectRegistry) UnRegister(serverURL *motan.URL) {

}
func (d *DirectRegistry) Available(serverURL *motan.URL) {

}
func (d *DirectRegistry) Unavailable(serverURL *motan.URL) {

}
func (d *DirectRegistry) GetRegisteredServices() []*motan.URL {
	return nil
}
func (d *DirectRegistry) StartSnapshot(conf *motan.SnapshotConf) {}

func parseURLs(url *motan.URL) []*motan.URL {
	urls := make([]*motan.URL, 0)
	if len(url.Host) > 0 && url.Port > 0 {
		urls = append(urls, url)
	} else if address, exist := url.Parameters[motan.AddressKey]; exist {
		for _, add := range strings.Split(address, ",") {
			hostport := strings.Split(add, ":")
			if len(hostport) == 2 {
				port, err := strconv.Atoi(hostport[1])
				if err == nil {
					u := &motan.URL{Host: hostport[0], Port: port}
					urls = append(urls, u)
				}
			}
		}
	} else {
		vlog.Warningf("direct registry parse fail.url:L %+v", url)
	}
	return urls
}
