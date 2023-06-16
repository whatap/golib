package percentutil

func TopFloat(src float32, top float32) float32 {
	if src > top {
		return top
	}
	return src
}
