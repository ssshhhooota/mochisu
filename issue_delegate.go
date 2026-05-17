package main

import (
	"io"

	"charm.land/bubbles/v2/list"
)

type issueDelegate struct {
	list.DefaultDelegate
}

func newIssueDelegate() issueDelegate {
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.SetSpacing(0)
	return issueDelegate{DefaultDelegate: d}
}

func (d issueDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	i, ok := item.(issue)
	if !ok {
		d.DefaultDelegate.Render(w, m, index, item)
		return
	}

	tmp := d.DefaultDelegate
	switch i.State {
	case "OPEN":
		tmp.Styles.NormalTitle = tmp.Styles.NormalTitle.Foreground(colorStateOpen)
		tmp.Styles.SelectedTitle = tmp.Styles.SelectedTitle.Foreground(colorStateOpen)
	case "CLOSED":
		tmp.Styles.NormalTitle = tmp.Styles.NormalTitle.Foreground(colorStateClosed)
		tmp.Styles.SelectedTitle = tmp.Styles.SelectedTitle.Foreground(colorStateClosed)
	}
	tmp.Render(w, m, index, item)
}
