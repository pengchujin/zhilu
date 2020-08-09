package DnsUtils

import (
	"log"
	"github.com/miekg/dns"
	"github.com/oschwald/geoip2-golang"
	"net"
	"context"
	"encoding/json"
	"github.com/m13253/dns-over-https/json-dns"
	"strconv"
	"time"
	"strings"
	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func QueryDns(name string, rrTypeStr string, ednsClientSubnet string) *jsonDNS.Response  {

	optRedis, err := redis.ParseURL("redis://localhost:6379/0")
	if err != nil {
			panic(err)
	}
	rdb := redis.NewClient(optRedis)

	db, err := geoip2.Open("GeoLite2-Country.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	comment := ""
	// format dnsType
	rrType := uint16(1)
	if rrTypeStr == "" {
		} else if v, err := strconv.ParseUint(rrTypeStr, 10, 16); err == nil {
			rrType = uint16(v)
		} else if v, ok := dns.StringToType[strings.ToUpper(rrTypeStr)]; ok {
			rrType = v
		} else {
			comment = "Wrong DNStype"
	}

	log.Println(comment)


	r := new(dns.Msg)
	ednsClientFamily := uint16(0)
	ednsClientAddress := net.IP(nil)
	ednsClientNetmask := uint8(255)

	if ednsClientSubnet != "" {
		if ednsClientSubnet == "0/0" {
			ednsClientSubnet = "0.0.0.0/0"
		}
		slash := strings.IndexByte(ednsClientSubnet, '/')
		if slash < 0 {
			ednsClientAddress = net.ParseIP(ednsClientSubnet)
			if ednsClientAddress == nil {
				comment = "Invalid argument value: \"edns_client_subnet\" = %q"
			}
			if ipv4 := ednsClientAddress.To4(); ipv4 != nil {
				ednsClientFamily = 1
				ednsClientAddress = ipv4
				ednsClientNetmask = 24
			} else {
				ednsClientFamily = 2
				ednsClientNetmask = 56
			}
		} else {
			ednsClientAddress = net.ParseIP(ednsClientSubnet[:slash])
			if ednsClientAddress == nil {
				comment = "Invalid argument value: \"edns_client_subnet\" = %q"
			}
			if ipv4 := ednsClientAddress.To4(); ipv4 != nil {
				ednsClientFamily = 1
				ednsClientAddress = ipv4
			} else {
				ednsClientFamily = 2
			}
			netmask, err := strconv.ParseUint(ednsClientSubnet[slash+1:], 10, 8)
			if err != nil {
				comment = "Invalid argument value: \"edns_client_subnet\" = %q"
			}
			ednsClientNetmask = uint8(netmask)
		}
	} 

	opt := new(dns.OPT)
	opt.Hdr.Name = "."
	opt.Hdr.Rrtype = dns.TypeOPT
	opt.SetUDPSize(dns.DefaultMsgSize)
	opt.SetDo(true)

	if ednsClientAddress != nil {
		edns0Subnet := new(dns.EDNS0_SUBNET)
		edns0Subnet.Code = dns.EDNS0SUBNET
		edns0Subnet.Family = ednsClientFamily
		edns0Subnet.SourceNetmask = ednsClientNetmask
		edns0Subnet.SourceScope = 0
		edns0Subnet.Address = ednsClientAddress
		opt.Option = append(opt.Option, edns0Subnet)
	}
	r.Extra = append(r.Extra, opt)
	name = dns.Fqdn(name)


	log.Println(name)
	r.SetQuestion(name, rrType)
	c := new(dns.Client)

	record, err := db.City(ednsClientAddress)
	countryCode := record.Country.IsoCode

	redisKey := string(name) +  strconv.Itoa(int(rrType)) + "." + string(countryCode)
	val, err := rdb.Get(ctx, redisKey).Result()
	log.Println("---------------------------------------", val, err, rrType, redisKey, "---------------------------------------")
	if err == redis.Nil {
		res, _, err := c.Exchange(r, "8.8.8.8:53")

		var	ttl time.Duration
		jsonRes := jsonDNS.Marshal(res)

		if err == nil && len(jsonRes.Answer) > 0{
			ttl = time.Duration(jsonRes.Answer[0].TTL) * 100000000
		}
		resRedis, err := json.Marshal(jsonRes)
		log.Println("TTLLLLLLLLLLLLLLLLLLL", ttl)
		error := rdb.Set(ctx, redisKey, string(resRedis), ttl).Err()
		if error != nil {
			log.Println(error, "redis rrrrrrrrrrrrrrrrrrrrrrrrrrrrr")
		}
		return jsonRes
	} else {
		jsonRes := new(jsonDNS.Response)
		log.Println("Cacheeeeeeeeeeeeeeeee  dns")
		err := json.Unmarshal([]byte(val), jsonRes)
		if err != nil {
			log.Println(err)
		}
		return jsonRes
	}

}
