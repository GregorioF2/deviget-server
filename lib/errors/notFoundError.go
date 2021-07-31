package errors

type NotFoundError struct {
	Err string
}

func (m *NotFoundError) Error() string {
	return m.Err
}
