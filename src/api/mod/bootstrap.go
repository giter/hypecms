package mod

// bs, hahahahaha.
import bs "modules/bootstrap"

func init() {
	modules["bootstrap"] = dyn{Hooks: bs.Hooks, Actions: bs.Actions, Views:bs.Views}
}