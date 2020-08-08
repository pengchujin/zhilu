package Services

import (
	"log"
	"github.com/miekg/dns"
	"github.com/m13253/dns-over-https/json-dns"
)

func GetDnsRes(name string) *jsonDNS.Response {
	jsonRes := new(jsonDNS.Response)
	r := new(dns.Msg)
	log.Println(name)
	r.SetQuestion("qq.com.", dns.TypeA)
	c := new(dns.Client)
	res, _, err := c.Exchange(r, "8.8.8.8:53")
	log.Println(jsonRes, res, err)
	jsonRes = jsonDNS.Marshal(res)
	return jsonRes
}