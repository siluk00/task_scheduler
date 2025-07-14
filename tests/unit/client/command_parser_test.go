package client_test

import (
	"testing"

	"github.com/siluk00/task_scheduler/cmd/client/commands"
	"github.com/stretchr/testify/assert"
)

func TestParseCreateCommand(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{"Valid command", []string{"create", "--id", "1", "--name", "Backup", "--command", "echo hello"}, false},
		{"Missing name", []string{"create", "--id", "2", "--command", "echo"}, true},
		{"Invalid JSON file", []string{"create", "--id", "3", "--name", "jsonFile", "--file", "nonexistent.json"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := commands.NewCreateCommand()
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			assert.Equal(t, tt.wantErr, err != nil)
			t.Log(err)
		})
	}
}
