package entry

import "encoding/json"

type Jsoner interface {
	Lock()
	Unlock()
}

func Json(v Jsoner) string {
	v.Lock()
	defer v.Unlock()
	bs, _ := json.Marshal(v)
	return string(bs)
}
