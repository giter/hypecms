package mod

import "modules/skeleton"

func init() {
	modules["skeleton"] = dyn{Hooks: skeleton.Hooks}
}