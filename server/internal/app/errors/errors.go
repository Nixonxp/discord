package errors

type ErrorWithCode interface {
	Code() string
	Message() string
}

type appError string

const (
	ErrWrongState  appError = "wrong application state"
	ErrMainOmitted appError = "main function is omitted"
	ErrShutdown    appError = "application is in shutdown state"
	ErrTermTimeout appError = "termination timeout"
)

func (e appError) Error() string {
	return string(e)
}

type ArrError []error

func (e ArrError) Error() string {
	if len(e) == 0 {
		return "something went wrong"
	}
	var s = "the following errors occurred:"
	for i := range e {
		s += "\n" + e[i].Error()
	}
	return s
}
