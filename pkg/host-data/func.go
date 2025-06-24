package hostdata

func bytesToString(p []byte) string {
	return string(clrSideSpaces(p))
}

func clrSideSpaces(p []byte) []byte {
	lPtr, rPtr := 0, len(p)-1

	for lPtr <= rPtr && p[lPtr] == ' ' {
		lPtr++
	}

	for rPtr >= lPtr && p[rPtr] == ' ' {
		rPtr--
	}

	if lPtr > rPtr {
		return []byte{}
	}

	return p[lPtr : rPtr+1]
}
