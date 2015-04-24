package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gosexy/redis"
	"github.com/huzorro/spfactor/sexredis"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	CACHE_TYPE_HASH = iota
	CACHE_TYPE_SET
	CACHE_TYPE_SET_NUM
)

const (
	SERVICE_TYPE_BOTH = iota
	SERVICE_TYPE_ONLY_MT
)

const (
	STATUS_ID_SUCCESS = "200"
)

type SpInfo struct {
	Spnum  string `json:"spnum"`
	Spname string `json:"spname"`
}

type SpConsign struct {
	Consignid   int    `json:"consignid"`
	Consignname string `json:"consignname"`
}

type SpService struct {
	Serviceid   int           `json:"serviceid"`
	Serviceword string        `json:"serviceword"`
	Servicename string        `json:"servicename"`
	Servicetype int           `json:"servicetype"`
	Servicefee  int           `json:"servicefee"`
	Serviceip   string        `json:"serviceip"`
	Servicerule SpServiceRule `json:"servicerule"`
	Referrule   SpReferRule   `json:"referrule"`
	Spnum       string        `json:"spnum"`
}

type SpCp struct {
	Cpid      int `json:"cpid"`
	Consignid int `json:"consignid"`
	Serviceid int `json:"serviceid"`
}

type SpMsisdn struct {
	Id           int    `json:"id"`
	Prefix       string `json:"prefix"`
	Provinceid   string `json:"provinceid"`
	Provincename string `json:"provincename"`
	Cityid       string `json:"cityid"`
	Cityname     string `json:"cityname"`
}
type SpSink struct {
	Sinkid       int            `json:"id"`
	Bnum         int            `json:"bnum"`
	Slice        int            `json:"slice"`
	Consignid    int            `json:"consignid"`
	Serviceid    int            `json:"serviceid"`
	Provincename string         `json:"provincename"`
	Response     SpSinkResponse `json:"response"`
	Tconsignid   int            `json:"tconsignid"`
	Key          string         `json:"key"`
	Cnum         int            `json:"cnum"`
}
type SpServiceRule struct {
	Spnum        string `json:"spnum"`
	Terminal     string `json:"terminal"`
	Statusid     string `json:"statusid"`
	Statuscode   string `json:"statuscode"`
	Serviceword  string `json:"serviceword"`
	Linkid       string `json:"linkid"`
	Cityid       string `json:"cityid"`
	Cityname     string `json:"cityname"`
	Provinceid   string `json:"provinceid"`
	Provincename string `json:"provincename`
	Timeline     string `json:"timeline"`
	Expendtime   string `json:"expendtime""`
}
type SpUser struct {
	SpInfo        SpInfo        `json:"spinfo"`
	SpConsign     SpConsign     `json:"spconsign"`
	SpService     SpService     `json:"spservice"`
	SpMsisdn      SpMsisdn      `json:"spmsisdn"`
	SpServiceRule SpServiceRule `json:"spservicerule"`
	SpOconsign    SpConsign     `json:"spoconsign"`
}

type SpReferRule struct {
	Key      []string `json:"key"`
	Statusid string   `json:"statusid"`
}

type SpSinkResponse struct {
	Response []string `json:"response"`
}

type Cache struct {
	db   *sql.DB
	pool *sexredis.RedisPool
}

func New() *Cache {
	c := new(Cache)
	return c
}

//func Classic(pool *sexredis.RedisPool, db *sql.DB) *Cache {
//	c := New()
//	c.SetRedisConnect(pool)
//	c.SetMysqlConnect(db)
//	return c
//}

func (self *Cache) SetRedisConnect(rc *sexredis.RedisPool) {
	self.pool = rc
}

func (self *Cache) SetMysqlConnect(db *sql.DB) {
	self.db = db
}

