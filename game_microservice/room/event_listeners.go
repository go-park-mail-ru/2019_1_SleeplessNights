package room

func (r *Room) changeTurn() {
	if (*r.active).ID() == r.p1.ID() {
		r.active = &r.p2
	} else {
		r.active = &r.p1
	}
}
