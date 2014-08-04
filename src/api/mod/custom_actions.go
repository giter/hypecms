package mod

import ca "modules/custom_actions"

func init() {
	modules["custom_actions"] = dyn{Hooks: ca.Hooks, Actions: ca.Actions}
}