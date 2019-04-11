package room

func (r *Room) changeTurn() {
	if (*r.active).ID() == r.p1.ID() {
		r.active = &r.p2
	} else {
		r.active = &r.p1
	}
}

//Maybe Remove this, GameField already has a validation
func (r *Room) validateCoordinates(x, y int) {

}
