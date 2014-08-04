package mod

import te "modules/template_editor"

func init() {
	modules["template_editor"] = dyn{Views: te.Views, Hooks: te.Hooks, Actions: te.Actions}
}