package profileui

import tea "github.com/charmbracelet/bubbletea"

type ProfileSelectedMsg struct {
	Profile Profile
}

type ProfileSelectCancelledMsg struct{}

type CreateProfileMsg struct{}

type CreateProfileFinishedMsg struct {
	root     string
	filename string
	err      error
}

type DeleteProfileMsg struct {
	Profile Profile
}

type CopyProfileMsg struct {
	Profile Profile
}

type RenameProfileMsg struct {
	Profile Profile
}

type EditProfileMsg struct {
	Profile Profile
}

type EditProfileFinishedMsg struct {
	Profile Profile
	err     error
}

type PreviewProfileMsg struct {
	Profile Profile
}

type ProfilesChangedMsg struct{}

func CreateChangeCmd() tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		return ProfilesChangedMsg{}
	})
}
