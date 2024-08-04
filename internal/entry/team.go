package entry

type Team interface {
	Base() *TeamBase
	ID() int8
}

type TeamBase struct {
	id     int8
	groups map[int64]Group
}

func NewTeam(id int8, g Group) Team {
	t := &TeamBase{
		id:     id,
		groups: make(map[int64]Group),
	}
	t.groups[g.ID()] = g
	return t
}

func (t *TeamBase) Base() *TeamBase {
	return t
}

func (t *TeamBase) ID() int8 {
	return t.id
}

func (t *TeamBase) GetGroups() []Group {
	res := make([]Group, len(t.groups))
	i := 0
	for _, g := range t.groups {
		res[i] = g
		i++
	}
	return res
}

func (t *TeamBase) AddGroup(g Group) {
	t.groups[g.ID()] = g
}

func (t *TeamBase) RemoveGroup(id int64) {
	delete(t.groups, id)
}
