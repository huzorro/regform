package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/huzorro/spfactor/sexredis"
	"log"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type StoreMsg struct {
	db *sql.DB
}
type StoreMo struct {
	db *sql.DB
}
type StoreMt struct {
	db *sql.DB
}
type IpToRule struct {
	cache *Cache
}
type MatchArea struct {
	cache *Cache
}

type MatchWord struct {
	cache *Cache
}

type MatchMo struct {
	cache *Cache
}

type CacheMo struct {
	cache *Cache
}

type SinkReport struct {
	cache *Cache
}

type UpdateMoUsers struct {
	cache *Cache
}

type UpdateMtUsers struct {
	cache *Cache
}

type FinalMoStat struct {
	db    *sql.DB
	cache *Cache
}

type FinalMtStat struct {
	db    *sql.DB
	cache *Cache
}

type SinkMtStat struct {
	db    *sql.DB
	cache *Cache
}
type MoQueue struct {
	p *sexredis.RedisPool
}
type MtQueue struct {
	p *sexredis.RedisPool
}

func (self *StoreMsg) SProcess(msg *sexredis.Msg) {
	log.Printf("Store msg process start ... %s", msg)

	stmtIn, err := self.db.Prepare("INSERT INTO sp_receive_log (msg) VALUES(?)")
	defer stmtIn.Close()
	if err != nil {
		log.Printf("Store msg Error:%s", err.Error())
		msg.Err = err
		return
	}
	result, err := stmtIn.Exec(msg.Content)

	if err != nil {
		log.Printf("INSERT INTO sp_receive_log(msg) VALUES(%s) Error:%s", msg, err.Error())
		msg.Err = err
		return
	}
	ids, _ := result.LastInsertId()
	log.Printf("Store msg process end ... INSERT INTO sp_receive_log(msg) VALUES(%s) [%d]", msg, ids)
}

func (self *IpToRule) SProcess(msg *sexredis.Msg) {

	log.Printf("IpToRule process start ... %+v", msg)
	m := make(map[string][]string)
	if str, ok := msg.Content.(string); ok {
		if err := json.Unmarshal([]byte(str), &m); err != nil {
			log.Printf("ip to rule Error:%s", err.Error())
			msg.Err = err
			return
		}
	} else {
		log.Printf("%s is not string", msg.Content)
		msg.Err = errors.New("Msg is not string")
		return
	}
	rip := m["rip"][0]
	service := SpService{}
	if err := self.cache.GetCache("sp_service:serviceip:"+rip, CACHE_TYPE_HASH, &service); err != nil {
		log.Printf("not found in cache %s", rip)
		msg.Err = errors.New("not found in cache")
		return
	}

	serviceRule := SpServiceRule{}

	st := reflect.ValueOf(&serviceRule).Elem()
	t := reflect.TypeOf(&service.Servicerule).Elem()
	v := reflect.ValueOf(&service.Servicerule).Elem()

	for i := 0; i < t.NumField(); i++ {
		n := t.Field(i).Name
		v := v.FieldByName(n).String()
		if s, ok := m[v]; ok {
			st.FieldByName(n).SetString(s[0])
		}
	}
	user := SpUser{}

	user.SpService = service
	user.SpServiceRule = serviceRule
	msg.Content = user
	log.Printf("IpToRule process end ... %+v", msg)
}
func (self *MatchArea) SProcess(msg *sexredis.Msg) {
	log.Printf("Match area process start ... %+v", msg)
	var (
		user SpUser
		ok   bool
	)
	//msg type ok ?
	if user, ok = msg.Content.(SpUser); !ok {
		log.Printf("Msg type error %+v", msg)
		msg.Err = errors.New("Msg type error")
		return
	}

	msisdn := SpMsisdn{}
	//terminal
	rule := user.SpServiceRule
	if ru := []rune(rule.Terminal); rule.Terminal != "" && len(ru) > 7 && len(ru) < 13 {
		//unicode
		//		ru := []rune(rule.Terminal)
		//		if len(ru) < 7 || len(ru) > 13 {
		//			log.Printf("terminal is too short or long %s", ru)
		//			return
		//		}

		tp := string(ru[0:7])

		if err := self.cache.GetCache("sp_msisdn:prefix:"+tp, CACHE_TYPE_HASH, &msisdn); err != nil {
			log.Printf("not found in cache %s", tp)
		}
		//cityid
	} else if rule.Cityid != "" {
		cid, _ := fmt.Printf("%04s", rule.Cityid)
		if err := self.cache.GetCache("sp_msisdn:cityid:0313", CACHE_TYPE_HASH, &msisdn); err != nil {
			log.Printf("not found in cache %s", cid)
		}
		//provinceid
	} else if rule.Provinceid != "" {
		pid, _ := fmt.Printf("%04s", rule.Provinceid)
		if err := self.cache.GetCache("sp_msisdn:provinceid:0311", CACHE_TYPE_HASH, &msisdn); err != nil {
			log.Printf("not found in cache %s", pid)
		}
	} else {
		log.Printf("not found area info %+v", rule)
	}

	user.SpMsisdn = msisdn
	msg.Content = user
	log.Printf("Match area process end ... %+v", msg)
}

