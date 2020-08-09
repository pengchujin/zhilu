package Services

import (
	"github.com/pengchujin/zhilu/DnsUtils"
	"github.com/m13253/dns-over-https/json-dns"
)

func GetDnsRes(name string, rrTypeStr string, ednsClientSubnet string) *jsonDNS.Response {
	jsonRes := DnsUtils.QueryDns(name, rrTypeStr, ednsClientSubnet)
	return jsonRes
}
