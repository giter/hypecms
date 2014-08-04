// Package skeleton implements a minimalistic but idiomatic plugin for hypeCMS.
package skeleton

import (
	"api/context"
	"routep"
	"gopkg.in/mgo.v2/bson"
)

// Create a type only to spare ourselves from typing map[string]interface{} every time.
type m map[string]interface{}

type H struct {
	uni *context.Uni
}
func Hooks(uni *context.Uni) *H {
	return &H{uni}
}

func (h *H) Front() (bool, error) {
	var hijacked bool
	if _, err := routep.Comp("/skeleton", h.uni.P); err == nil {
		hijacked = true                                    	// This stops the main front loop from executing any further modules.
	}
	return hijacked, nil
}

func (h *H) Install(id bson.ObjectId) error {
	skeleton_options := m{
		"example": "any value",
	}
	q := m{"_id": id}
	upd := m{
		"$addToSet": m{
			"Hooks.Front": "skeleton",
		},
		"$set": m{
			"Modules.skeleton": skeleton_options,
		},
	}
	return h.uni.Db.C("options").Update(q, upd)
}

func (h *H) Uninstall(id bson.ObjectId) error {
	q := m{"_id": id}
	upd := m{
		"$pull": m{
			"Hooks.Front": "skeleton",
		},
		"$unset": m{
			"Modules.skeleton": 1,
		},
	}
	return h.uni.Db.C("options").Update(q, upd)
}