func (self *MatchWord) SProcess(msg *sexredis.Msg) {
	log.Printf("Match word process start ... %+v", msg)
	var (
		user SpUser
		ok   bool
	)
	//msg type ok ?
	if user, ok = msg.Content.(SpUser); !ok {
		log.Printf("Msg type error %s", user)
		msg.Err = errors.New("Msg type error")
		return
	}
	service := SpService{}
	//fuzzy match
	rule := user.SpServiceRule
	for i := len(rule.Serviceword); i > 0; i-- {
		if err := self.cache.GetCache("sp_service:spnum:"+rule.Spnum+":serviceword:"+rule.Serviceword[:i], CACHE_TYPE_HASH, &service); err == nil {
			break
		}
	}

	user.SpService = service
	cp := SpCp{}

	if err := self.cache.GetCache("sp_cp:serviceid:"+strconv.Itoa(service.Serviceid), CACHE_TYPE_HASH, &cp); err != nil {
		log.Printf("not found in cache serviceid:%d", service.Serviceid)
		msg.Err = errors.New("not found in cache")
		return
	}
	info := SpInfo{}
	if err := self.cache.GetCache("sp_info:spnum:"+service.Spnum, CACHE_TYPE_HASH, &info); err != nil {
		log.Printf("not found in cache spnum:%s", service.Spnum)
		msg.Err = errors.New("not found in cache")
		return
	}
	consign := SpConsign{}
	if err := self.cache.GetCache("sp_consign:consignid:"+strconv.Itoa(cp.Consignid), CACHE_TYPE_HASH, &consign); err != nil {
		log.Printf("not found in cache consignid:%d", cp.Consignid)
		msg.Err = errors.New("not found in cache")
		return
	}
	user.SpConsign = consign
	user.SpOconsign = consign //set original consign
	user.SpInfo = info
	msg.Content = user
	log.Printf("Match word process end ... %+v", msg)
}

func (self *CacheMo) SProcess(msg *sexredis.Msg) {
	log.Printf("Cache mo process start ... %+v", msg)
	var (
		user SpUser
		ok   bool
	)
	//msg type ok ?
	if user, ok = msg.Content.(SpUser); !ok {
		log.Printf("Msg type error %s", user)
		msg.Err = errors.New("Msg type error")
		return
	}
	//cache mo
	if err := self.cache.SetMo(&user); err != nil {
		log.Printf("cache mo fails %s", err)
		msg.Err = errors.New("cache mo fails")
		return
	}
	log.Printf("Cache mo process end ... %+v", msg)
}

