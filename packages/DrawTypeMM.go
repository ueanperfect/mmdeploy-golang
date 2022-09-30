package packages

type singlePlatteMM struct {
	r, g, b int
}

type LineMM struct {
	Point1 PointMM
	Point2 PointMM
}

type PointMM struct {
	X float64
	Y float64
}
