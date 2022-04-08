package util

func FilterExactValue(target []int32, value int32) []int32 {
	var retv []int32

	for _, v := range target {
		if v != value {
			retv = append(retv, v)
		}
	}
	return retv
}
