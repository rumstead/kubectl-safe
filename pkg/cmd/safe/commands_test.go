package safe

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestIsVerbSafe(t *testing.T) {
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
			got, err := IsVerbSafe(tt.args.verb)
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
			want: &Commands{safeCmds: map[string]Void{"get": Empty, "list": Empty, "version": Empty}}, wantErr: false},
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

func Test_getSafeCommands(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    *Commands
		wantErr bool
	}{
		{name: "EnvOneValue", value: "one", want: &Commands{safeCmds: map[string]Void{"one": Empty}}, wantErr: false},
		{name: "EnvManyValues", value: "one,two,three", want: &Commands{safeCmds: map[string]Void{"one": Empty, "two": Empty, "three": Empty}}, wantErr: false},
		{name: "EnvNoValues", value: "", want: &DefaultSafeCommands, wantErr: false},
		{name: "EnvEndInComma", value: "one,", want: &Commands{safeCmds: map[string]Void{"one": Empty}}, wantErr: false},
		{name: "File", value: fmt.Sprintf("%svalid-commands.txt", getTestDataDir()), want: &Commands{safeCmds: map[string]Void{"get": Empty, "list": Empty, "version": Empty}}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(KubectlSafeCommands, tt.value)
			got, err := getSafeCommands()
			if (err != nil) != tt.wantErr {
				t.Errorf("getSafeCommands() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getSafeCommands() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_paresSafeCommands(t *testing.T) {
	type args struct {
		commands string
	}
	tests := []struct {
		name    string
		args    args
		want    *Commands
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := paresSafeCommands(tt.args.commands)
			if (err != nil) != tt.wantErr {
				t.Errorf("paresSafeCommands() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paresSafeCommands() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func getTestDataDir() string {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	return basepath + "/testdata/"
}