func (self *Cache) SetSpInfo() error {

	stmtOut, err := self.db.Prepare("SELECT spnum, spname FROM sp_info")
	if err != nil {
		return err
	}
	rc, err := self.pool.Get()
	if err != nil {
		return err
	}
	defer func() {
		stmtOut.Close()
		self.pool.Close(rc)
	}()
	rows, err := stmtOut.Query()
	if err != nil {
		return err
	}
	spinfo := SpInfo{}
	for rows.Next() {
		if err := rows.Scan(&spinfo.Spnum, &spinfo.Spname); err != nil {
			return err
		}
		jsonstr, err := json.Marshal(spinfo)
		if err != nil {
			return err
		}
		if _, err := rc.HMSet("sp_info:spnum:"+spinfo.Spnum, "spnum", spinfo.Spnum, "spname", spinfo.Spname, "json", jsonstr); err != nil {
			return err
		}
	}
	return nil
}

func (self *Cache) SetSpConsign() error {
	stmtOut, err := self.db.Prepare("SELECT consignid, consignname FROM sp_consign")
	if err != nil {
		return err
	}
	rc, err := self.pool.Get()
	if err != nil {
		return err
	}
	defer func() {
		stmtOut.Close()
		self.pool.Close(rc)
	}()
	rows, err := stmtOut.Query()
	if err != nil {
		return err
	}
	spConsign := SpConsign{}
	for rows.Next() {
		if err := rows.Scan(&spConsign.Consignid, &spConsign.Consignname); err != nil {
			return err
		}
		jsonstr, err := json.Marshal(spConsign)
		if err != nil {
			return err
		}
		if _, err := rc.HMSet("sp_consign:consignid:"+strconv.Itoa(spConsign.Consignid), "consignid", spConsign.Consignid, "consignname", spConsign.Consignname, "json", jsonstr); err != nil {
			return err
		}
	}
	return nil

}

func (self *Cache) SetSpService() error {
	stmtOut, err := self.db.Prepare("SELECT serviceid, serviceword, servicename, servicetype, servicefee, serviceip, servicerule, referrule, spnum FROM sp_service")
	if err != nil {
		return err
	}
	rc, err := self.pool.Get()
	if err != nil {
		return err
	}
	defer func() {
		stmtOut.Close()
		self.pool.Close(rc)
	}()

	rows, err := stmtOut.Query()
	if err != nil {
		return err
	}
	spService := SpService{}
	var serviceRuleStr, referRuleStr string
	for rows.Next() {
		if err := rows.Scan(&spService.Serviceid, &spService.Serviceword, &spService.Servicename,
			&spService.Servicetype, &spService.Servicefee, &spService.Serviceip, &serviceRuleStr, &referRuleStr, &spService.Spnum); err != nil {
			return err
		}
		serviceRule := SpServiceRule{}
		if err := json.Unmarshal([]byte(serviceRuleStr), &serviceRule); err != nil {
			return err
		}
		spService.Servicerule = serviceRule
		referRule := SpReferRule{}
		if err := json.Unmarshal([]byte(referRuleStr), &referRule); err != nil {
			return err
		}
		spService.Referrule = referRule
		jsonstr, err := json.Marshal(spService)
		if err != nil {
			return err
		}

		if _, err := rc.HMSet("sp_service:serviceid:"+strconv.Itoa(spService.Serviceid),
			"serviceid", spService.Serviceid,
			"serviceword", spService.Serviceword, "servicename", spService.Servicename,
			"servicetype", spService.Servicetype, "servicefee", spService.Servicefee,
			"serviceip", spService.Serviceip, "servicerule", serviceRuleStr, "referrule", referRuleStr,
			"spnum", spService.Spnum, "json", jsonstr); err != nil {
			return err
		}
		if _, err := rc.HMSet("sp_service:spnum:"+spService.Spnum+":serviceword:"+spService.Serviceword,
			"serviceid", spService.Serviceid,
			"serviceword", spService.Serviceword, "servicename", spService.Servicename,
			"servicetype", spService.Servicetype, "servicefee", spService.Servicefee,
			"serviceip", spService.Serviceip, "servicerule", serviceRuleStr, "referrule", referRuleStr,
			"spnum", spService.Spnum, "json", jsonstr); err != nil {
			return err
		}
		if _, err := rc.HMSet("sp_service:serviceip:"+spService.Serviceip,
			"serviceid", spService.Serviceid,
			"serviceword", spService.Serviceword, "servicename", spService.Servicename,
			"servicetype", spService.Servicetype, "servicefee", spService.Servicefee,
			"serviceip", spService.Serviceip, "servicerule", serviceRuleStr, "referrule", referRuleStr,
			"spnum", spService.Spnum, "json", jsonstr); err != nil {
			return err
		}
	}
	return nil
}
func (self *Cache) SetSpCp() error {
	stmtOut, err := self.db.Prepare("SELECT cpid, consignid, serviceid FROM sp_cp")
	if err != nil {
		return err
	}
	rc, err := self.pool.Get()
	if err != nil {
		return err
	}
	defer func() {
		stmtOut.Close()
		self.pool.Close(rc)
	}()
	rows, err := stmtOut.Query()
	if err != nil {
		return err
	}
	cp := SpCp{}
	for rows.Next() {
		if err := rows.Scan(&cp.Cpid, &cp.Consignid, &cp.Serviceid); err != nil {
			return err
		}
		jsonstr, err := json.Marshal(cp)
		if err != nil {
			return err
		}
		if _, err := rc.HMSet("sp_cp:serviceid:"+strconv.Itoa(cp.Serviceid), "cpid", cp.Cpid,
			"consignid", cp.Consignid, "serviceid", cp.Serviceid, "json", jsonstr); err != nil {
			return err
		}
	}
	return nil
}

