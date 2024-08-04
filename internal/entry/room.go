package entry

type Room interface {
	Base() *RoomBase
	ID() int64
}

type RoomBase struct {
	id    int64
	teams map[int8]Team
}

func NewRoomBase(id int64, t Team) *RoomBase {
	r := &RoomBase{
		id:    id,
		teams: make(map[int8]Team),
	}
	r.teams[t.ID()] = t
	return r
}

func (r *RoomBase) Base() *RoomBase {
	return r
}

func (r *RoomBase) ID() int64 {
	return r.id
}

func (r *RoomBase) GetTeams() []Team {
	res := make([]Team, 0, len(r.teams))
	for _, t := range r.teams {
		res = append(res, t)
	}
	return res
}

func (r *RoomBase) AddTeam(t Team) {
	r.teams[t.ID()] = t
}

func (r *RoomBase) RemoveTeam(id int8) {
	delete(r.teams, id)
}
