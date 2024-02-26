package build

import "fmt"

func (o *SyncStore[T]) Exec(key string, f func() (T, error)) (T, error) {
	val, err := o.sf.Do(key, func() (_ interface{}, err error) {

		defer func() {
			if rErr := recover(); rErr != nil {
				err = retrieveError(key, rErr)
			}
		}()
		v, err, ok := o.results.Load(key)
		if ok {
			return v, err
		}
		v, err = f()
		o.results.Store(key, v, err)
		return v, err
	})

	var defaultT T
	if err != nil {
		return defaultT, err
	}
	switch t := val.(type) {
	case error:
		return defaultT, t
	case T:
		return t, nil
	default:
		return defaultT, err
	}
}

func (m *syncMap[T]) Store(k any, v T, err error) {
	if err != nil {
		m.Map.Store(k, err)
	} else {
		m.Map.Store(k, v)
	}
}

func (m *syncMap[T]) Load(k any) (v T, err error, ok bool) {
	val, found := m.Map.Load(k)
	if !found {
		return
	}
	ok = true
	switch t := val.(type) {
	case error:
		err = t
	case T:
		v = t
	}
	return
}

func (e StoreError) Error() string {
	return e.message
}

func retrieveError(key string, i interface{}) StoreError {
	return StoreError{
		message: fmt.Sprintf("internal error retrieving cached result for key %s: %v", key, i),
	}
}
