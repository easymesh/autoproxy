package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/easymesh/autoproxy/engin"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type RemoteItem struct {
	Name     string
	Address  string
	Protocal string
	Auth     bool
	User     string
	Password string
}

func TestUrlGet() string {
	url := DataStringValueGet("remotetest")
	if url == "" {
		url = "https://google.com"
	}
	return url
}

func TestUrlSet(url string)  {
	DataStringValueSet("remotetest", url)
}

var remoteCache []RemoteItem

func remoteGet() []RemoteItem {
	if remoteCache == nil {
		list := make([]RemoteItem, 0)
		value := DataStringValueGet("remotelist")
		if value != "" {
			err := json.Unmarshal([]byte(value), &list)
			if err != nil {
				logs.Error("json marshal fail",err.Error())
			}
		}
		remoteCache = list
	}
	return remoteCache
}

func remoteSync()  {
	value, err := json.Marshal(remoteCache)
	if err != nil {
		logs.Error("json marshal fail",err.Error())
	} else {
		DataStringValueSet("remotelist", string(value))
	}
}

func RemoteIndexSet(name string)  {
	err := DataStringValueSet("remoteIndex", name)
	if err != nil {
		logs.Error(err.Error())
	}
}

func RemoteIndexGet() int {
	name := DataStringValueGet("remoteIndex")
	list := remoteGet()
	for idx, v := range list {
		if v.Name == name {
			return idx
		}
	}
	return 0
}

func RemoteCurName() string {
	name := DataStringValueGet("remoteIndex")
	list := remoteGet()
	for _, v := range list {
		if v.Name == name {
			return v.Name
		}
	}
	return list[0].Name
}

func RemoteList() []RemoteItem {
	return remoteGet()
}

func RemoteOptions() []string {
	var output []string
	list := remoteGet()
	for _, v := range list {
		output = append(output, v.Name)
	}
	if len(output) == 0 {
		output = append(output, "")
	}
	return output
}

func RemoteFind(name string) RemoteItem {
	list := remoteGet()
	for _, v := range list {
		if v.Name == name {
			return v
		}
	}
	return RemoteItem{
		Name: name, Protocal: "HTTPS",
	}
}

func RemoteGet() RemoteItem {
	list := remoteGet()
	if len(list) > 0 {
		return list[0]
	}
	return RemoteItem{
		Protocal: "HTTPS",
	}
}

func RemoteDelete(name string)  {
	defer remoteSync()
	for i, v := range remoteCache {
		if v.Name == name {
			remoteCache = append(remoteCache[:i], remoteCache[i+1:]...)
			return
		}
	}
}

func RemoteUpdate(item RemoteItem) {
	defer remoteSync()
	for i, v := range remoteCache {
		if v.Name == item.Name {
			remoteCache[i] = item
			return
		}
	}
	remoteCache = append(remoteCache, item)
}

func ProtocalOptions() []string {
	return []string{
		"HTTP","HTTPS",
	}
}

var curRemoteItem RemoteItem

func TestEngin(testhttps string, item *RemoteItem) (time.Duration, error) {
	now := time.Now()
	if !engin.IsConnect(item.Address, 5) {
		return 0, fmt.Errorf("remote address connnect %s fail", item.Address)
	}

	urls, err := url.Parse(testhttps)
	if err != nil {
		logs.Error("%s raw url parse fail, %s", testhttps, err.Error())
		return 0, err
	}

	var auth *engin.AuthInfo
	if item.Auth {
		auth = &engin.AuthInfo{
			User: item.User,
			Token: item.Password,
		}
	}

	var tls bool
	if strings.ToLower(item.Protocal) == "https" {
		tls = true
	}

	forward, err := engin.NewHttpsProtcal(item.Address, 10, auth, tls)
	if err != nil {
		logs.Error("new remote http proxy fail, %s", err.Error())
		return 0, err
	}

	defer forward.Close()

	request, err := http.NewRequest("GET", testhttps, nil)
	if err != nil {
		logs.Error("%s raw url parse fail, %s", testhttps, err.Error())
		return 0, err
	}

	if strings.ToLower(urls.Scheme) == "https" {
		conn, err := forward.Https(engin.Address(urls), request)
		if err != nil {
			logs.Error("remote server %s forward %s fail, %s",
				item.Address, urls.RawPath, err.Error())
			return 0, err
		}
		conn.Close()
	} else {
		rsp, err := forward.Http(request)
		if err != nil {
			logs.Error("remote server %s forward %s fail, %s",
				item.Address, urls.RawPath, err.Error())
			return 0, err
		}
		rsp.Body.Close()
	}

	return time.Now().Sub(now), nil
}

var remoteDlg *walk.Dialog

