package main

import (
	"github.com/astaxie/beego/logs"
	"github.com/easymesh/autoproxy"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"time"
)

var mainWindow *walk.MainWindow

var mainWindowWidth = 300
var mainWindowHeight = 200

func waitWindows()  {
	for  {
		if mainWindow != nil && mainWindow.Visible() {
			mainWindow.SetSize(walk.Size{
				mainWindowWidth,
				mainWindowHeight})
			break
		}
		time.Sleep(100*time.Millisecond)
	}
	NotifyInit()
}

func MainWindowsClose()  {
	if mainWindow != nil {
		mainWindow.Close()
		mainWindow = nil
	}
}

func statusUpdate()  {
	StatUpdate(StatGet())
}

func init()  {
	go func() {
		waitWindows()
		for  {
			statusUpdate()
			time.Sleep(time.Second)
		}
	}()
}

var isAuth *walk.RadioButton
var protocal  *walk.RadioButton

func mainWindows() {
	CapSignal(CloseWindows)
	cnt, err := MainWindow{
		Title:   "AutoProxy " + autoproxy.VersionGet(),
		Icon: ICON_Main,
		AssignTo: &mainWindow,
		MinSize: Size{mainWindowWidth, mainWindowHeight-1},
		Size: Size{mainWindowWidth, mainWindowHeight-1},
		Layout:  VBox{},
		MenuItems: MenuBarInit(),
		StatusBarItems: StatusBarInit(),
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: ConsoleWidget(),
			},
			Composite{
				Layout: Grid{Columns: 2},
				Children: ButtonWight(),
			},
		},
	}.Run()

	if err != nil {
		logs.Error(err.Error())
	} else {
		logs.Info("main windows exit %d", cnt)
	}

	if err:= recover();err != nil{
		logs.Error(err)
	}

	CloseWindows()
}

func CloseWindows()  {
	if ServerRunning() {
		ServerShutdown()
	}
	MainWindowsClose()
	NotifyExit()
}