func (self *MatchMo) SProcess(msg *sexredis.Msg) {
	log.Printf("Match mo process start ... %+v", msg)
	var (
		user SpUser
		ok   bool
	)
	//msg type ok ?
	if user, ok = msg.Content.(SpUser); !ok {
		log.Printf("Msg type error %s", user)
		msg.Err = errors.New("Msg type error")
		return
	}
	//match mo
	//	jsonStr, err := json.Marshal(user.SpServiceRule)
	//	if err != nil {
	//		log.Printf("Json Marshal error %s", err)
	//		msg.Err = errors.New("Json Marshal error")
	//		return
	//	}
	//	rMap := make(map[string]string)

	//	if err := json.Unmarshal(jsonStr, rMap); err != nil {
	//		log.Printf("Json Unmarshal error %s", err)
	//		msg.Err = errors.New("Json Unmarshal error")
	//		return
	//	}

	//	//	end := strings.Join(user.SpService.Referrule.Key, ":")
	//	var end string
	//	for _, v := range user.SpService.Referrule.Key {
	//		end += v + ":" + rMap[v] + ":"
	//	}
	//	keyPrefix := "sp_mo:" + end[:len(end)-1]
	var end string
	vr := reflect.ValueOf(&user.SpServiceRule).Elem()
	for _, v := range user.SpService.Referrule.Key {
		end += v + ":" + vr.FieldByName(strings.Title(v)).String() + ":"
	}
	keyPrefix := "sp_mo:" + end[:len(end)-1]

	userMo := SpUser{}
	if err := self.cache.GetCache(keyPrefix, CACHE_TYPE_HASH, &userMo); err != nil {
		log.Printf("not found in cache %s", keyPrefix)
		msg.Err = errors.New("not found in cache")
		return
	}
	if user.SpServiceRule.Statusid == user.SpService.Referrule.Statusid {
		userMo.SpServiceRule.Statusid = STATUS_ID_SUCCESS
	} else {
		userMo.SpServiceRule.Statusid = user.SpServiceRule.Statusid
	}
	userMo.SpServiceRule.Statuscode = user.SpServiceRule.Statuscode
	msg.Content = userMo
	log.Printf("Match mo process end ... %+v", msg)
}

func (self *SinkReport) SProcess(msg *sexredis.Msg) {
	log.Printf("Sink report process start ... %+v", msg)
	var (
		user SpUser
		ok   bool
	)
	//msg type ok ?
	if user, ok = msg.Content.(SpUser); !ok {
		log.Printf("Msg type error %s", user)
		msg.Err = errors.New("Msg type error")
		return
	}
	if user.SpServiceRule.Statusid != STATUS_ID_SUCCESS {
		log.Printf("Sink report process end ... %+v", msg)
		return
	}

	sink, rnum, err := self.cache.IncrSinkCnum(&user)
	if err != nil {
		log.Printf("Incr sink cnum fails")
		return
	}
	//exec sink
	if rnum > int64(sink.Bnum-sink.Slice) && rnum <= int64(sink.Bnum) {
		consign := SpConsign{}
		if err := self.cache.GetCache("sp_consign:consignid:"+strconv.Itoa(sink.Tconsignid), CACHE_TYPE_HASH, &consign); err != nil {
			log.Printf("not found in cache %d", user.SpConsign.Consignid)
			return
		}
		user.SpConsign = consign
		//set random status id
		r := sink.Response.Response
		if len(r) > 0 {
			user.SpServiceRule.Statusid = r[rand.Intn(len(r))]
		}
		//reset counter
	} else if rnum > int64(sink.Bnum) {
		log.Println(rnum, sink.Bnum, sink.Slice)
		keyPrefix := "sp_sink:consignid:" + strconv.Itoa(user.SpConsign.Consignid)
		keyServiceid := ":serviceid:" + strconv.Itoa(user.SpService.Serviceid)
		keyProvince := ":provincename:" + user.SpMsisdn.Provincename

		log.Println(keyPrefix + keyServiceid + keyProvince)
		if _, _, err := self.cache.ResetSinkCnum(&user); err != nil {
			log.Printf("rest sink cnum fails %s", err)
		}
	}
	msg.Content = user
	log.Printf("Sink report process end ... %+v", msg)
}
func (self *UpdateMoUsers) SProcess(msg *sexredis.Msg) {
	log.Printf("update mo users process start ... %+v", msg)
	var (
		user SpUser
		ok   bool
	)
	//msg type ok ?
	if user, ok = msg.Content.(SpUser); !ok {
		log.Printf("Msg type error %s", user)
		msg.Err = errors.New("Msg type error")
		return
	}

	if _, _, err := self.cache.UpdateMoUsers(&user); err != nil {
		log.Printf("update mo users fails %s", err)
		return
	}
	log.Printf("update mo users process end ... %+v", msg)
}

func (self *UpdateMtUsers) SProcess(msg *sexredis.Msg) {
	log.Printf("update mt users process start ... %+v", msg)
	var (
		user SpUser
		ok   bool
	)
	//msg type ok ?
	if user, ok = msg.Content.(SpUser); !ok {
		log.Printf("Msg type error %s", user)
		msg.Err = errors.New("Msg type error")
		return
	}
	if user.SpServiceRule.Statusid != STATUS_ID_SUCCESS {
		log.Printf("update mt users process end ... %+v", msg)
		return
	}
	if _, _, err := self.cache.UpdateMtUsers(&user); err != nil {
		log.Printf("update mt users fails %s", err)
		return
	}
	log.Printf("update mt users process end ... %+v", msg)
}

