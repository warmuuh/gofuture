package gofuture

type Req interface{}
type Res interface{}

type Future struct {
	resChan chan Res
	errChan chan error
	res     Res
	err     error
	C       chan bool
}

func Defer(f func(Req) Res) func(Req) Future {
	c := make(chan Res, 1)
	return func(r Req) Future {
		go func() {
			c <- f(r)
		}()
		return Future{
			resChan: c,
		}
	}

}

func ResolvedFuture(r Res) Future {
	c := make(chan Res, 1)
	c <- r
	return Future{
		resChan: c,
	}

}

func FailedFuture(err error) Future {
	c := make(chan error, 1)
	c <- err
	return Future{
		errChan: c,
	}

}
func (f *Future) Get() (Res, error) {

	closeAll := func() {
		if f.errChan != nil {
			close(f.errChan)
		}
		if f.resChan != nil {
			close(f.resChan)
		}

	}

	select {
	case res, ok := <-f.resChan:
		if ok {
			f.res = res
			closeAll()
			if futRes, ok := f.res.(Future); ok {
				f.res, f.err = futRes.Get()
			}
		}
	case err, ok := <-f.errChan:
		if ok {
			f.err = err
			closeAll()
		}
	}

	return f.res, f.err
}

func (f Future) Then(success func(Req) Future, failed func(error) Future) Future {
	return Defer(func(Req) Res {
		res, err := f.Get()
		if err != nil {
			if failed != nil {
				return failed(err)
			}
			return FailedFuture(err)
		}

		if success != nil {
			return success(res)
		}
		return ResolvedFuture(res)
	})(nil)
}
func (f *Future) Success(success func(Req) Future) Future {
	return f.Then(success, nil)
}
func (f *Future) Recover(failed func(error) Future) Future {
	return f.Then(nil, failed)
}
