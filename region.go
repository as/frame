package frame

// It must hold that q0 <= q1
func region3(r, q0, q1 int64) int {
	if r <= q0 {
		return -1
	}
	if r > q1 {
		return 1
	}
	return 0
}

func coInsert(r0, r1, q0, q1 int64) (int64, int64) {
	dx := r1 - r0 + 1
	if dx == 0 {
		return q0, q1
	}
	switch region3(r0, q0, q1) {
	case -1:
		q0 += dx
		q1 += dx
	case 0:
		q1 += dx
	case 1:
		// nop
	}
	return q0, q1
}
