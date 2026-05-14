package main

import (
	"charm.land/bubbles/v2/list"
	"charm.land/bubbles/v2/viewport"

	"charm.land/huh/v2"
)

type focus int

const (
	focusList focus = iota
	focusPreview
)

type repo struct {
	NameWithOwner string `json:"nameWithOwner"`
}

type issue struct {
	Name   string `json:"title"`
	Body   string `json:"body"`
	Number int    `json:"number"`
	URL    string `json:"url"`
}

type model struct {
	repos         []repo
	selectedRepo  string
	selectedIssue issue
	form          *huh.Form
	screen        screen
	issues        list.Model
	preview       viewport.Model
	focus         focus
	width, height int
}

func (i issue) Title() string       { return i.Name }
func (i issue) Description() string { return "" }
func (i issue) FilterValue() string { return i.Name }
func newRepoForm(repos []repo) *huh.Form {
	options := make([]huh.Option[string], 0, len(repos))
	for _, r := range repos {
		options = append(options, huh.NewOption(r.NameWithOwner, r.NameWithOwner))
	}
	return huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Pick a repo.").
				Options(options...).
				Key("repo"),
		),
	)
}

func newModel(repos []repo) model {
	m := model{
		repos: repos,
	}
	m.form = newRepoForm(repos)
	d := list.NewDefaultDelegate()
	d.ShowDescription = false
	d.SetSpacing(0)
	l := list.New([]list.Item{}, d, 0, 0)
	l.SetStatusBarItemName("issues", "issue")
	l.SetShowHelp(true)
	l.SetShowTitle(true)
	l.Title = "issues"
	l.InfiniteScrolling = true
	l.SetShowPagination(false)

	m.issues = l

	vp := viewport.New()
	vp.SoftWrap = true
	m.preview = vp

	return m

}