func RemoteServer()  {
	var acceptPB, cancelPB *walk.PushButton

	var remote, protocal *walk.ComboBox
	var auth *walk.RadioButton
	var user, passwd, address, testurl *walk.LineEdit
	var testbut *walk.PushButton

	curRemoteItem = RemoteGet()

	updateHandler := func() {
		protocal.SetText(curRemoteItem.Protocal)
		address.SetText(curRemoteItem.Address)
		auth.SetChecked(curRemoteItem.Auth)
		user.SetEnabled(curRemoteItem.Auth)
		passwd.SetEnabled(curRemoteItem.Auth)
		user.SetText(curRemoteItem.User)
		passwd.SetText(curRemoteItem.Password)
	}

	_, err := Dialog{
		AssignTo: &remoteDlg,
		Title: LangValue("remoteproxy"),
		Icon: walk.IconShield(),
		DefaultButton: &acceptPB,
		CancelButton: &cancelPB,
		Size: Size{250, 300},
		MinSize: Size{250, 300},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: LangValue("remoteproxy") + ":",
					},
					ComboBox{
						AssignTo: &remote,
						Editable: true,
						CurrentIndex:  0,
						Model:         RemoteOptions(),
						OnBoundsChanged: func() {
							curRemoteItem = RemoteFind(remote.Text())
							updateHandler()
						},
						OnCurrentIndexChanged: func() {
							curRemoteItem = RemoteFind(remote.Text())
							updateHandler()
						},
						OnEditingFinished: func() {
							curRemoteItem = RemoteFind(remote.Text())
							updateHandler()
						},
					},

					Label{
						Text: LangValue("remoteaddress") + ":",
					},

					LineEdit{
						AssignTo: &address,
						Text: curRemoteItem.Address,
						OnEditingFinished: func() {
							curRemoteItem.Address = address.Text()
						},
					},

					Label{
						Text: LangValue("protocal") + ":",
					},
					ComboBox{
						AssignTo: &protocal,
						Model: ProtocalOptions(),
						Value: curRemoteItem.Protocal,
						OnCurrentIndexChanged: func() {
							curRemoteItem.Protocal = protocal.Text()
						},
					},

					Label{
						Text: LangValue("whetherauth") + ":",
					},
					RadioButton{
						AssignTo: &auth,
						OnBoundsChanged: func() {
							auth.SetChecked(curRemoteItem.Auth)
						},
						OnClicked: func() {
							auth.SetChecked(!curRemoteItem.Auth)
							curRemoteItem.Auth = !curRemoteItem.Auth

							user.SetEnabled(curRemoteItem.Auth)
							passwd.SetEnabled(curRemoteItem.Auth)
						},
					},

					Label{
						Text: LangValue("user") + ":",
					},

					LineEdit{
						AssignTo: &user,
						Text: curRemoteItem.User,
						Enabled: curRemoteItem.Auth,
						OnEditingFinished: func() {
							curRemoteItem.User = user.Text()
						},
					},

					Label{
						Text: LangValue("password") + ":",
					},

					LineEdit{
						AssignTo: &passwd,
						Text: curRemoteItem.Password,
						Enabled: curRemoteItem.Auth,
						OnEditingFinished: func() {
							curRemoteItem.Password = passwd.Text()
						},
					},

					PushButton{
						AssignTo: &testbut,
						Text: LangValue("test"),
						OnClicked: func() {
							go func() {
								testbut.SetEnabled(false)
								delay, err := TestEngin(testurl.Text(), &curRemoteItem)
								if err != nil {
									ErrorBoxAction(remoteDlg, err.Error())
								} else {
									info := fmt.Sprintf("%s, %s %dms",
										LangValue("testpass"),
										LangValue("delay"), delay/time.Millisecond )
									InfoBoxAction(remoteDlg, info)
								}
								testbut.SetEnabled(true)
							}()
						},
					},

					LineEdit{
						AssignTo: &testurl,
						Text: TestUrlGet(),
						OnEditingFinished: func() {
							TestUrlSet(testurl.Text())
						},
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					PushButton{
						AssignTo: &acceptPB,
						Text:     LangValue("save"),
						OnClicked: func() {
							if curRemoteItem.Auth {
								if curRemoteItem.User == "" || curRemoteItem.Password == "" {
									ErrorBoxAction(remoteDlg, LangValue("inputuserandpasswd"))
									return
								}
							}
							if curRemoteItem.Name == "" || curRemoteItem.Address == "" {
								ErrorBoxAction(remoteDlg, LangValue("inputnameandaddress"))
								return
							}
							RemoteUpdate(curRemoteItem)
							remoteDlg.Accept()
							ConsoleRemoteUpdate()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      LangValue("delete"),
						OnClicked: func() {
							RemoteDelete(curRemoteItem.Name)
							remote.SetModel(RemoteOptions())
							remote.SetCurrentIndex(RemoteIndexGet())
							ConsoleRemoteUpdate()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      LangValue("cancel"),
						OnClicked: func() {
							remoteDlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(mainWindow)

	if err != nil {
		logs.Error(err.Error())
	}
}