func (self *Cache) SetSpMsisdn() error {
	stmtOut, err := self.db.Prepare("SELECT id, prefix, provinceid, provincename, cityid, cityname FROM sp_msisdn")
	if err != nil {
		return err
	}
	rc, err := self.pool.Get()
	if err != nil {
		return err
	}
	defer func() {
		stmtOut.Close()
		self.pool.Close(rc)
	}()
	rows, err := stmtOut.Query()
	if err != nil {
		return err
	}
	msisdn := SpMsisdn{}
	for rows.Next() {
		if err := rows.Scan(&msisdn.Id, &msisdn.Prefix, &msisdn.Provinceid, &msisdn.Provincename, &msisdn.Cityid, &msisdn.Cityname); err != nil {
			return err
		}
		jsonstr, err := json.Marshal(msisdn)
		if err != nil {
			return err
		}
		if _, err := rc.HMSet("sp_msisdn:prefix:"+msisdn.Prefix, "id", msisdn.Id, "prefix", msisdn.Prefix,
			"provinceid", msisdn.Provinceid, "provincename", msisdn.Provincename, "cityid", msisdn.Cityid, "cityname", msisdn.Cityname,
			"json", jsonstr); err != nil {
			return err
		}
	}
	return nil
}

func (self *Cache) SetCityToProvince() error {
	stmtOut, err := self.db.Prepare("SELECT provinceid, provincename, cityid, cityname FROM sp_msisdn GROUP BY cityid, cityname")
	if err != nil {
		return err
	}
	rc, err := self.pool.Get()
	if err != nil {
		return err
	}
	defer func() {
		stmtOut.Close()
		self.pool.Close(rc)
	}()
	rows, err := stmtOut.Query()
	if err != nil {
		return err
	}
	msisdn := SpMsisdn{}
	for rows.Next() {
		if err := rows.Scan(&msisdn.Provinceid, &msisdn.Provincename, &msisdn.Cityid, &msisdn.Cityname); err != nil {
			return err
		}
		jsonstr, err := json.Marshal(msisdn)
		if err != nil {
			return err
		}
		if _, err := rc.HMSet("sp_msisdn:cityid:"+msisdn.Cityid, "provinceid", msisdn.Provinceid, "provincename",
			msisdn.Provincename, "cityid", msisdn.Cityid, "cityname", msisdn.Cityname,
			"json", jsonstr); err != nil {
			return err
		}
	}
	return nil
}

func (self *Cache) SetProvince() error {
	stmtOut, err := self.db.Prepare("SELECT provinceid, provincename FROM sp_msisdn GROUP BY provinceid, provincename")
	if err != nil {
		return err
	}
	rc, err := self.pool.Get()
	if err != nil {
		return err
	}
	defer func() {
		stmtOut.Close()
		self.pool.Close(rc)
	}()
	rows, err := stmtOut.Query()
	if err != nil {
		return err
	}
	msisdn := SpMsisdn{}
	for rows.Next() {
		if err := rows.Scan(&msisdn.Provinceid, &msisdn.Provincename); err != nil {
			return err
		}
		jsonstr, err := json.Marshal(msisdn)
		if err != nil {
			return err
		}
		if _, err := rc.HMSet("sp_msisdn:provinceid:"+msisdn.Provinceid, "provinceid", msisdn.Provinceid, "provincename",
			msisdn.Provincename, "json", jsonstr); err != nil {
			return err
		}
	}
	return nil
}

