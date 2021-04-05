package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
)

var statusFlow *walk.StatusBarItem

func UpdateStatFlow(flow string)  {
	if statusFlow != nil {
		statusFlow.SetText(flow)
	}
}

func UpdateStatFlag(image *walk.Icon)  {
	if statusFlow != nil {
		statusFlow.SetIcon(image)
	}
}

func StatusBarInit() []StatusBarItem {
	return []StatusBarItem{
		{
			AssignTo: &statusFlow,
			Icon: ICON_Network_Disable,
			ToolTipText: LangValue("realtimeflow"),
			Width: 80,
		},
	}
}
