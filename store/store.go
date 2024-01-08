package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

var files = make(map[string]map[string]map[string]interface{})

func writeFile(name string, data []byte) {
	psp := strings.Split(name, "/")
	if len(psp) > 1 {
		os.MkdirAll(strings.Join(psp[:len(psp)-1], "/"), 0755)
	}
	os.WriteFile(name, data, 0755)
}

func read(path string) map[string]map[string]interface{} {
	var output map[string]map[string]interface{}
	d, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return map[string]map[string]interface{}{}
	}
	json.Unmarshal(d, &output)
	return output
}

func addEntryTo(path string, id string, data ...map[string]interface{}) {
	if files[path] == nil {
		files[path] = read(path)
	}
	if len(data) > 0 && data[0] != nil {
		files[path][id] = data[0]
	} else {
		files[path][id] = map[string]interface{}{}
	}
	d, _ := json.Marshal(files[path])
	writeFile(path, d)
}

func removeEntryFrom(path string, data string) {
	if files[path] == nil {
		files[path] = read(path)
	}
	delete(files[path], data)
	d, _ := json.Marshal(files[path])
	writeFile(path, d)
}

func get(path string, data string) (bool, map[string]interface{}) {
	if files[path] == nil {
		files[path] = read(path)
	}
	if d, ok := files[path][data]; ok {
		return true, d
	}
	return false, nil
}

func getAll(path string) map[string]map[string]interface{} {
	if files[path] == nil {
		files[path] = read(path)
	}
	return files[path]
}

func AddCornered(id int64, data ...map[string]interface{}) {
	addEntryTo("data/cornered.json", fmt.Sprint(id), data...)
}

func RemoveCornered(id int64) {
	removeEntryFrom("data/cornered.json", fmt.Sprint(id))
}

func GetCorner(id int64) (bool, map[string]interface{}) {
	return get("data/cornered.json", fmt.Sprint(id))
}

func GetCorners() map[string]map[string]interface{} {
	return getAll("data/cornered.json")
}

func GetMuted() map[string]map[string]interface{} {
	return getAll("data/muted.json")
}

func AddMuted(id int64) {
	addEntryTo("data/muted.json", fmt.Sprint(id))
}

func RemoveMuted(id int64) {
	removeEntryFrom("data/muted.json", fmt.Sprint(id))
}

func IsMuted(id int64) bool {
	d, _ := get("data/muted.json", fmt.Sprint(id))
	return d
}

func Warnings(id int64) int {
	_, d := get("data/warnings.json", fmt.Sprint(id))
	if d == nil {
		return 0
	}
	w, _ := d["warnings"].(float64)
	return int(w)
}

func Warn(id int64) {
	addEntryTo("data/warnings.json", fmt.Sprint(id), map[string]interface{}{
		"warnings": float64(Warnings(id) + 1),
	})
}

func Unwarn(id int64) {
	w := Warnings(id)
	if w == 0 {
		return
	}
	addEntryTo("data/warnings.json", fmt.Sprint(id), map[string]interface{}{
		"warnings": float64(w - 1),
	})
}

func SetWarnings(id int64, w uint) {
	addEntryTo("data/warnings.json", fmt.Sprint(id), map[string]interface{}{
		"warnings": float64(w),
	})
}