func (self *Cache) SetSpSink() error {
	stmtOut, err := self.db.Prepare("SELECT sinkid, bnum, slice, consignid, serviceid, provincename, response,tconsignid, CONCAT('consignid:', consignid, IF(serviceid, CONCAT(':serviceid:', serviceid), ''), IF(provincename <> '', CONCAT(':provincename:', provincename), '')) AS key1 FROM sp_sink")
	if err != nil {
		return err
	}
	rc, err := self.pool.Get()
	if err != nil {
		return err
	}
	defer func() {
		stmtOut.Close()
		self.pool.Close(rc)
	}()
	rows, err := stmtOut.Query()
	if err != nil {
		return err
	}
	sink := SpSink{}
	var responseStr string
	for rows.Next() {
		if err := rows.Scan(&sink.Sinkid, &sink.Bnum, &sink.Slice, &sink.Consignid,
			&sink.Serviceid, &sink.Provincename, &responseStr, &sink.Tconsignid, &sink.Key); err != nil {
			return err
		}
		response := SpSinkResponse{}
		if err := json.Unmarshal([]byte(responseStr), &response); err != nil {
			return err
		}
		sink.Response = response
		jsonstr, err := json.Marshal(sink)
		if err != nil {
			return err
		}
		if _, err := rc.HMSet("sp_sink:"+sink.Key, "sinkid", sink.Sinkid, "bnum", sink.Bnum,
			"slice", sink.Slice, "consignid", sink.Consignid, "serviceid", sink.Serviceid,
			"provincename", sink.Provincename, "response", responseStr, "tconsignid",
			sink.Tconsignid, "key", sink.Key, "cnum", sink.Cnum, "json", jsonstr); err != nil {
			return err
		}
	}
	return nil
}
func (self *Cache) UpdateSetByKey(key string, user *SpUser) (rnum int64, err error) {
	rc, err := self.pool.Get()
	if err != nil {
		return 0, err
	}
	defer func() {
		self.pool.Close(rc)
	}()
	return rc.SAdd(key, user.SpServiceRule.Terminal)
}

func (self *Cache) UpdateMoUsers(user *SpUser) (rnum int64, rbool bool, err error) {
	key := time.Now().Format("20060102") + ":mousers:consignid:" + strconv.Itoa(user.SpConsign.Consignid) +
		":serviceid:" + strconv.Itoa(user.SpService.Serviceid) +
		":provinceid:" + user.SpMsisdn.Provinceid
	if rnum, err := self.UpdateSetByKey(key, user); err != nil {
		return rnum, false, err
	} else {
		rc, err := self.pool.Get()
		if err != nil {
			return rnum, false, err
		}
		defer func() {
			self.pool.Close(rc)
		}()
		rb, err := rc.Expire(key, 60*60*72)
		return rnum, rb, err
	}
}

func (self *Cache) UpdateMtUsers(user *SpUser) (rnum int64, rbool bool, err error) {
	key := time.Now().Format("20060102") + ":mtusers:consignid:" + strconv.Itoa(user.SpConsign.Consignid) +
		":serviceid:" + strconv.Itoa(user.SpService.Serviceid) +
		":provinceid:" + user.SpMsisdn.Provinceid
	if rnum, err := self.UpdateSetByKey(key, user); err != nil {
		return rnum, false, err
	} else {
		rc, err := self.pool.Get()
		if err != nil {
			return rnum, false, err
		}
		defer func() {
			self.pool.Close(rc)
		}()
		rb, err := rc.Expire(key, 60*60*72)
		return rnum, rb, err
	}
}

func (self *Cache) GetSetsBulkByKey(key string) (rnum int64, err error) {
	rc, err := self.pool.Get()
	if err != nil {
		return 0, err
	}
	defer func() {
		self.pool.Close(rc)
	}()
	return rc.SCard(key)
}

