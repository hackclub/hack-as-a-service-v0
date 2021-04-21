package db

func removeDuplicates(a, b []uint) ([]uint, []uint) {
	ma := make(map[uint]bool)
	mb := make(map[uint]bool)

	for _, item := range a {
		ma[item] = true
	}

	for _, item := range b {
		mb[item] = true
	}

	var newA []uint = nil
	for _, item := range a {
		if _, ok := mb[item]; !ok {
			newA = append(newA, item)
		}
	}

	var newB []uint = nil
	for _, item := range b {
		if _, ok := ma[item]; !ok {
			newB = append(newB, item)
		}
	}

	return newA, newB
}
