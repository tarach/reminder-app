package remindercreate

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var daysOrder = [7]string{"M", "T", "W", "Th", "F", "Sa", "Su"}

// View builds the reminder creation form UI.
type View struct {
	content fyne.CanvasObject

	nameEntry    *widget.Entry
	typeSelect   *widget.Select
	enabledCheck *widget.Check
	dayChecks    [7]*widget.Check
	alertEntry   *widget.Entry
	counterInit  *widget.Entry
	counterFmt   *widget.Select

	onSave   func(Input)
	onCancel func()
}

// NewView creates a new create view. Use SetCallbacks to bind Save/Cancel.
func NewView() *View {
	v := &View{}
	v.buildForm()
	return v
}

// SetCallbacks sets the Save and Cancel callbacks.
func (v *View) SetCallbacks(onSave func(Input), onCancel func()) {
	v.onSave = onSave
	v.onCancel = onCancel
}

// Content returns the root canvas object for the create screen.
func (v *View) Content() fyne.CanvasObject {
	if v.content == nil {
		return container.NewVBox()
	}
	return v.content
}

func (v *View) buildForm() {
	v.nameEntry = widget.NewEntry()
	v.nameEntry.SetPlaceHolder("Reminder name")

	v.typeSelect = widget.NewSelect([]string{"alert", "counter"}, func(string) {})
	v.typeSelect.SetSelected("alert")

	v.enabledCheck = widget.NewCheck("Enabled", func(bool) {})
	v.enabledCheck.SetChecked(true)

	for i := range v.dayChecks {
		day := daysOrder[i]
		v.dayChecks[i] = widget.NewCheck(day, func(bool) {})
	}

	v.alertEntry = widget.NewEntry()
	v.alertEntry.SetPlaceHolder("Time (e.g. 07:00)")

	v.counterInit = widget.NewEntry()
	v.counterInit.SetPlaceHolder("Initial value (e.g. 00:00)")

	v.counterFmt = widget.NewSelect([]string{"mm:ss", "hh:mm", "number"}, func(string) {})
	v.counterFmt.SetSelected("mm:ss")

	saveBtn := widget.NewButton("Save", func() {
		if v.onSave != nil {
			v.onSave(v.gatherInput())
		}
	})
	saveBtn.Importance = widget.HighImportance

	cancelBtn := widget.NewButton("Cancel", func() {
		if v.onCancel != nil {
			v.onCancel()
		}
	})

	form := container.NewVBox(
		widget.NewLabel("Name"),
		v.nameEntry,
		widget.NewLabel("Type"),
		v.typeSelect,
		v.enabledCheck,
		widget.NewLabel("Active days"),
		container.NewHBox(
			v.dayChecks[0], v.dayChecks[1], v.dayChecks[2], v.dayChecks[3],
			v.dayChecks[4], v.dayChecks[5], v.dayChecks[6],
		),
		widget.NewLabel("Alert time (for type alert)"),
		v.alertEntry,
		widget.NewLabel("Counter initial value (for type counter)"),
		v.counterInit,
		widget.NewLabel("Counter display format"),
		v.counterFmt,
		container.NewHBox(saveBtn, cancelBtn),
	)

	v.content = container.NewScroll(container.NewPadded(form))
}

func (v *View) gatherInput() Input {
	var active [7]bool
	for i := range v.dayChecks {
		active[i] = v.dayChecks[i].Checked
	}
	return Input{
		Name:                 v.nameEntry.Text,
		Type:                 v.typeSelect.Selected,
		Enabled:              v.enabledCheck.Checked,
		ActiveDays:           active,
		AlertValue:            v.alertEntry.Text,
		CounterInitialValue:   v.counterInit.Text,
		CounterDisplayFormat:  v.counterFmt.Selected,
	}
}
