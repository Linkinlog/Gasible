package yamlParser

import (
	// "os"
	"testing"
)

func TestParseGas(t *testing.T) {
	// We need to make a YAML file but just make it
	//	err, file := makeDefaultYAML()
	//	if err != nil {
	//		panic(err)
	//	}
	//	// Then we need to parse it and make sure we get what we expect
	//	gas, err := ParseGas(file)
	//	if err != nil {
	//		panic(err)
	//	}
	// What do we expect?
	// TODO: Make test cases so we can go over everything in
	// * CreateDefaults() and confirm it equals what we expected
}

//func makeDefaultYAML() (error, string) {
//	file, err := os.CreateTemp("", "temp")
//	if err != nil {
//		return err, ""
//	}
//	defer os.Remove(file.Name())
//
//	_, err = file.WriteString(defaultYAML)
//	if err != nil {
//		return err, ""
//	}
//	return nil, file.Name()
//}

// const defaultYAML = `
// pkg-manager: dnf
// packages:
//     - python3-pip
//     - util-linux-user
//     - wget
//     - neovim
//     - zsh
//     - docker
//     - gh
// installer: true
// teamviewer: true
// ssh: true
// git: true
// hostname: development-station
// staticIP: 192.168.4.20
// mask: 255.255.255.0
// TeamViewerCreds:
//     user: username
//     pass: password
// `
