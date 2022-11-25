package zipper

func defaultZip[Type comparable](x Zippable[Type], y Zippable[Type]) *Result[Type] {
	if x.Key() != y.Key() {
		return nil
	}
	return &Result[Type]{
		First:  x,
		Second: y,
		Key:    x.Key(),
	}
}
