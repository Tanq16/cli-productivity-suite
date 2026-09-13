package runner

import "github.com/tanq16/cli-productivity-suite/utils"

const packagesUnit = utils.Unit("packages")

type jobResult struct {
	name string
	err  error
}
