package errors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/go-chi/render"
)

// This has been inspired by Upspin
// https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html
// https://upspin.googlesource.com/upspin/+/master/errors/errors.go

// PathName name of the items path
type PathName string

// Error is the type that implements the error interface.
// It contains a number of fields, each of different type.
// An Error value may leave some values unset.
type Error struct {
	// Path is the path name of the item being accessed.
	Path PathName `json:"path,omitempty"`
	// Op is the operation being performed, usually the name of the method
	// being invoked (Get, Put, etc.). It should not contain an at sign @.
	Op Op `json:"op,omitempty"`
	// Kind is the class of error, such as permission failure,
	// or "Other" if its class is unknown or irrelevant.
	Kind Kind `json:"kind,omitempty"`
	// The underlying error that triggered this one, if any.
	Err error `json:"error,omitempty"`
}

// Op describes an operation, usually as the package and method,
// such as "key/server.Lookup".
type Op string

// Separator is the string used to separate nested errors. By
// default, to make errors easier on the eye, nested errors are
// indented on a new line. A server may instead choose to keep each
// error on a single line by modifying the separator string, perhaps
// to ":: ".
var Separator = ":\n\t"

// Kind defines the kind of error this is, mostly for use by systems
// such as FUSE that must act differently depending on the error.
type Kind uint8

// Kinds of errors.
//
// The values of the error kinds are common between both
// clients and servers. Do not reorder this list or remove
// any items since that will change their values.
// New items must be added only to the end.
const (
	Other      Kind = iota // Unclassified error. This value is not printed in the error message.
	Invalid                // Invalid operation for this type of item.
	Parameters             // Invalid parameters provided.
	Permission             // Permission denied.
	IO                     // External I/O error such as network failure.
	Exist                  // Item already exists.
	NotExist               // Item does not exist.
	Internal               // Internal error or inconsistency.
)

func (k Kind) String() string {
	switch k {
	case Other:
		return "other error"
	case Invalid:
		return "invalid operation"
	case Parameters:
		return "invalid parameters provided"
	case Permission:
		return "permission denied"
	case IO:
		return "I/O error"
	case Exist:
		return "item already exists"
	case NotExist:
		return "item does not exist"
	case Internal:
		return "internal error"
	}
	return "unknown error kind"
}

func (e *Error) isZero() bool {
	return e.Path == "" && e.Op == "" && e.Kind == 0 && e.Err == nil
}

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}

func (e *Error) Error() string {

	b := new(bytes.Buffer)

	if e.Op != "" {
		pad(b, ": ")
		b.WriteString(string(e.Op))
	}

	if e.Path != "" {
		pad(b, ": ")
		b.WriteString(string(e.Path))
	}

	if e.Kind != 0 {
		pad(b, ": ")
		b.WriteString(e.Kind.String())
	}

	if e.Err != nil {
		// Indent on new line if we are cascading non-empty Upspin errors.
		if prevErr, ok := e.Err.(*Error); ok {
			if !prevErr.isZero() {
				pad(b, Separator)
				b.WriteString(e.Err.Error())
			}
		} else {
			pad(b, ": ")
			b.WriteString(e.Err.Error())
		}
	}

	if b.Len() == 0 {
		return "no error"
	}

	return b.String()
}

type errorOutDTO struct {
	Value string `json:"value"`
}

func (we errorOutDTO) Error() string {
	return we.Value
}

// MarshalJSON converts Error to json
func (e *Error) MarshalJSON() ([]byte, error) {
	copy := *e
	_, isError := e.Err.(*Error)
	if e.Err != nil && !isError {
		copy.Err = &errorOutDTO{
			Value: e.Err.Error(),
		}
	}
	return json.Marshal(copy)
}

type errorInDTO struct {
	Path  PathName    `json:"path,omitempty"`
	Op    Op          `json:"op,omitempty"`
	Kind  Kind        `json:"kind,omitempty"`
	Err   *errorInDTO `json:"error,omitempty"`
	Value string      `json:"value,omitempty"`
}

func convert(f *errorInDTO) *Error {

	if f.Value != "" {
		return &Error{
			Kind: Other,
			Err:  Str(f.Value),
		}
	}

	return &Error{
		Op:   f.Op,
		Kind: f.Kind,
		Path: f.Path,
		Err:  convert(f.Err),
	}
}

