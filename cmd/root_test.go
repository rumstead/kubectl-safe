package cmd

import (
	"testing"
)

func Test_findVerb(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "SimpleVerb", args: []string{"get", "pods"}, want: "get"},
		{name: "FlagBeforeVerb", args: []string{"--context=prod", "delete", "pod", "foo"}, want: "delete"},
		{name: "FlagWithSpaceBeforeVerb", args: []string{"--context", "prod", "delete", "pod", "foo"}, want: "delete"},
		{name: "NamespaceFlagBeforeVerb", args: []string{"-n", "kube-system", "get", "pods"}, want: "get"},
		{name: "MultipleFlagsBeforeVerb", args: []string{"--context", "prod", "-n", "default", "apply", "-f", "foo.yaml"}, want: "apply"},
		{name: "NoArgs", args: []string{}, want: ""},
		{name: "OnlyFlags", args: []string{"--help"}, want: ""},
		{name: "BoolFlagBeforeVerb", args: []string{"--v=5", "get", "pods"}, want: "get"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := findVerb(tt.args)
			if got != tt.want {
				t.Errorf("findVerb() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_extractYesFlag(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantArgs  []string
		wantFound bool
	}{
		{name: "NoYesFlag", args: []string{"delete", "pod", "foo"}, wantArgs: []string{"delete", "pod", "foo"}, wantFound: false},
		{name: "YesLong", args: []string{"--yes", "delete", "pod", "foo"}, wantArgs: []string{"delete", "pod", "foo"}, wantFound: true},
		{name: "YesShort", args: []string{"delete", "pod", "foo", "-y"}, wantArgs: []string{"delete", "pod", "foo"}, wantFound: true},
		{name: "BothYesFlags", args: []string{"--yes", "delete", "-y", "pod"}, wantArgs: []string{"delete", "pod"}, wantFound: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotArgs, gotFound := extractYesFlag(tt.args)
			if gotFound != tt.wantFound {
				t.Errorf("extractYesFlag() found = %v, want %v", gotFound, tt.wantFound)
			}
			if len(gotArgs) != len(tt.wantArgs) {
				t.Errorf("extractYesFlag() args = %v, want %v", gotArgs, tt.wantArgs)
				return
			}
			for i := range gotArgs {
				if gotArgs[i] != tt.wantArgs[i] {
					t.Errorf("extractYesFlag() args[%d] = %v, want %v", i, gotArgs[i], tt.wantArgs[i])
				}
			}
		})
	}
}
