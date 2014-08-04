package mod

import "modules/user"

func init() {
	modules["user"] = dyn{Hooks: user.Hooks, Actions: user.Actions}
}