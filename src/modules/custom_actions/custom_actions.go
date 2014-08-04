package custom_actions

import (
	"fmt"
	"api/context"
	ca_model "modules/custom_actions/model"
	"jsonp"
	"gopkg.in/mgo.v2/bson"
)

type m map[string]interface{}

func (a *A) Execute() error {
	uni := a.uni
	action_name := uni.Req.Form["action"][0]
	action, has := jsonp.GetM(uni.Opt, "Modules.custom_actions.actions."+action_name)
	if !has {
		return fmt.Errorf("Can't find action %v in custom actions module.", action_name)
	}
	db := uni.Db
	user := uni.Dat["_user"].(map[string]interface{})
	opt := uni.Opt
	inp := map[string][]string(uni.Req.Form)
	typ := action["type"].(string)
	var r error
	switch typ {
	case "vote":
		r = ca_model.Vote(db, user, action, inp)
	case "respond_content":
		r = ca_model.RespondContent(db, user, action, inp, opt)
	default:
		r = fmt.Errorf("Unkown action %v at RunAction.", action_name)
	}
	return r
}


func (h *H) Install(id bson.ObjectId) error {
	custom_action_options := m{}
	q := m{"_id": id}
	upd := m{
		"$set": m{
			"Modules.custom_actions": custom_action_options,
		},
	}
	return h.uni.Db.C("options").Update(q, upd)
}

func (h *H) Uninstall(id bson.ObjectId) error {
	q := m{"_id": id}
	upd := m{
		"$unset": m{
			"Modules.custom_actions": 1,
		},
	}
	return h.uni.Db.C("options").Update(q, upd)
}

type A struct {
	uni *context.Uni
}

func Actions(uni *context.Uni) *A {
	return &A{uni}
}

type H struct {
	uni *context.Uni
}

func Hooks(uni *context.Uni) *H {
	return &H{uni}
}