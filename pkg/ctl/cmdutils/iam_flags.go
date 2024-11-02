package cmdutils

import (
	"github.com/spf13/pflag"
)

// AddIAMServiceAccountFilterFlags add common `--include` and `--exclude` flags for filtering iamserviceaccounts
func AddIAMServiceAccountFilterFlags(fs *pflag.FlagSet, includeGlobs, excludeGlobs *[]string) {
	fs.StringSliceVar(includeGlobs, "include", nil,
		"iamserviceaccounts to include (list of globs), e.g.: 'default/s3-reader,*/dynamo-*'")

	fs.StringSliceVar(excludeGlobs, "exclude", nil,
		"iamserviceaccounts to exclude (list of globs), e.g.: 'default/s3-reader,*/dynamo-*'")
}
