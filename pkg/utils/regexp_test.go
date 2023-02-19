package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegexpExtra(t *testing.T) {
	text := "state=hKFo2SA5eEZPZTRjVjJESVhNOUYtZ1pUZEdVVWRIeW1UekNRV6Fur3VuaXZlcnNhbC1sb2dpbqN0aWTZIGVXaDJ1Vm1RRFRDTUJMbDZsMjhwREFTR0J3eWVMRXNZo2NpZNkgVGRKSWNiZTE2V29USHROOTVueXl3aDVFNHlPbzZJdEc\" aria-label=\"\">Sign up</a"
	state, ok := RegexpExtra(text, `state=([a-zA-Z0-9]*)`, 1)
	assert.True(t, ok)
	assert.Equal(t, "hKFo2SA5eEZPZTRjVjJESVhNOUYtZ1pUZEdVVWRIeW1UekNRV6Fur3VuaXZlcnNhbC1sb2dpbqN0aWTZIGVXaDJ1Vm1RRFRDTUJMbDZsMjhwREFTR0J3eWVMRXNZo2NpZNkgVGRKSWNiZTE2V29USHROOTVueXl3aDVFNHlPbzZJdEc", state)
}
