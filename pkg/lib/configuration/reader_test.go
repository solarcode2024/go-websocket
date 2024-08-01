package configuration

import "testing"

func TestReader(t *testing.T) {
	configuration, err := Reader("dummy.yaml")
	if err != nil {
		t.Fatalf("error reading config file: %v", err.Error())
	}

	t.Logf("configuration: %v", configuration)
}