func (self *Cache) GetTodayMoUsers(user *SpUser) (rnum int64, err error) {
	key := time.Now().Format("20060102") + ":mousers:consignid:" + strconv.Itoa(user.SpConsign.Consignid) +
		":serviceid:" + strconv.Itoa(user.SpService.Serviceid) +
		":provinceid:" + user.SpMsisdn.Provinceid
	return self.GetSetsBulkByKey(key)
}

func (self *Cache) GetYdayMoUsers(user *SpUser) (rnum int64, err error) {
	day, _ := time.ParseDuration("-24h")
	key := time.Now().Add(day).Format("20060102") + ":mousers:consignid:" + strconv.Itoa(user.SpConsign.Consignid) +
		":serviceid:" + strconv.Itoa(user.SpService.Serviceid) +
		":provinceid:" + user.SpMsisdn.Provinceid
	return self.GetSetsBulkByKey(key)
}

func (self *Cache) GetTodayMtUsers(user *SpUser) (rnum int64, err error) {
	key := time.Now().Format("20060102") + ":mtusers:consignid:" + strconv.Itoa(user.SpConsign.Consignid) +
		":serviceid:" + strconv.Itoa(user.SpService.Serviceid) +
		":provinceid:" + user.SpMsisdn.Provinceid
	return self.GetSetsBulkByKey(key)
}

func (self *Cache) GetYdayMtUsers(user *SpUser) (rnum int64, err error) {
	day, _ := time.ParseDuration("-24h")
	key := time.Now().Add(day).Format("20060102") + ":mtusers:consignid:" + strconv.Itoa(user.SpConsign.Consignid) +
		":serviceid:" + strconv.Itoa(user.SpService.Serviceid) +
		":provinceid:" + user.SpMsisdn.Provinceid
	return self.GetSetsBulkByKey(key)
}

func (self *Cache) IncrHashFieldByKey(key string, field string, num int) (rnum int64, err error) {
	rc, err := self.pool.Get()
	if err != nil {
		return 0, err
	}
	defer func() {
		self.pool.Close(rc)
	}()
	return rc.HIncrBy(key, field, num)
}

func (self *Cache) UpdateSinkCnum(user *SpUser, cnum int) (spSink *SpSink, rnum int64, err error) {
	keyPrefix := "sp_sink:consignid:" + strconv.Itoa(user.SpConsign.Consignid)
	keyServiceid := ":serviceid:" + strconv.Itoa(user.SpService.Serviceid)
	keyProvince := ":provincename:" + user.SpMsisdn.Provincename
	sink := &SpSink{}
	if err := self.GetCache(keyPrefix+keyServiceid+keyProvince, CACHE_TYPE_HASH, sink); err == nil {
		cnum, err := self.IncrHashFieldByKey(keyPrefix+keyServiceid+keyProvince, "cnum", cnum)
		return sink, cnum, err
	} else if err := self.GetCache(keyPrefix+keyServiceid, CACHE_TYPE_HASH, sink); err == nil {
		cnum, err := self.IncrHashFieldByKey(keyPrefix+keyServiceid, "cnum", cnum)
		return sink, cnum, err
	} else if err := self.GetCache(keyPrefix+keyProvince, CACHE_TYPE_HASH, sink); err == nil {
		cnum, err := self.IncrHashFieldByKey(keyPrefix+keyProvince, "cnum", cnum)
		return sink, cnum, err
	} else if err := self.GetCache(keyPrefix, CACHE_TYPE_HASH, sink); err == nil {
		cnum, err := self.IncrHashFieldByKey(keyPrefix, "cnum", cnum)
		return sink, cnum, err
	} else {
		return nil, 0, errors.New("update sink cnum fails >> not found key")
	}
}

func (self *Cache) IncrSinkCnum(user *SpUser) (spSink *SpSink, rnum int64, err error) {
	return self.UpdateSinkCnum(user, 1)
}

