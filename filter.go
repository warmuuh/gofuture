package common

type Filter func(Req, Service) Future

func (f Filter) AndThen(next Filter) Filter {
	return func(r Req, s Service) Future {
		return f(r, next.And(s))
	}
}
func (f Filter) And(s Service) Service {
	return func(r Req) Future {
		return f(r, s)
	}
}
