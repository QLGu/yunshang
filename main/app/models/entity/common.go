package entity

type regionData map[string]string

func (e regionData) GetById(id string) string {
	v, ok := e[id]
	if !ok {
		return ""
	}
	return v
}
