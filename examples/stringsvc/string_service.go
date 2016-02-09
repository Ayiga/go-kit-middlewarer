package stringsvc

//go:generate go-kit-middlewarer -type=StringService

// StringService from go-kit/kit's example
type StringService interface {
	// Uppercase returns an uppercase version of the given string, or an error.
	Uppercase(str string) (upper string, err error)

	// Count returns the length of the given string
	Count(str string) (count int)
}
