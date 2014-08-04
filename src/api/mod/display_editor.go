package mod

import de "modules/display_editor"

func init() {
	modules["display_editor"] = dyn{Views: de.Views, Hooks: de.Hooks, Actions: de.Actions}
}