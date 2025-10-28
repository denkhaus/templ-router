package router

// ErrorTemplate represents an error.templ file for error page presentation
type ErrorTemplate struct {
	// FilePath is the full path to the error.templ file
	FilePath string

	// DirectoryPath is the directory containing this error template
	DirectoryPath string

	// ErrorTypes is a list of error types handled by this template
	ErrorTypes []string

	// ParentErrorTemplate is the path to parent error template if this doesn't override completely
	ParentErrorTemplate string

	// PrecedenceLevel is the level of precedence (closer templates override further ones)
	PrecedenceLevel int

	// ErrorMessages contains mapping of error codes to specific messages
	ErrorMessages map[int]string
}