// UnmarshalJSON implementation
func (e *Error) UnmarshalJSON(b []byte) error {

	var v errorInDTO
	err := json.Unmarshal(b, &v)
	if err != nil {
		return err
	}

	*e = *convert(&v)

	return nil
}

// E builds an error value from its arguments.
// There must be at least one argument or E panics.
// The type of each argument determines its meaning.
// If more than one argument of a given type is presented,
// only the last one is recorded.
//
// The types are:
//	errors.PathName
//		The Upspin path name of the item being accessed.
//	errors.Op
//		The operation being performed, usually the method
//		being invoked (Get, Put, etc.).
//	string
//		Treated as an error message and assigned to the
//		Err field after a call to errors.Str. To avoid a common
//		class of misuse, if the string contains an @, it will be
//		treated as a PathName or UserName, as appropriate. Use
//		errors.Str explicitly to avoid this special-casing.
//	errors.Kind
//		The class of error, such as permission failure.
//	error
//		The underlying error that triggered this one.
//
// If the error is printed, only those items that have been
// set to non-zero values will appear in the result.
//
// If Kind is not specified or Other, we set it to the Kind of
// the underlying error.
//
func E(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case PathName:
			e.Path = arg
		case Op:
			e.Op = arg
		case string:
			e.Err = Str(arg)
		case Kind:
			e.Kind = arg
		case *Error:
			// Make a copy
			copy := *arg
			e.Err = &copy
		case error:
			e.Err = arg
		default:
			_, file, line, _ := runtime.Caller(1)
			log.Printf("errors.E: bad call from %s:%d: %v", file, line, args)
			return Errorf("unknown type %T, value %v in error call", arg, arg)
		}
	}
	// Populate stack information (only in debug mode).
	// e.populateStack()
	prev, ok := e.Err.(*Error)
	if !ok {
		return e
	}
	// The previous error was also one of ours. Suppress duplications
	// so the message won't contain the same kind, file name or user name
	// twice.
	if prev.Path == e.Path {
		prev.Path = ""
	}
	if prev.Kind == e.Kind {
		prev.Kind = Other
	}
	// If this error has Kind unset or Other, pull up the inner one.
	if e.Kind == Other {
		e.Kind = prev.Kind
		prev.Kind = Other
	}
	return e
}

// Recreate the errors.New functionality of the standard Go errors package
// so we can create simple text errors when needed.

// Str returns an error that formats as the given text. It is intended to
// be used as the error-typed argument to the E function.
func Str(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// Errorf is equivalent to fmt.Errorf, but allows clients to import only this
// package for all error handling.
func Errorf(format string, args ...interface{}) error {
	return &errorString{fmt.Sprintf(format, args...)}
}

// HandlerErrorAwareFunc handler implementation which is error aware.
// intended to be used together with errors.WrapHandler
type HandlerErrorAwareFunc func(http.ResponseWriter, *http.Request) error

func getStatusCode(e *Error) int {
	switch e.Kind {
	case Permission:
		return 403
	case IO:
		return 400
	case Invalid:
		return 400
	case Parameters:
		return 400
	case Exist:
		return 409
	case NotExist:
		return 404
	case Other:
		return 500
	case Internal:
		return 500
	}
	panic("Unknown error kind. Always provide an error kind!")
}

// WrapHandler takes care of wrapping a HandlerErrorAwareFunc into an
// http.HandlerFunc. If the inner func errors, it will render the request
func WrapHandler(handler HandlerErrorAwareFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			RenderError(w, r, err)
		}
	}
}

// RenderError applies an error to a request
func RenderError(w http.ResponseWriter, r *http.Request, err error) {
	e, ok := err.(*Error)
	if ok {
		render.Status(r, getStatusCode(e))
	} else {
		render.Status(r, 500)
	}
	render.JSON(w, r, err)
}

// DecodeError decodes an http response to an Error
func DecodeError(resp *http.Response) error {
	const op Op = "errors.decode"

	output := &Error{}
	err := json.NewDecoder(resp.Body).Decode(&output)
	if err != nil {
		return E(op, Internal, err)
	}

	return output
}
