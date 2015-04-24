package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gosexy/redis"
	"github.com/huzorro/spfactor/sexredis"
	"testing"
)

var cache *Cache
var (
	user    *SpUser
	info    *SpInfo
	consign *SpConsign
	service *SpService
	msisdn  *SpMsisdn
)

func init() {
	rc := &sexredis.RedisPool{make(chan *redis.Client, 20), func() (*redis.Client, error) {
		client := redis.New()
		err := client.Connect("localhost", uint(6379))
		return client, err
	}}

	mc, _ := sql.Open("mysql", "sp:woai840511~@tcp(127.0.0.1:3306)/sp?charset=utf8")
	mc.SetMaxOpenConns(1)
	cache = New()
	cache.SetRedisConnect(rc)
	cache.SetMysqlConnect(mc)

	user = &SpUser{}
	info = &SpInfo{}
	consign = &SpConsign{}
	service = &SpService{}
	msisdn = &SpMsisdn{}

	cache.GetCache("sp_info:spnum:10669501", CACHE_TYPE_HASH, info)
	cache.GetCache("sp_consign:consignid:1", CACHE_TYPE_HASH, consign)
	cache.GetCache("sp_service:spnum:10669501:serviceword:t", CACHE_TYPE_HASH, service)
	cache.GetCache("sp_msisdn:prefix:1390000", CACHE_TYPE_HASH, msisdn)

	user.SpConsign = consign
	user.SpInfo = info
	user.SpMsisdn = msisdn
	user.SpService = service
	user.SpServiceRule = &SpServiceRule{"10669501", "13900005678", "ok", "deliver", "t", "123456", "0313", "zhangjiakou", "0311", "hebeisheng"}
}

func TestSetSpInfo(t *testing.T) {
	if err := cache.SetSpInfo(); err != nil {
		t.Error(err)
	}
}

func TestGetCacheSpInfo(t *testing.T) {
	spInfo := &SpInfo{}

	if err := cache.GetCache("sp_info:spnum:10669501", CACHE_TYPE_HASH, spInfo); err != nil {
		t.Error(err)
	} else {
		t.Log(spInfo)
		t.Log(spInfo.Spnum, spInfo.Spname)
	}
}

func TestSetSpConsign(t *testing.T) {
	if err := cache.SetSpConsign(); err != nil {
		t.Error(err)
	}
}

func TestGetCacheSPConsign(t *testing.T) {
	consign := &SpConsign{}
	if err := cache.GetCache("sp_consign:consignid:1", CACHE_TYPE_HASH, consign); err != nil {
		t.Error(err)
	} else {
		t.Log(consign)
		t.Log(consign.Consignid, consign.Consignname)
	}
}

func TestSetSpService(t *testing.T) {
	if err := cache.SetSpService(); err != nil {
		t.Error(err)
	}
}

func TestGetCacheSpServiceByID(t *testing.T) {
	service := &SpService{}
	if err := cache.GetCache("sp_service:serviceid:1", CACHE_TYPE_HASH, service); err != nil {
		t.Error(err)
	} else {
		t.Log(service)
		t.Log(service.Serviceid, service.Serviceword, service.Servicename,
			service.Servicetype, service.Servicefee,
			service.Serviceip, service.Servicerule, service.Referrule,
			service.Spnum)
		t.Log(service.Servicerule.Spnum, service.Servicerule.Terminal)
	}
}

func TestGetCacheSpServiceBySpnumAndWord(t *testing.T) {
	service := &SpService{}
	if err := cache.GetCache("sp_service:spnum:10669501:serviceword:t", CACHE_TYPE_HASH, service); err != nil {
		t.Error(err)
	} else {
		t.Log(service)
		t.Log(service.Serviceid, service.Serviceword, service.Servicename, service.Servicetype, service.Servicefee,
			service.Serviceip, service.Servicerule, service.Referrule,
			service.Spnum)
	}
}

func TestGetCacheSpServiceByIp(t *testing.T) {
	service := &SpService{}
	if err := cache.GetCache("sp_service:serviceip:192.168.1.222", CACHE_TYPE_HASH, service); err != nil {
		t.Error(err)
	} else {
		t.Log(service)
		t.Log(service.Serviceid, service.Serviceword, service.Servicename, service.Servicetype, service.Servicefee,
			service.Serviceip, service.Servicerule, service.Referrule,
			service.Spnum)
	}
}

func TestSetSpCp(t *testing.T) {
	if err := cache.SetSpCp(); err != nil {
		t.Error(err)
	}
}

func TestGetCacheCpByServiceid(t *testing.T) {
	cp := &SpCp{}
	if err := cache.GetCache("sp_cp:serviceid:1", CACHE_TYPE_HASH, cp); err != nil {
		t.Error(err)
	} else {
		t.Log(cp)
		t.Log(cp.Cpid, cp.Consignid, cp.Serviceid)
	}
}

