package tutorials

import (
	"image/color"

	"github.com/ravenlab/fyne"
	"github.com/ravenlab/fyne/canvas"
	"github.com/ravenlab/fyne/container"
	"github.com/ravenlab/fyne/layout"
	"github.com/ravenlab/fyne/theme"
	"github.com/ravenlab/fyne/widget"
)

type browser struct {
	current int

	name *widget.Select
	icon *widget.Icon
}

func (b *browser) setIcon(index int) {
	if index < 0 || index > len(icons)-1 {
		return
	}
	b.current = index

	b.name.SetSelected(icons[index].name)
	b.icon.SetResource(icons[index].icon)
}

// iconScreen loads a panel that shows the various icons available in Fyne
func iconScreen(_ fyne.Window) fyne.CanvasObject {
	b := &browser{}

	prev := widget.NewButtonWithIcon("", theme.NavigateBackIcon(), func() {
		b.setIcon(b.current - 1)
	})
	next := widget.NewButtonWithIcon("", theme.NavigateNextIcon(), func() {
		b.setIcon(b.current + 1)
	})
	b.name = widget.NewSelect(iconList(), func(name string) {
		for i, icon := range icons {
			if icon.name == name {
				if b.current != i {
					b.setIcon(i)
				}
				break
			}
		}
	})
	b.name.SetSelected(icons[b.current].name)
	buttons := container.NewHBox(prev, next)
	bar := container.NewBorder(nil, nil, buttons, nil, b.name)

	background := canvas.NewRasterWithPixels(checkerPattern)
	background.SetMinSize(fyne.NewSize(280, 280))
	b.icon = widget.NewIcon(icons[b.current].icon)

	return fyne.NewContainerWithLayout(layout.NewBorderLayout(
		bar, nil, nil, nil), bar, background, b.icon)
}

func checkerPattern(x, y, _, _ int) color.Color {
	x /= 20
	y /= 20

	if x%2 == y%2 {
		return theme.BackgroundColor()
	}

	return theme.ShadowColor()
}

func iconList() []string {
	var ret []string
	for _, icon := range icons {
		ret = append(ret, icon.name)
	}

	return ret
}

var icons = []struct {
	name string
	icon fyne.Resource
}{
	{"CancelIcon", theme.CancelIcon()},
	{"ConfirmIcon", theme.ConfirmIcon()},
	{"DeleteIcon", theme.DeleteIcon()},
	{"SearchIcon", theme.SearchIcon()},
	{"SearchReplaceIcon", theme.SearchReplaceIcon()},

	{"CheckButtonIcon", theme.CheckButtonIcon()},
	{"CheckButtonCheckedIcon", theme.CheckButtonCheckedIcon()},
	{"RadioButtonIcon", theme.RadioButtonIcon()},
	{"RadioButtonCheckedIcon", theme.RadioButtonCheckedIcon()},

	{"ContentAddIcon", theme.ContentAddIcon()},
	{"ContentRemoveIcon", theme.ContentRemoveIcon()},
	{"ContentClearIcon", theme.ContentClearIcon()},
	{"ContentCutIcon", theme.ContentCutIcon()},
	{"ContentCopyIcon", theme.ContentCopyIcon()},
	{"ContentPasteIcon", theme.ContentPasteIcon()},
	{"ContentRedoIcon", theme.ContentRedoIcon()},
	{"ContentUndoIcon", theme.ContentUndoIcon()},

	{"InfoIcon", theme.InfoIcon()},
	{"QuestionIcon", theme.QuestionIcon()},
	{"WarningIcon", theme.WarningIcon()},

	{"DocumentCreateIcon", theme.DocumentCreateIcon()},
	{"DocumentPrintIcon", theme.DocumentPrintIcon()},
	{"DocumentSaveIcon", theme.DocumentSaveIcon()},

	{"FileIcon", theme.FileIcon()},
	{"FileApplicationIcon", theme.FileApplicationIcon()},
	{"FileAudioIcon", theme.FileAudioIcon()},
	{"FileImageIcon", theme.FileImageIcon()},
	{"FileTextIcon", theme.FileTextIcon()},
	{"FileVideoIcon", theme.FileVideoIcon()},
	{"FolderIcon", theme.FolderIcon()},
	{"FolderNewIcon", theme.FolderNewIcon()},
	{"FolderOpenIcon", theme.FolderOpenIcon()},
	{"HomeIcon", theme.HomeIcon()},
	{"HelpIcon", theme.HelpIcon()},
	{"SettingsIcon", theme.SettingsIcon()},

	{"ViewFullScreenIcon", theme.ViewFullScreenIcon()},
	{"ViewRestoreIcon", theme.ViewRestoreIcon()},
	{"ViewRefreshIcon", theme.ViewRefreshIcon()},
	{"VisibilityIcon", theme.VisibilityIcon()},
	{"VisibilityOffIcon", theme.VisibilityOffIcon()},
	{"ZoomFitIcon", theme.ZoomFitIcon()},
	{"ZoomInIcon", theme.ZoomInIcon()},
	{"ZoomOutIcon", theme.ZoomOutIcon()},

	{"MoveDownIcon", theme.MoveDownIcon()},
	{"MoveUpIcon", theme.MoveUpIcon()},

	{"NavigateBackIcon", theme.NavigateBackIcon()},
	{"NavigateNextIcon", theme.NavigateNextIcon()},

	{"MenuDropDown", theme.MenuDropDownIcon()},
	{"MenuDropUp", theme.MenuDropUpIcon()},

	{"MailAttachmentIcon", theme.MailAttachmentIcon()},
	{"MailComposeIcon", theme.MailComposeIcon()},
	{"MailForwardIcon", theme.MailForwardIcon()},
	{"MailReplyIcon", theme.MailReplyIcon()},
	{"MailReplyAllIcon", theme.MailReplyAllIcon()},
	{"MailSendIcon", theme.MailSendIcon()},

	{"MediaFastForward", theme.MediaFastForwardIcon()},
	{"MediaFastRewind", theme.MediaFastRewindIcon()},
	{"MediaPause", theme.MediaPauseIcon()},
	{"MediaPlay", theme.MediaPlayIcon()},
	{"MediaRecord", theme.MediaRecordIcon()},
	{"MediaReplay", theme.MediaReplayIcon()},
	{"MediaSkipNext", theme.MediaSkipNextIcon()},
	{"MediaSkipPrevious", theme.MediaSkipPreviousIcon()},

	{"VolumeDown", theme.VolumeDownIcon()},
	{"VolumeMute", theme.VolumeMuteIcon()},
	{"VolumeUp", theme.VolumeUpIcon()},
}
