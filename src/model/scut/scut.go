// Package scut contains a somewhat ugly but useful collection of frequently appearing patterns to allow faster prototyping.
// Methods here are mainly related to view- or conroller-like parts.
package scut

import (
	"fmt"
	"numcon"
	"io/ioutil"
	"gopkg.in/mgo.v2/bson"
	"path/filepath"
	"sort"
	"strings"
)

// Converts all bson.ObjectId s to string. Usually called before displaying a database query result.
// Input is the result from the database.
func IdsToStrings(v interface{}) {
	switch value := v.(type) {
	case bson.M:
		for i, mem := range value {
			if id, is_id := mem.(bson.ObjectId); is_id {
				value[i] = id.Hex()
			} else {
				IdsToStrings(mem)
			}
		}
	case map[string]interface{}:
		for i, mem := range value {
			if id, is_id := mem.(bson.ObjectId); is_id {
				value[i] = id.Hex()
			} else {
				IdsToStrings(mem)
			}
		}
	case []interface{}:
		for i, mem := range value {
			if id, is_id := mem.(bson.ObjectId); is_id {
				value[i] = id.Hex()
			} else {
				IdsToStrings(mem)
			}
		}
	}
}

// A more generic version of abcKeys. Takes a map[string]interface{} and puts every element of that into an []interface{}, ordered by keys alphabetically.
// TODO: find the intersecting parts between the two functions and refactor.
func OrderKeys(d map[string]interface{}) []interface{} {
	keys := []string{}
	for i, _ := range d {
		keys = append(keys, i)
	}
	sort.Strings(keys)
	ret := []interface{}{}
	for _, v := range keys {
		if ma, is_ma := d[v].(map[string]interface{}); is_ma {
			// RETHINK: What if a key field gets overwritten? Should we name it _key?
			ma["key"] = v
		}
		ret = append(ret, d[v])
	}
	return ret
}

// TODO: ret should contain the rules, so we can display/js validate based on them too.
// Extract module should be modified to not blow up when encountering an unkown rule field, so we can embed metainformation (like text or input, WYSIWYG editor, etc) in the rule too.
//
// Takes a dat map[string]interface{}, and puts every element of that which is defined in r to a slice, sorted by the keys ABC order.
// prior parameter can override the default abc ordering, so keys in prior will be the first ones in the slice, if those keys exist.
func abcKeys(rule map[string]interface{}, dat map[string]interface{}, prior []string) []map[string]interface{} {
	ret := []map[string]interface{}{}
	already_in := map[string]struct{}{}
	for _, v := range prior {
		if _, contains := rule[v]; contains {
			item := map[string]interface{}{v: 1, "key": v}
			if dat != nil {
				item["value"] = dat[v]
			}
			ret = append(ret, item)
			already_in[v] = struct{}{}
		}
	}
	keys := []string{}
	for i, v := range rule {
		// If the value is not false
		if boo, is_boo := v.(bool); !is_boo || boo == true {
			keys = append(keys, i)
		}
	}
	sort.Strings(keys)
	for _, v := range keys {
		if _, in := already_in[v]; !in {
			item := map[string]interface{}{v: 1, "key": v}
			if dat != nil {
				item["value"] = dat[v]
			}
			ret = append(ret, item)
		}
	}
	return ret
}

// Takes an extraction/validation rule, a document and from that creates a slice which can be easily displayed by a templating engine as a html form.
// Takes interface{}s and not map[string]interface{}s to include type checking here, and avoid that boilerplate in caller. 
func RulesToFields(rule interface{}, dat interface{}) ([]map[string]interface{}, error) {
	rm, rm_ok := rule.(map[string]interface{})
	if !rm_ok {
		return nil, fmt.Errorf("Rule is not a map[string]interface{}.")
	}
	datm, datm_ok := dat.(map[string]interface{})
	if !datm_ok && dat != nil {
		return nil, fmt.Errorf("Dat is not a map[string]interface{}.")
	}
	return abcKeys(rm, datm, []string{"title", "name", "slug"}), nil
}

// Gives you back the type of the currently used template (either "private" or public).
func TemplateType(opt map[string]interface{}) string {
	_, priv := opt["TplIsPrivate"]
	var ttype string
	if priv {
		ttype = "private"
	} else {
		ttype = "public"
	}
	return ttype
}

