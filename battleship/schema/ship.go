package schema

type Ship struct {
	ship_type		string
	positions 		[] Position
}

type Position struct {
	x		int
	y		int
}