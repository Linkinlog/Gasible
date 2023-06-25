package modules

import (
	"testing"
)

func Test_genericPackageManager_addToInstaller(t *testing.T) {
	type fields struct {
		Enabled  bool
		Settings packageManagerConfig
	}
	type args struct {
		packageMap map[string][]string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			fields: struct {
				Enabled  bool
				Settings packageManagerConfig
			}{
				Enabled: true,
				Settings: packageManagerConfig{
					Manager:      "apt-get",
					Packages:     []string{"foo", "bar"},
					Dependencies: []string{"foo", "bar"},
				},
			},
			args: args{
				packageMap: map[string][]string{
					"apt-get": {"foo", "bar"},
				},
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: struct {
				Enabled  bool
				Settings packageManagerConfig
			}{
				Enabled: true,
				Settings: packageManagerConfig{
					Manager:      "failure",
					Packages:     []string{"foo", "bar"},
					Dependencies: []string{"foo", "bar"},
				},
			},
			args: args{
				packageMap: map[string][]string{
					"failure": {"foo", "bar"},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			packageMan := &genericPackageManager{
				name:     tt.fields.Settings.Manager,
				Enabled:  tt.fields.Enabled,
				Settings: tt.fields.Settings,
			}
			if err := packageMan.addToInstaller(tt.args.packageMap); (err != nil) != tt.wantErr {
				t.Errorf("addToInstaller() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
