package modules_test

import (
	"github.com/Linkinlog/gasible/modules"
	"testing"
)

func Test_genericPackageManager_addToInstaller(t *testing.T) {
	type fields struct {
		Enabled  bool
		Settings modules.PackageManagerConfig
	}
	type args struct {
		packageMap modules.PackageManagerMap
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test all package managers can be added",
			fields: struct {
				Enabled  bool
				Settings modules.PackageManagerConfig
			}{
				Enabled: true,
				Settings: modules.PackageManagerConfig{
					Manager:  "apt-get",
					Packages: []string{"foo", "bar"},
				},
			},
			args: args{
				packageMap: modules.PackageManagerMap{
					&modules.Aptitude: {"foo", "bar"},
					&modules.Brew:     {"foo", "bar"},
					&modules.Dnf:      {"foo", "bar"},
					&modules.Pacman:   {"foo", "bar"},
					&modules.Zypper:   {"foo", "bar"},
				},
			},
			wantErr: false,
		},
		{ // todo these dont work
			name: "test that the correct packages get added",
			fields: struct {
				Enabled  bool
				Settings modules.PackageManagerConfig
			}{
				Enabled: true,
				Settings: modules.PackageManagerConfig{
					Manager:  "failure",
					Packages: []string{"foo", "bar"},
				},
			},
			args: args{
				packageMap: modules.PackageManagerMap{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packageMan := &modules.GenericPackageManager{
				Name:     tt.fields.Settings.Manager,
				Enabled:  tt.fields.Enabled,
				Settings: tt.fields.Settings,
			}
			packageMan.AddToInstaller(tt.args.packageMap)
		})
	}
}
