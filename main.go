package main

import (
	"errors"
	"log"
	"time"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/storage"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

func showTime(clock *widget.Label) {
	formatted := time.Now().Format("03:04:05")
	clock.SetText(formatted)
}

func toolbarView() *fyne.Container {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentCreateIcon(), func() {
			log.Println("New document")
		}),
		widget.NewToolbarSeparator(),
		widget.NewToolbarAction(theme.ContentCutIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentCopyIcon(), func() {}),
		widget.NewToolbarAction(theme.ContentPasteIcon(), func() {}),
		widget.NewToolbarSpacer(),
		widget.NewToolbarAction(theme.HelpIcon(), func() {
			log.Println("Display help")
		}),
	)

	content := fyne.NewContainerWithLayout(layout.NewBorderLayout(toolbar, nil, nil, nil),
		toolbar, widget.NewLabel("Content"))
	return content
}
func fileOpened(f fyne.URIReadCloser) {
	if f == nil {
		log.Println("Cancelled")
		return
	}

	ext := f.URI().Extension()
	if ext == ".png" {
		//showImage(f)
	} else if ext == ".txt" {
		//showText(f)
	}
	err := f.Close()
	if err != nil {
		fyne.LogError("Failed to close stream", err)
	}
}

func fileSaved(f fyne.URIWriteCloser) {
	if f == nil {
		log.Println("Cancelled")
		return
	}

	log.Println("Save to...", f.URI())
}
func dialogView(win fyne.Window) *widget.Group {
	return widget.NewGroup("Dialogs",
		widget.NewButton("Info", func() {
			dialog.ShowInformation("Information", "You should know this thing...", win)
		}),
		widget.NewButton("Error", func() {
			err := errors.New("a dummy error message")
			dialog.ShowError(err, win)
		}),
		widget.NewButton("Confirm", func() {
			cnf := dialog.NewConfirm("Confirmation", "Are you enjoying this demo?", nil, win)
			cnf.SetDismissText("Nah")
			cnf.SetConfirmText("Oh Yes!")
			cnf.Show()
		}),
		widget.NewButton("Progress", func() {
			prog := dialog.NewProgress("MyProgress", "Nearly there...", win)

			go func() {
				num := 0.0
				for num < 1.0 {
					time.Sleep(50 * time.Millisecond)
					prog.SetValue(num)
					num += 0.01
				}

				prog.SetValue(1)
				prog.Hide()
			}()

			prog.Show()
		}),
		widget.NewButton("ProgressInfinite", func() {
			prog := dialog.NewProgressInfinite("MyProgress", "Closes after 5 seconds...", win)

			go func() {
				time.Sleep(time.Second * 5)
				prog.Hide()
			}()

			prog.Show()
		}),
		widget.NewButton("File Open With Filter (.txt or .png)", func() {
			fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
				if err == nil && reader == nil {
					return
				}
				if err != nil {
					dialog.ShowError(err, win)
					return
				}

				fileOpened(reader)
			}, win)
			fd.SetFilter(storage.NewExtensionFileFilter([]string{".png", ".txt"}))
			fd.Show()
		}),
		widget.NewButton("File Save", func() {
			dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
				if err != nil {
					dialog.ShowError(err, win)
					return
				}

				fileSaved(writer)
			}, win)
		}),
		widget.NewButton("Custom Dialog (Login Form)", func() {
			username := widget.NewEntry()
			password := widget.NewPasswordEntry()
			content := widget.NewForm(widget.NewFormItem("Username", username),
				widget.NewFormItem("Password", password))

			dialog.ShowCustomConfirm("Login...", "Log In", "Cancel", content, func(b bool) {
				if !b {
					return
				}

				log.Println("Please Authenticate", username.Text, password.Text)
			}, win)
		}),
	)
}

func progressView() *widget.Box {
	progress := widget.NewProgressBar()
	infinite := widget.NewProgressBarInfinite()

	go func() {
		for i := 0.0; i <= 1.0; i += 0.1 {
			time.Sleep(time.Millisecond * 250)
			progress.SetValue(i)
		}
	}()
	return widget.NewVBox(progress, infinite)
}

func formView(myWindow fyne.Window) *widget.Form {
	entry := widget.NewEntry()
	textArea := widget.NewMultiLineEntry()

	form := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{"Entry", entry}},
		OnSubmit: func() { // optional, handle form submission
			log.Println("Form submitted:", entry.Text)
			log.Println("multiline:", textArea.Text)
			myWindow.Close()
		},
	}

	// we can also append items
	form.Append("Text", textArea)
	return form
}

func choicesView() *widget.Box {
	check := widget.NewCheck("Optional", func(value bool) {
		log.Println("Check set to", value)
	})
	radio := widget.NewRadio([]string{"Option 1", "Option 2"}, func(value string) {
		log.Println("Radio set to", value)
	})
	combo := widget.NewSelect([]string{"Option 1", "Option 2"}, func(value string) {
		log.Println("Select set to", value)
	})
	return widget.NewVBox(check, radio, combo)
}

func entryView() *widget.Box {
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter text...")

	content := widget.NewVBox(input, widget.NewButton("Save", func() {
		log.Println("Content was:", input.Text)
	}))
	return content
}

func boxView() *widget.Box {
	content := widget.NewVBox(widget.NewLabel("The top row of VBox"),
		widget.NewHBox(
			widget.NewLabel("Label 1"),
			widget.NewLabel("Label 2")))

	content.Append(widget.NewButton("Add more items", func() {
		content.Append(widget.NewLabel("Appended"))
	}))
	button := widget.NewButton("Split!", func() {
		log.Println("Splitt!")
	})
	content.Append(button)
	return content
}
func fileP(x fyne.URIReadCloser, err error) {

}

func splitView(myWindow fyne.Window) *widget.Box {

	filePath := "???"
	callback := func(x fyne.URIReadCloser, err error) {
		log.Println("sss")
		log.Println(x.URI())
		log.Println(err)
		filePath = x.URI().Extension()
		log.Println(filePath)
	}
	callback2 := func(x fyne.URIWriteCloser, err error) {
		log.Println("sss")
		log.Println(x.URI())
		log.Println(err)

	}
	dial := dialog.NewFileOpen(callback, myWindow)

	fileButton := widget.NewButton("Select File!", func() {
		//NewFileOpen(callback, myWindow)
		dial.Show()
		log.Println("Split!")
	})
	splitButton := widget.NewButton("Split Button", func() {
		log.Println("Split!")
		dialog.NewFileSave(callback2, myWindow).Show()
	})
	content := widget.NewVBox(widget.NewLabel("Split Menu"),
		fileButton,
		splitButton)
	widget.NewLabel(filePath)
	return content

}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("The Splitter")

	tabs := widget.NewTabContainer(
		widget.NewTabItem("Split", splitView(myWindow)),
		widget.NewTabItem("Toolbar", toolbarView()),
		widget.NewTabItem("progress", progressView()),
		widget.NewTabItem("Form", formView(myWindow)),
		widget.NewTabItem("Choices", choicesView()),
		widget.NewTabItem("Box", boxView()),
		widget.NewTabItem("Entry", entryView()))

	//tabs.Append(widget.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))

	tabs.SetTabLocation(widget.TabLocationTop)
	myWindow.Resize(fyne.NewSize(600, 300))
	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