func (self *Cache) UpdateHashByKey(key string, filed string, v interface{}) (bool, error) {
	rc, err := self.pool.Get()
	if err != nil {
		return false, err
	}
	defer func() {
		self.pool.Close(rc)
	}()
	return rc.HSet(key, filed, v)
}
func (self *Cache) ResetSinkCnum(user *SpUser) (spSink *SpSink, ret bool, err error) {
	keyPrefix := "sp_sink:consignid:" + strconv.Itoa(user.SpConsign.Consignid)
	keyServiceid := ":serviceid:" + strconv.Itoa(user.SpService.Serviceid)
	keyProvince := ":provincename:" + user.SpMsisdn.Provincename

	sink := &SpSink{}
	if err := self.GetCache(keyPrefix+keyServiceid+keyProvince, CACHE_TYPE_HASH, sink); err == nil {
		r, err := self.UpdateHashByKey(keyPrefix+keyServiceid+keyProvince, "cnum", 1)
		return sink, r, err
	} else if err := self.GetCache(keyPrefix+keyServiceid, CACHE_TYPE_HASH, sink); err == nil {
		r, err := self.UpdateHashByKey(keyPrefix+keyServiceid, "cnum", 1)
		return sink, r, err
	} else if err := self.GetCache(keyPrefix+keyProvince, CACHE_TYPE_HASH, sink); err == nil {
		r, err := self.UpdateHashByKey(keyPrefix+keyProvince, "cnum", 1)
		return sink, r, err
	} else if err := self.GetCache(keyPrefix, CACHE_TYPE_HASH, sink); err == nil {
		r, err := self.UpdateHashByKey(keyPrefix, "cnum", 1)
		return sink, r, err
	} else {
		return nil, false, errors.New("reset sink cnum fails >> not found key")
	}
}
func (self *Cache) SetMo(user *SpUser) error {
	var end string
	vr := reflect.ValueOf(&user.SpServiceRule).Elem()
	for _, v := range user.SpService.Referrule.Key {
		end += v + ":" + vr.FieldByName(strings.Title(v)).String() + ":"
	}
	keyPrefix := "sp_mo:" + end[:len(end)-1]

	rc, err := self.pool.Get()
	if err != nil {
		return err
	}
	defer func() {
		self.pool.Close(rc)
	}()
	jsonstr, err := json.Marshal(user)
	if err != nil {
		return err
	}
	if _, err := rc.HSet(keyPrefix, "json", jsonstr); err != nil {
		return err
	}
	if _, err := rc.Expire(keyPrefix+end, 60*60*72); err != nil {
		return err
	}
	return nil
}

func (self *Cache) RbacNodeToMap() (map[string]*SpStatNode, error) {
	stmtOut, err := self.db.Prepare("SELECT id, name, node FROM sp_node_privilege")
	defer stmtOut.Close()
	if err != nil {
		return nil, err
	}
	result, err := stmtOut.Query()
	if err != nil {
		return nil, err
	}
	nMap := make(map[string]*SpStatNode)
	for result.Next() {
		node := &SpStatNode{}

		if err := result.Scan(&node.Id, &node.Name, &node.Node); err != nil {
			return nil, err
		} else {
			nMap[node.Node] = node
		}
	}
	return nMap, nil
}

func (self *Cache) RbacMenuToSlice() ([]*SpStatMenu, error) {
	stmtOut, err := self.db.Prepare("SELECT id, title, name FROM sp_menu_template")
	defer stmtOut.Close()
	if err != nil {
		return nil, err
	}
	result, err := stmtOut.Query()
	if err != nil {
		return nil, err
	}
	var ms []*SpStatMenu
	for result.Next() {
		menu := &SpStatMenu{}

		if err := result.Scan(&menu.Id, &menu.Title, &menu.Name); err != nil {
			return nil, err
		} else {
			ms = append(ms, menu)
		}
	}
	return ms, nil
}

func (self *Cache) GetCache(key string, cacheType int, response interface{}) error {
	var (
		jsonstr string
		err     error
		rc      *redis.Client
	)
	rc, err = self.pool.Get()
	defer func() {
		self.pool.Close(rc)
	}()

	if err != nil {
		return err
	}

	switch cacheType {
	case CACHE_TYPE_HASH:
		if jsonstr, err = rc.HGet(key, "json"); err != nil {
			return err
		}
	default:
		return errors.New("Invalid cache type")
	}

	if err = json.Unmarshal([]byte(jsonstr), response); err != nil {
		return err
	}
	return nil
}
