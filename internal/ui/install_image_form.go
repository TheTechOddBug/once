package ui

import (
	tea "charm.land/bubbletea/v2"
)

const (
	installImageField = iota
	installRegistryUsernameField
	installRegistryPasswordField
)

type InstallImageSubmitMsg struct {
	ImageRef         string
	RegistryUsername string
	RegistryPassword string
}
type InstallImageBackMsg struct{}

type InstallImageForm struct {
	form Form
}

func NewInstallImageForm() InstallImageForm {
	usernameField := NewTextField("(optional)")
	passwordField := NewTextField("(optional)")
	passwordField.SetEchoPassword()

	m := InstallImageForm{
		form: NewForm("Next",
			FormItem{
				Label:    "Image",
				Field:    NewTextField("user/repo:tag"),
				Required: true,
			},
			FormItem{Label: "Registry Username", Field: usernameField},
			FormItem{Label: "Registry Password", Field: passwordField},
		),
	}

	m.form.OnSubmit(func(f *Form) tea.Cmd {
		ref := f.TextField(installImageField).Value()
		if expanded, ok := expandAlias(ref); ok {
			ref = expanded
		}
		return func() tea.Msg {
			return InstallImageSubmitMsg{
				ImageRef:         ref,
				RegistryUsername: f.TextField(installRegistryUsernameField).Value(),
				RegistryPassword: f.TextField(installRegistryPasswordField).Value(),
			}
		}
	})
	m.form.OnCancel(func(f *Form) tea.Cmd {
		return func() tea.Msg { return InstallImageBackMsg{} }
	})

	return m
}

func (m InstallImageForm) Init() tea.Cmd {
	return m.form.Init()
}

func (m InstallImageForm) Update(msg tea.Msg) (InstallImageForm, tea.Cmd) {
	var cmd tea.Cmd
	m.form, cmd = m.form.Update(msg)
	return m, cmd
}

func (m InstallImageForm) View() string {
	return m.form.View()
}
