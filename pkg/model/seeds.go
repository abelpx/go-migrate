package model

type Seeds struct {
	Err error
}

func (s *Seeds) Error() string {
	if s.Err == nil {
		return ""
	}

	return s.Err.Error()
}

func NewSeed(err error) *Seeds {
	return &Seeds{
		Err: err,
	}
}
