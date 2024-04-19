package util

type CycledStringArray struct {
	Strs   []string
	CurIdx int
}

func CreateCycledStringArray(strs []string) *CycledStringArray {
	csa := &CycledStringArray{}
	csa.Strs = append(csa.Strs, strs...)
	return csa
}

func (csa *CycledStringArray) GetCurrentString() string {
	idx := csa.CurIdx
	str := csa.Strs[idx]
	csa.CurIdx = (idx + 1) % len(csa.Strs)
	return str
}
