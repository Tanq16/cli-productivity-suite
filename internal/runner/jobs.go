package runner

type jobResult struct {
	name    string
	message string
	err     error
}

type checkResult struct {
	name    string
	current string
	latest  string
	status  string
}