// Gives you back the name of the current template in use.
func TemplateName(opt map[string]interface{}) string {
	tpl, has_tpl := opt["Template"]
	if !has_tpl {
		tpl = "default"
	}
	return tpl.(string)
}

// Decides if a given relative filepath (filep) is a possible module filepath.
// This may be deprecated in the future since it seems so restrictive.
func PossibleModPath(filep string) bool {
	sl := strings.Split(filep, "/")
	return len(sl) >= 2
}

// TODO: Implement file caching here.
// Reads the fi relative filepath from either the current template, or the fallback module tpl folder if fi has at least one slash in it.
// file_reader is optional, falls back to simple ioutil.ReadFile if not given. file_reader will be a custom file_reader with caching soon.
func GetFile(root, fi string, opt map[string]interface{}, host string, file_reader func(string) ([]byte, error)) ([]byte, error) {
	if file_reader == nil {
		file_reader = ioutil.ReadFile
	}
	p := GetTPath(opt, host)
	b, err := file_reader(filepath.Join(root, p, fi))
	if err == nil {
		return b, nil
	}
	if !PossibleModPath(fi) {
		return nil, fmt.Errorf("Not found.")
	}
	mp := GetModTPath(fi)
	return file_reader(filepath.Join(root, mp[0], mp[1]))
}

// Observes opt and gives you back the path of your template eg
// "templates/public/template_name" or "templates/private/hostname/template_name"
func GetTPath(opt map[string]interface{}, host string) string {
	templ := TemplateName(opt)
	ttype := TemplateType(opt)
	if ttype == "public" {
		return filepath.Join("templates", ttype, templ)
	}
	return filepath.Join("templates", ttype, host, templ)
}

// Inp:	"admin/this/that.txt"
// []string{ "modules/admin/tpl", "this/that.txt"}
func GetModTPath(filename string) []string {
	sl := []string{}
	p := strings.Split(filename, "/")
	sl = append(sl, filepath.Join("modules", p[0], "tpl"))
	sl = append(sl, strings.Join(p[1:], "/"))
	return sl
}

func NotAdmin(user interface{}) bool {
	return Ulev(user) < 300
}

func IsAdmin(user interface{}) bool {
	return Ulev(user) >= 300
}

func IsModerator(user interface{}) bool {
	return Ulev(user) >= 200
}

func IsRegistered(user interface{}) bool {
	return Ulev(user) >= 100
}

func IsGuest(user interface{}) bool {
	ulev := Ulev(user)
	return (ulev > 0 && ulev < 100)
}

func IsStranger(user interface{}) bool {
	return Ulev(user) == 0
}

func SolvedPuzzles(user interface{}) bool {
	return Ulev(user) > 1
}

// Gives back the user level.
func Ulev(useri interface{}) int {
	if useri == nil {
		return 0 // useri should never be nil BTW
	}
	user := useri.(map[string]interface{})
	ulev, has := user["level"]
	if !has {
		return 0
	}
	return numcon.IntP(ulev)
}

// Merges b into a (overwriting members in a.
func Merge(a map[string]interface{}, b map[string]interface{}) {
	for i, v := range b {
		a[i] = v
	}
}

// CanonicalHost(uni.Req.Host, uni.Opt)
// Gives you back the canonical address of the site so it can be made available from different domains.
func Host(host string, opt map[string]interface{}) string {
	alias_whitelist, has_alias_whitelist := opt["host_alias_whitelist"]
	if has_alias_whitelist {
		awm := alias_whitelist.(map[string]interface{})
		if _, allowed := awm[host]; !allowed && len(awm) > 0 { // To prevent entirely locking yourself out of the site. Still can introduce problems if misused.
			panic(fmt.Sprintf("Unapproved host alias %v.", host))
		}
	}
	canon_host, has_canon := opt["canonical_host"]
	if !has_canon {
		return host
	}
	return canon_host.(string)
}

func OnlyAdmin(dat map[string]interface{}) {
	if Ulev(dat["_user"]) < 300 {
		panic("Only an admin can do this operation.")
	}
}
