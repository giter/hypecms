package mod

import c "modules/content"

func init() {
	modules["content"] = dyn{Views: c.Views, Hooks: c.Hooks, Actions: c.Actions}
}