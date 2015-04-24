package main

import (
	"database/sql"
	"encoding/json"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/martini-contrib/sessions"
	"log"
	"net/http"
)

type SpStatUser struct {
	Id       int
	UserName string
	Password string
	Role     *SpStatRole
	Access   *SpStatAccess
}

type SpStatRole struct {
	Id        int
	Name      string
	Privilege int
	Menu      int
}

type SpStatAccess struct {
	Id    int
	Group []string
	Rule  int
}

type SpStatNode struct {
	Id   int
	Name string
	Node string
}

type SpStatMenu struct {
	Id    int
	Title string
	Name  string
}
type LoginStatus struct {
	Status int
	Text   string
}

type StatByServiceid struct {
	Serviceid   int
	Servicename string
	Data        StatData
}

type StatByConsignid struct {
	Consignid   int
	Consignname string
	Data        StatData
}
type StatByProvinceid struct {
	Provinceid   string
	Provincename string
	Data         StatData
}
type StatBySpnum struct {
	Spnum  string
	Spname string
	Data   StatData
}

type StatTable struct {
	Consign  []StatByConsignid
	Service  []StatByServiceid
	Province []StatByProvinceid
}
type StatData struct {
	Monums, Mousers, Mtnums, Mtusers, Fee int
}
type IFilter interface {
	Filter(params ...interface{})
}

type RBAC struct{}

const (
	SESSION_KEY_QUSER = "Quser"
	ERROR_PAGE_NAME   = "/500.html"
	LOGIN_PAGE_NAME   = "/login"
	NODE_LOGIN_PRI    = 1 << 0
	NODE_RLOGIN_PRI   = 1 << 6

	ROLE_GUESS_ID = 1 << 0
	//定义匿名用户访问权限
	ROLE_GUESS_PRI = 123
)
const (
	GROUP_PRI_ALL = 1 << iota
	GROUP_PRI_ALLOW
	GROUP_PRI_BAN
)

func (self *RBAC) Filter() martini.Handler {
	return func(r *http.Request, w http.ResponseWriter, log *log.Logger, db *sql.DB, session sessions.Session, nMap map[string]*SpStatNode) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		path := r.URL.Path
		user := SpStatUser{}

		value := session.Get(SESSION_KEY_QUSER)
		if value == nil {
			role := SpStatRole{}
			role.Id = ROLE_GUESS_ID
			role.Privilege = ROLE_GUESS_PRI
			user.Role = &role
		} else {

			if v, ok := value.([]byte); ok {
				json.Unmarshal(v, &user)
			}
		}
		if v, ok := nMap[path]; ok {
			if (user.Role.Privilege & v.Id) != v.Id {
				http.Redirect(w, r, LOGIN_PAGE_NAME, 301)
				return
			}
			return
		}
		log.Printf("access unauth page")
		http.Redirect(w, r, LOGIN_PAGE_NAME, 301)
	}
}
