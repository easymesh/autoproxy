package pages

import (
	"fmt"
	"github.com/GoAdminGroup/go-admin/context"
	"github.com/GoAdminGroup/go-admin/modules/db"
	"github.com/GoAdminGroup/go-admin/modules/logger"
	form2 "github.com/GoAdminGroup/go-admin/plugins/admin/modules/form"
	"github.com/GoAdminGroup/go-admin/plugins/admin/modules/table"
	"github.com/GoAdminGroup/go-admin/template/types"
	"github.com/GoAdminGroup/go-admin/template/types/form"
	editType "github.com/GoAdminGroup/go-admin/template/types/table"
	"github.com/astaxie/beego/logs"
	"github.com/easymesh/autoproxy/console/models"
	"strings"
	"sync"
	"time"
)


func DomainTableGet(ctx *context.Context) (table.Table) {
	profile := table.NewDefaultTable(table.DefaultConfigWithDriver("sqlite"))

	info := profile.GetInfo()
	info.AddField("ID", "id", db.Int).FieldFilterable()
	info.AddField("Tag", "tag", db.Varchar).FieldFilterable()
	info.AddField("Enable", "enable", db.Integer).FieldDisplay(func(model types.FieldModel) interface{} {
		return model.Value
	}).FieldEditAble(editType.Switch).FieldEditOptions(types.FieldOptions{
		{Value: "1", Text: "1"},
		{Value: "0", Text: "0"},
	}).FieldFilterable(types.FilterType{FormType: form.SelectSingle}).FieldFilterOptions(types.FieldOptions{
		{Value: "1", Text: "1"},
		{Value: "0", Text: "0"},
	})
	info.AddField("Domains", "domains", db.Varchar).FieldFilterable()
	info.SetTable("domains").SetTitle("Domain").SetDescription("edit proxy domain table")

	addFrom := profile.GetForm()
	addFrom.AddField("Tag", "tag", db.Varchar, form.Text).
		FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
			if value.Value.Value() == "" {
				logger.Error("tag is null", value.Value.Value())
				return ""
			}
			return value.Value.Value()
		})
	addFrom.AddField("Domains", "domains", db.Varchar, form.TextArea).
		FieldPostFilterFn(func(value types.PostFieldModel) interface{} {
			if value.Value.Value() == "" {
				logger.Error("domains is null", value.Value.Value())
				return ""
			}
			return value.Value.Value()
		}).FieldDefault("").FieldValue("")
	addFrom.AddField("Enable", "enable", db.Tinyint, form.Number).FieldDefault("1").FieldHide()
	addFrom.SetPostValidator(func(values form2.Values) error {
		if values.IsSingleUpdatePost() {
			if !values.Has("enable") {
				return fmt.Errorf("account single only enable update")
			}
			return nil
		}
		if len(values.Get("tag")) < 3 {
			return fmt.Errorf("tag should more than 3 characters")
		}
		if len(values.Get("domains")) < 1 {
			return fmt.Errorf("domains should more than 1 characters")
		}
		return nil
	})
	addFrom.SetPostHook(func(values form2.Values) error {
		go func() {
			time.Sleep(time.Second)
			DomainInit()
		}()
		return nil
	})
	addFrom.SetTable("domains").SetTitle("Domain").SetDescription("edit proxy domain table")
	return profile
}

type DomainCtrl struct {
	sync.RWMutex
	cache  map[string]string
	domain []string
}

var forwardCtrl DomainCtrl

func DomainInit() {
	var domainList []string

	domains := models.DomainGet()
	for _, v := range domains {
		if v.Enable == 0 {
			continue
		}
		domainList = append(domainList, strings.Split(v.Domains, ";")...)
	}

	forwardCtrl.Lock()
	forwardCtrl.cache  = make(map[string]string, 1024)
	forwardCtrl.domain = domainList
	forwardCtrl.Unlock()

	logger.Info("domain cache reset success")
}

func domainGet(address string) string {
	domain := address
	idx := strings.Index(address, ":")
	if idx != -1 {
		domain = address[:idx]
	}
	return domain
}

func domainMatch(domain string, match string) bool {
	begin := strings.Index(match, "*")
	end := strings.Index(match[begin+1:], "*")
	if end != -1 {
		end += begin+1
	}
	if begin != -1 && end == -1 {
		// suffix match
		return strings.HasSuffix(domain, match[begin+1:])
	}
	if begin == -1 && end != -1 {
		// prefix match
		return strings.HasPrefix(domain, match[:end])
	}
	if begin == -1 && end == -1 {
		// full match
		if domain == match {
			return true
		} else {
			return false
		}
	}
	idx := strings.Index(domain, match[begin+1: end])
	if idx == -1 {
		return false
	}
	return true
}

// address: www.baidu.com:80 or www.baidu.com:443
func routeMatch(address string) string {
	domain := domainGet(address)
	for _, v := range forwardCtrl.domain {
		if domainMatch(domain, v) {
			forwardCtrl.cache[address] = v
			logs.Info("route address %s match to domain %s", address, v)
			return v
		}
	}
	logs.Info("route address %s no match", address)
	forwardCtrl.cache[address] = ""
	return ""
}

func DomainCheck(address string) bool {
	forwardCtrl.RLock()
	result, flag := forwardCtrl.cache[address]
	forwardCtrl.RUnlock()

	if flag == false {
		forwardCtrl.Lock()
		result = routeMatch(address)
		forwardCtrl.Unlock()
	}

	if result == "" {
		return false
	}

	return true
}