// +build !kar

package kar

import "log"

func init() {
	log.Fatal("github.com/omeid/kar improted without correct tag. Do you have proper build tags for your kar file?")
}
