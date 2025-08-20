package config_test

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShell(t *testing.T) {

	var out = bytes.NewBuffer([]byte{})
	bash := exec.Command("bash", "-s")
	in, err := bash.StdinPipe()
	assert.NoError(t, err)
	bash.Stdout = out
	bash.Start()
	// bash.Env = append(bash.Env, "WC_TEST=abc")
	in.Write([]byte("export WC_TEST=abc; echo Hallo $WC_TEST!"))
	in.Close()
	_ = bash.Wait()

	fmt.Println(out.String())

}