func TestSetSpMsisdn(t *testing.T) {
	if err := cache.SetSpMsisdn(); err != nil {
		t.Error(err)
	}
}

func TestGetCacheSpMsisdn(t *testing.T) {
	msisdn := &SpMsisdn{}
	if err := cache.GetCache("sp_msisdn:prefix:1390000", CACHE_TYPE_HASH, msisdn); err != nil {
		t.Error(err)
	} else {
		t.Log(msisdn)
		t.Log(msisdn.Id, msisdn.Prefix, msisdn.Provinceid, msisdn.Provincename, msisdn.Cityid, msisdn.Cityname)
	}
}

func TestSetCityToProvince(t *testing.T) {
	if err := cache.SetCityToProvince(); err != nil {
		t.Error(err)
	}
}

func TestGetCacheCityToProvince(t *testing.T) {
	msisdn := &SpMsisdn{}
	if err := cache.GetCache("sp_msisdn:cityid:0313", CACHE_TYPE_HASH, msisdn); err != nil {
		t.Error(err)
	} else {
		t.Log(msisdn)
		t.Log(msisdn.Id, msisdn.Prefix, msisdn.Provinceid, msisdn.Provincename, msisdn.Cityid, msisdn.Cityname)
	}
}

func TestSetProvince(t *testing.T) {
	if err := cache.SetProvince(); err != nil {
		t.Error(err)
	}
}

func TestGetCacheProvince(t *testing.T) {
	msisdn := &SpMsisdn{}
	if err := cache.GetCache("sp_msisdn:provinceid:0311", CACHE_TYPE_HASH, msisdn); err != nil {
		t.Error(err)
	} else {
		t.Log(msisdn)
		t.Log(msisdn.Id, msisdn.Prefix, msisdn.Provinceid, msisdn.Provincename, msisdn.Cityid, msisdn.Cityname)
		if msisdn.Prefix == "" {
			t.Log("prefix is empty")
		}
	}
}

func TestSetSpSink(t *testing.T) {
	if err := cache.SetSpSink(); err != nil {
		t.Error(err)
	}
}

func TestGetCacheSpSink(t *testing.T) {
	sink := &SpSink{}

	if err := cache.GetCache("sp_sink:consignid:1:serviceid:1:provincename:hebeisheng", CACHE_TYPE_HASH, sink); err == nil {
		t.Log(sink)
		t.Log(sink.Sinkid, sink.Bnum, sink.Slice, sink.Consignid, sink.Provincename, sink.Serviceid, sink.Response, sink.Tconsignid, sink.Key, sink.Cnum)
	} else if err := cache.GetCache("sp_sink:consignid:1:serviceid:1", CACHE_TYPE_HASH, sink); err == nil {
		t.Log(sink)
		t.Log(sink.Sinkid, sink.Bnum, sink.Slice, sink.Consignid, sink.Provincename, sink.Serviceid, sink.Response, sink.Tconsignid, sink.Key, sink.Cnum)
	} else if err := cache.GetCache("sp_sink:consignid:1:provincename:hebeisheng", CACHE_TYPE_HASH, sink); err == nil {
		t.Log(sink)
		t.Log(sink.Sinkid, sink.Bnum, sink.Slice, sink.Consignid, sink.Provincename, sink.Serviceid, sink.Response, sink.Tconsignid, sink.Key, sink.Cnum)
	} else {
		t.Error(err)
	}

}

func TestUpdateMoUsers(t *testing.T) {
	recode, rbool, err := cache.UpdateMoUsers(user)
	t.Log(recode, rbool)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateMtUsers(t *testing.T) {
	recode, rbool, err := cache.UpdateMtUsers(user)
	t.Log(recode, rbool)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateSinkCnum(t *testing.T) {
	recode, err := cache.UpdateSinkCnum(user)
	t.Log(recode, err)
	if err != nil {
		t.Error(err)
	}
}

func TestGetTodayMoUsers(t *testing.T) {
	rnum, err := cache.GetTodayMoUsers(user)
	t.Log(rnum)
	if err != nil {
		t.Error(err)
	}
}

func TestGetTodayMtUsers(t *testing.T) {
	rnum, err := cache.GetTodayMtUsers(user)
	t.Log(rnum)
	if err != nil {
		t.Error(err)
	}
}

func TestGetYdayMoUsers(t *testing.T) {
	rnum, err := cache.GetYdayMoUsers(user)
	t.Log(rnum)
	if err != nil {
		t.Error(err)
	}
}

func TestGetYdayMtUsers(t *testing.T) {
	rnum, err := cache.GetYdayMtUsers(user)
	t.Log(rnum)
	if err != nil {
		t.Error(err)
	}
}
