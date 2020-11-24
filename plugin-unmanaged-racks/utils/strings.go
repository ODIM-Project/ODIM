package utils

type Finder []string

func (f Finder) Find(v string) bool {
	for _, e := range f {
		if e == v {
			return true
		}
	}

	return false
}
