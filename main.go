package main

import (
	"delivery/internal/services/document_builder"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	builder := document_builder.New()

	w := a.NewWindow("Генератор документов")
	w.Resize(fyne.Size{
		Width:  600.00,
		Height: 600.00,
	})
	checkButton := addFileOpenButton(w, "Добавить счет на оплату", builder.SetCheckFilePath)
	decipherButton := addFileOpenButton(w, "Добавить расшифровку", builder.SetDecipherFilePath)
	go func() {
		for {
			if checkButton.Text == "Счёт добавлен" && decipherButton.Text == "Расшифровка добавлена" {
				return
			}
			if builder.HasCheckFilePath() && checkButton.Text != "Счёт добавлен" {
				checkButton.Text = "Счёт добавлен"
				checkButton.Refresh()
				checkButton.Disable()
			}
			if builder.HasDecipherFilePath() && decipherButton.Text != "Расшифровка добавлена" {
				decipherButton.Text = "Расшифровка добавлена"
				decipherButton.Refresh()
				decipherButton.Disable()
			}
		}
	}()
	dates := widget.NewEntry()
	dates.SetPlaceHolder("Заполнить даты поездки")
	year := widget.NewEntry()
	year.SetPlaceHolder("Укажите год")
	widgetContainer := container.NewVBox(
		checkButton,
		decipherButton,
		dates,
		year,
		widget.NewButton("Сформировать", func() {
			if dates.Text != "" {
				builder.SetRawDates(dates.Text)
				dates.Disable()
			}
			if year.Text != "" {
				builder.SetYear(year.Text)
				year.Disable()
			}
			if !builder.Validate() {
				dialog.ShowInformation("Ошибка", "Заполните все поля", w)
			}
			err := builder.Init()
			if err != nil {
				fmt.Println(err)
				dialog.ShowInformation("Ошибка", fmt.Sprintf("Ошибка при заполнении полей: %v", err), w)
			}
			err = builder.Build()
			if err != nil {
				fmt.Println(err)
				dialog.ShowInformation("Ошибка", fmt.Sprintf("Ошибка при заполнении полей: %v", err), w)
			}
			dialog.ShowInformation("Документы сформированы", "Документы успешно сформированы", w)
		}),
	)

	w.SetContent(widgetContainer)
	w.ShowAndRun()
}

func addFileOpenButton(parent fyne.Window, label string, pathFunc func(path string)) *widget.Button {
	return widget.NewButton(label, func() {
		dialog.ShowFileOpen(func(list fyne.URIReadCloser, err error) {
			if list == nil {
				return
			}
			if err != nil {
				dialog.ShowError(err, parent)
				return
			}
			filePath := list.URI().Path()
			if filePath != "" {
				pathFunc(filePath)
			}
		}, parent)
	})
}
