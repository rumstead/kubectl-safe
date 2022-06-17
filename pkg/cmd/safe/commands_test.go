package safe

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func Test_isVerbSafe(t *testing.T) {
	type args struct {
		verb string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{name: "DefaultSafeCommandTrue", args: args{verb: "get"}, want: true, wantErr: false},
		{name: "DefaultSafeCommandFalse", args: args{verb: "delete"}, want: false, wantErr: false},
		{name: "EmptyVerb", args: args{verb: ""}, want: true, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isVerbSafe(tt.args.verb)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsVerbSafe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsVerbSafe() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getCommandsFromFile(t *testing.T) {
	type args struct {
		commands string
	}
	tests := []struct {
		name    string
		args    args
		want    *Commands
		wantErr bool
	}{
		{name: "ValidFile", args: args{commands: "./testdata/valid-commands.txt"},
			want: &Commands{cmds: map[string]Void{"get": Empty, "list": Empty, "version": Empty}}, wantErr: false},
		{name: "FileDoesntExist", args: args{commands: "foo"},
			want: nil, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCommandsFromFile(tt.args.commands)
			if (err != nil) != tt.wantErr {
				t.Errorf("getCommandsFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getCommandsFromFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseCommands(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    *Commands
		env     string
		wantErr bool
	}{
		// safe commands
		{name: "EnvOneValue", value: "one", want: &Commands{cmds: map[string]Void{"one": Empty}}, wantErr: false, env: KubectlSafeCommands},
		{name: "EnvManyValues", value: "one,two,three", want: &Commands{cmds: map[string]Void{"one": Empty, "two": Empty, "three": Empty}}, wantErr: false, env: KubectlSafeCommands},
		{name: "EnvNoValues", value: "", want: &EmptyCommands, wantErr: false, env: KubectlSafeCommands},
		{name: "EnvEndInComma", value: "one,", want: &Commands{cmds: map[string]Void{"one": Empty}}, wantErr: false, env: KubectlSafeCommands},
		{name: "File", value: fmt.Sprintf("%svalid-commands.txt", getTestDataDir()), want: &Commands{cmds: map[string]Void{"get": Empty, "list": Empty, "version": Empty}}, wantErr: false, env: KubectlSafeCommands},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.env, tt.value)
			got, err := parseCommands(tt.env)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCommands() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCommands() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isDryRun(t *testing.T) {
	type args struct {
		cmd []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "NotDryRun", args: args{cmd: []string{"delete", "pod", "foo"}}, want: false},
		{name: "DryRun", args: args{cmd: []string{"create", "pod", "foo", "--dry-run=server"}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isDryRun(tt.args.cmd); got != tt.want {
				t.Errorf("isDryRun() = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestDataDir() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return basepath + "/testdata/"
}

func TestIsSafe(t *testing.T) {
	type args struct {
		verb string
		args []string
	}
	tests := []struct {
		name     string
		args     args
		want     bool
		wantErr  bool
		env      string
		envValue string
	}{
		{name: "DryRunNotSafe", args: args{
			verb: "delete",
			args: []string{"delete", "pod", "--dry-run=client"},
		}, want: true, wantErr: false, env: "foo", envValue: ""},
		{name: "NotDryRunSafe", args: args{
			verb: "get",
			args: []string{"get", "pod"},
		}, want: true, wantErr: false, env: "foo", envValue: ""},
		{name: "SetUnsafeCommands", args: args{
			verb: "delete",
			args: []string{"delete", "pod"},
		}, want: false, wantErr: false, env: KubectlUnsafeCommands, envValue: "delete"},
		{name: "SetUnsafeCommandsDiffer", args: args{
			verb: "delete",
			args: []string{"delete", "pod"},
		}, want: true, wantErr: false, env: KubectlUnsafeCommands, envValue: "apply"},
		{name: "OverrideSafeCommands", args: args{
			verb: "delete",
			args: []string{"delete", "pod"},
		}, want: true, wantErr: false, env: KubectlSafeCommands, envValue: "delete"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(tt.env, tt.envValue)
			got, err := IsSafe(tt.args.verb, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsSafe() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsSafe() got = %v, want %v", got, tt.want)
			}
		})
	}
}