func (self *StoreMo) SProcess(msg *sexredis.Msg) {
	log.Printf("store mo process start ... %+v", msg)
	var (
		user SpUser
		ok   bool
	)
	//msg type ok ?
	if user, ok = msg.Content.(SpUser); !ok {
		log.Printf("Msg type error %s", user)
		msg.Err = errors.New("Msg type error")
		return
	}
	stmtOut, err := self.db.Prepare("INSERT INTO sp_mo_log (linkid, spnum, spname, msg, serviceword, servicename, " +
		"servicefee, servicetype, terminal, consignid, consignname, provinceid, provincename, " +
		"cityid, cityname, statusid, statuscode) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	defer stmtOut.Close()
	if err != nil {
		log.Printf("Error:%s", err.Error())
		msg.Err = err
		return
	}
	result, err := stmtOut.Exec(user.SpServiceRule.Linkid, user.SpInfo.Spnum, user.SpInfo.Spname, user.SpServiceRule.Serviceword, user.SpService.Serviceword,
		user.SpService.Servicename, user.SpService.Servicefee, strconv.Itoa(user.SpService.Servicetype),
		user.SpServiceRule.Terminal, user.SpConsign.Consignid, user.SpConsign.Consignname,
		user.SpMsisdn.Provinceid, user.SpMsisdn.Provincename,
		user.SpMsisdn.Cityid, user.SpMsisdn.Cityname,
		user.SpServiceRule.Statusid, user.SpServiceRule.Statuscode)

	if err != nil {
		log.Printf("INSERT INTO sp_mo_log Error:%s", msg, err.Error())
		msg.Err = err
		return
	}
	ids, _ := result.LastInsertId()
	log.Printf("Store mo process end ... INSERT INTO sp_mo_log %+v [%d]", msg, ids)
}

func (self *StoreMt) SProcess(msg *sexredis.Msg) {
	log.Printf("store mt process start ... %+v", msg)
	var (
		user SpUser
		ok   bool
	)
	//msg type ok ?
	if user, ok = msg.Content.(SpUser); !ok {
		log.Printf("Msg type error %s", user)
		msg.Err = errors.New("Msg type error")
		return
	}

	stmtOut, err := self.db.Prepare("INSERT INTO	sp_mt_log(linkid, spnum, spname, msg, serviceword, servicename, " +
		"servicefee,servicetype, terminal, expendtime, timeline, consignid, consignname, " +
		"provinceid, provincename, cityid, cityname, statusid, statuscode) " +
		"VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")

	defer stmtOut.Close()
	if err != nil {
		log.Printf("Error:%s", err.Error())
		msg.Err = err
		return
	}
	//insert original mt if sink  then statusid is random response
	result, err := stmtOut.Exec(user.SpServiceRule.Linkid, user.SpInfo.Spnum, user.SpInfo.Spname,
		user.SpServiceRule.Serviceword, user.SpServiceRule.Serviceword, user.SpService.Servicename,
		user.SpService.Servicefee, strconv.Itoa(user.SpService.Servicetype),
		user.SpServiceRule.Terminal, user.SpServiceRule.Expendtime, user.SpServiceRule.Timeline, user.SpOconsign.Consignid, user.SpOconsign.Consignname,
		user.SpMsisdn.Provinceid, user.SpMsisdn.Provincename,
		user.SpMsisdn.Cityid, user.SpMsisdn.Cityname,
		user.SpServiceRule.Statusid, user.SpServiceRule.Statuscode)

	if err != nil {
		log.Printf("INSERT INTO sp_mt_log Error:%s", msg, err.Error())
		msg.Err = err
		return
	}
	ids, _ := result.LastInsertId()
	log.Printf("Store mt process end ... INSERT INTO sp_mt_log %+v [%d]", msg, ids)

	if user.SpConsign.Consignid == user.SpOconsign.Consignid {
		return
	}

	//sink mt store
	result, err = stmtOut.Exec(user.SpServiceRule.Linkid, user.SpInfo.Spnum, user.SpInfo.Spname,
		user.SpServiceRule.Serviceword, user.SpServiceRule.Serviceword, user.SpService.Servicename,
		user.SpService.Servicefee, strconv.Itoa(user.SpService.Servicetype),
		user.SpServiceRule.Terminal, user.SpServiceRule.Expendtime, user.SpServiceRule.Timeline, user.SpConsign.Consignid, user.SpConsign.Consignname,
		user.SpMsisdn.Provinceid, user.SpMsisdn.Provincename,
		user.SpMsisdn.Cityid, user.SpMsisdn.Cityname,
		STATUS_ID_SUCCESS, user.SpServiceRule.Statuscode)

	if err != nil {
		log.Printf("INSERT INTO sp_mt_log Error:%s", msg, err.Error())
		msg.Err = err
		return
	}
	ids, _ = result.LastInsertId()
	log.Printf("Store sink mt process end ... INSERT INTO sp_mt_log %+v [%d]", msg, ids)
}

func (self *FinalMoStat) SProcess(msg *sexredis.Msg) {
	log.Printf("final mo stat process start ... %+v", msg)
	var (
		user SpUser
		ok   bool
	)
	//msg type ok ?
	if user, ok = msg.Content.(SpUser); !ok {
		log.Printf("Msg type error %s", user)
		msg.Err = errors.New("Msg type error")
		return
	}
	stmtOut, err := self.db.Prepare("INSERT INTO sp_final_stat (spnum, spname, consignid, consignname," +
		"serviceid, servicename, provinceid, provincename, monums, mousers, day) " +
		"VALUES (?,?, ?, ?, ?, ?, ?, ?, ?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE monums = monums + 1, mousers = VALUES(mousers);")
	defer stmtOut.Close()
	if err != nil {
		log.Printf("Error:%s", err.Error())
		msg.Err = err
		return
	}
	mus, err := self.cache.GetTodayMoUsers(&user)
	result, err := stmtOut.Exec(user.SpInfo.Spnum, user.SpInfo.Spname, user.SpConsign.Consignid, user.SpConsign.Consignname,
		user.SpService.Serviceid, user.SpService.Servicename,
		user.SpMsisdn.Provinceid, user.SpMsisdn.Provincename, 1, mus, time.Now().Format("2006-01-02"))
	if err != nil {
		log.Printf("INSERT INTO sp_final_stat ON DUPLICATE KEY UPDATE Error:%s", msg, err.Error())
		msg.Err = err
		return
	}
	ids, _ := result.LastInsertId()
	log.Printf("final mo stat process end ... INSERT INTO sp_final_stat ON DUPLICATE KEY UPDATE %+v [%d]", msg, ids)
}

func (self *FinalMtStat) SProcess(msg *sexredis.Msg) {
	log.Printf("final mt stat process start ... %+v", msg)
	var (
		user SpUser
		ok   bool
		fee  int
		mtN  int
	)
	//msg type ok ?
	if user, ok = msg.Content.(SpUser); !ok {
		log.Printf("Msg type error %s", user)
		msg.Err = errors.New("Msg type error")
		return
	}
	stmtOut, err := self.db.Prepare("INSERT INTO sp_final_stat (spnum, spname, consignid, consignname," +
		"serviceid, servicename, provinceid, provincename, mtnums, mtusers, fee, day) " +
		"VALUES (?,?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE mtnums = mtnums + VALUES(mtnums), mtusers = VALUES(mtusers), fee = fee + VALUES(fee);")
	defer stmtOut.Close()
	if err != nil {
		log.Printf("Error:%s", err.Error())
		msg.Err = err
		return
	}
	mus, err := self.cache.GetTodayMtUsers(&user)
	if user.SpServiceRule.Statusid == STATUS_ID_SUCCESS || (user.SpConsign.Consignid != user.SpOconsign.Consignid) {
		fee = user.SpService.Servicefee
		mtN = 1
	}

	result, err := stmtOut.Exec(user.SpInfo.Spnum, user.SpInfo.Spname, user.SpConsign.Consignid, user.SpConsign.Consignname,
		user.SpService.Serviceid, user.SpService.Servicename,
		user.SpMsisdn.Provinceid, user.SpMsisdn.Provincename, mtN, mus, fee, time.Now().Format("2006-01-02"))
	if err != nil {
		log.Printf("INSERT INTO sp_final_stat ON DUPLICATE KEY UPDATE Error:%s", msg, err.Error())
		msg.Err = err
		return
	}
	ids, _ := result.LastInsertId()
	log.Printf("final mt stat process end ... INSERT INTO sp_final_stat ON DUPLICATE KEY UPDATE %+v [%d]", msg, ids)
}

func (self *SinkMtStat) SProcess(msg *sexredis.Msg) {
	log.Printf("sink mt stat process start ... %+v", msg)
	var (
		user SpUser
		ok   bool
	)
	//msg type ok ?
	if user, ok = msg.Content.(SpUser); !ok {
		log.Printf("Msg type error %s", user)
		msg.Err = errors.New("Msg type error")
		return
	}
	if user.SpConsign.Consignid == user.SpOconsign.Consignid {
		log.Printf("sink mt stat process end ... %+v", msg)
		return
	}
	stmtOut, err := self.db.Prepare("INSERT INTO sp_sink_stat (spnum, spname, consignid, consignname," +
		"serviceid, servicename, provinceid, provincename, mtnums, mtusers, fee, day) " +
		"VALUES (?,?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) " +
		"ON DUPLICATE KEY UPDATE mtnums = mtnums + VALUES(mtnums), mtusers = VALUES(mtusers), fee = fee + VALUES(fee);")
	defer stmtOut.Close()
	if err != nil {
		log.Printf("Error:%s", err.Error())
		msg.Err = err
		return
	}
	mus, err := self.cache.GetTodayMtUsers(&user)

	result, err := stmtOut.Exec(user.SpInfo.Spnum, user.SpInfo.Spname, user.SpOconsign.Consignid, user.SpOconsign.Consignname,
		user.SpService.Serviceid, user.SpService.Servicename,
		user.SpMsisdn.Provinceid, user.SpMsisdn.Provincename, 1, mus, user.SpService.Servicefee, time.Now().Format("2006-01-02"))
	if err != nil {
		log.Printf("INSERT INTO sp_sink_stat ON DUPLICATE KEY UPDATE Error:%s", msg, err.Error())
		msg.Err = err
		return
	}
	ids, _ := result.LastInsertId()
	log.Printf("sink mt stat process end ... INSERT INTO sp_sink_stat ON DUPLICATE KEY UPDATE %+v [%d]", msg, ids)
}

func (self *MoQueue) SProcess(msg *sexredis.Msg) {
	log.Printf("put mo in queue process start ... %+v", msg)
	var (
		user SpUser
		ok   bool
	)
	//msg type ok ?
	if user, ok = msg.Content.(SpUser); !ok {
		log.Printf("Msg type error %s", user)
		msg.Err = errors.New("Msg type error")
		return
	}
	rc, err := self.p.Get()
	if err != nil {
		log.Printf("get redis connection of redis pool fails %s", err)
		return
	}
	defer self.p.Close(rc)

	queue := sexredis.New()
	queue.SetRClient(MO_REST_QUEUE_NAME+":"+strconv.Itoa(user.SpConsign.Consignid), rc)

	if jsonStr, err := json.Marshal(user.SpServiceRule); err != nil {
		log.Printf("Marshal sp service rule of spuser fails %s", err)
		return
	} else {
		if _, err := queue.Put(jsonStr); err != nil {
			log.Printf("put mo into redis queue fails %s", err)
			return
		}
	}

	log.Printf("put mo in queue process end ... %+v", msg)
}

func (self *MtQueue) SProcess(msg *sexredis.Msg) {
	log.Printf("put mt in queue process start ... %+v", msg)
	var (
		user SpUser
		ok   bool
	)
	//msg type ok ?
	if user, ok = msg.Content.(SpUser); !ok {
		log.Printf("Msg type error %s", user)
		msg.Err = errors.New("Msg type error")
		return
	}
	rc, err := self.p.Get()
	if err != nil {
		log.Printf("get redis connection of redis pool fails %s", err)
		return
	}
	defer self.p.Close(rc)

	queue := sexredis.New()
	queue.SetRClient(MT_REST_QUEUE_NAME+":"+strconv.Itoa(user.SpOconsign.Consignid), rc)
	if jsonStr, err := json.Marshal(user.SpServiceRule); err != nil {
		log.Printf("Marshal sp service rule of spuser fails %s", err)
		return
	} else {
		if _, err := queue.Put(jsonStr); err != nil {
			log.Printf("put mt into redis queue fails %s", err)
		}
	}
	log.Printf("put mt in queue process end ... %+v", msg)

}
