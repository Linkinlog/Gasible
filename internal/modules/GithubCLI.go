package modules

import (
	"bufio"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/Linkinlog/gasible/internal/app"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
)

// init
// This should really just handle registering the module in the registry.
func init() {
	ToBeRegistered = append(ToBeRegistered, &github{
		name:     "GitHub",
		Enabled:  true,
		Settings: githubSettings{},
	})
	ToBeInstalled[&brew] = []string{"gh"}
}

// github implements the module interface, so we can execute system commands.
type github struct {
	name        string
	Enabled     bool
	Settings    githubSettings
	application *app.App
}

// githubSettings is the settings struct for the github module, allows user to set the token env key.
type githubSettings struct {
	token       string `yaml:"-"`
	TokenEnvKey string `yaml:"token-env-key"`
	SshKeyPath  string `yaml:"-"`
}

// SetApp sets the application field as the app that is passed in.
func (gh *github) SetApp(app *app.App) {
	gh.application = app
}

// GetName returns the name field of the github struct.
func (gh *github) GetName() string { return gh.name }

// ParseConfig takes in a map that ideally contains a YAML structure, to be marshalled into the config.
func (gh *github) ParseConfig(rawConfig map[string]interface{}) error {
	configBytes, err := yaml.Marshal(rawConfig)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(configBytes, gh)
	if err != nil {
		return err
	}

	return nil
}

// Config returns the shallow-copied module config from our module's config.
func (gh *github) Config() app.ModuleConfig {
	return app.ModuleConfig{
		Enabled:  gh.Enabled,
		Settings: gh.Settings,
	}
}

// TearDown will uninstall `gh` and remove any created keys from GitHub.
func (gh *github) TearDown() error {
	sshKeyDelErr := gh.removeSSHKeys()
	if sshKeyDelErr != nil {
		return sshKeyDelErr
	}
	err := gh.authLogout()
	if err != nil {
		return err
	}
	gh.uninstallGH()
	return nil
}

// Setup will install `gh` and log in to the CLI application, then add an SSH key to GitHub.
func (gh *github) Setup() error {
	// TODO figure out other package managers
	if _, err := exec.LookPath("apt-get"); err == nil {
		gh.installGH()
	}
	// else check if gh is installed, so we don't explode if it's not

	gh.getTokenFromUser()
	// use token to login
	err := gh.authLogin()
	if err != nil {
		return err
	}
	// generate / prompt for the ssh key,
	// then add the ssh key to gh
	sshErr := gh.addSSHKey("Gasible-Generated-Key")
	if sshErr != nil {
		return sshErr
	}
	return nil
}

// Update will run the update command on the chosen package manager.
func (gh *github) Update() error {
	err := gh.upgradeGH()
	if err != nil {
		return err
	} else {
		return nil
	}
}

// system returns the Syscall that is currently in use by the module registry.
func (gh *github) system() *SysCall {
	sysCallMod := gh.application.ModuleRegistry.GetModule("SysCall")
	return sysCallMod.(*SysCall)
}

// getTokenFromUser will check the TokenEnvKey for a key or prompt the user for one if one isn't found.
func (gh *github) getTokenFromUser() {
	if ghToke, ok := os.LookupEnv(gh.Settings.TokenEnvKey); ok {
		gh.Settings.token = ghToke
		return
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter GitHub token: ")
	token, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("getTokenFromUser error: %q", err)
	} else {
		// Trim the newline from the token
		gh.Settings.token = strings.TrimSuffix(token, "\n")
	}
}

// authLogin runs the auth login --with-token command to authenticate with gh.
func (gh *github) authLogin() error {
	// use token to run `gh auth login --with-token`
	resp, err := gh.system().ExecWithInput("gh", []string{"auth", "login", "--with-token"}, gh.Settings.token, false)
	if err != nil {
		return fmt.Errorf("authLogin error: %w \n more details: %s", err, string(resp))
	} else {
		return nil
	}
}

// authLogout runs the auth logout --hostname github.com to de-authenticate with gh.
func (gh *github) authLogout() error {
	resp, err := gh.system().Exec("gh", []string{"auth", "logout", "--hostname", "github.com"}, false)
	if err != nil {
		log.Fatal(resp, err)
		return err
	}
	return nil
}

// addSSHKey will add our new SSH Key locally and to GitHub.
func (gh *github) addSSHKey(title string) error {
	// create a new ssh key or specify an existing one.
	if gh.Settings.SshKeyPath == "" {
		keyPath, sshErr := generateSSHKeys("github-gasible")
		if sshErr != nil {
			log.Fatal(sshErr)
		}
		gh.Settings.SshKeyPath = keyPath + ".pub"
	}
	// use it with `gh ssh-key add "FILEPATH" --title "TITLE"`.
	out, err := gh.system().Exec("gh", []string{"ssh-key", "add", gh.Settings.SshKeyPath, "--title", title}, false)
	if err != nil {
		var outputErr = errors.New(string(out))
		return errors.Join(err, outputErr)
	} else {
		return nil
	}
}

// getSSHKeyIDs is to retrieve all SSH keys that we believe have been made by this program.
func (gh *github) getSSHKeyIDs(sshKeyName string) ([]string, error) {
	out, err := gh.system().Exec("gh", []string{"ssh-key", "list"}, false)
	if err != nil {
		return []string{string(out)}, err
	}

	var ids []string
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 4 && fields[0] == sshKeyName {
			ids = append(ids, fields[4]) // Assuming ID is in the second column
		}
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("ssh key with name %s not found", sshKeyName)
	}

	return ids, nil
}

// removeSSHKeys will remove all SSH keys from GitHub that we believe we created.
func (gh *github) removeSSHKeys() error {
	ids, err := gh.getSSHKeyIDs("Gasible-Generated-Key")
	if err != nil {
		return err
	}
	for _, id := range ids {
		out, runErr := gh.system().Exec("gh", []string{"ssh-key", "delete", id, "-y"}, false)
		if runErr != nil {
			return errors.Join(runErr, errors.New(string(out)))
		}
	}
	return nil
}

// installGh installs the gh cli application.
func (gh *github) installGH() {
	// Step 1: Check if curl is installed, if not Install it
	if _, err := exec.LookPath("curl"); err != nil {
		// curl is not installed, Install it
		if _, execErr := gh.system().Exec("apt-get", []string{"update"}, true); execErr != nil {
			_, printErr := fmt.Fprintf(os.Stderr, "Failed to Update apt package list: %v\n", execErr)
			if printErr != nil {
				return
			}
			return
		}

		if _, execErr := gh.system().Exec("apt-get", []string{"install", "curl", "-y"}, true); execErr != nil {
			_, printErr := fmt.Fprintf(os.Stderr, "Failed to Install curl: %v\n", execErr)
			if printErr != nil {
				return
			}
			return
		}
	}

	// Step 2: Fetch the GPG key for the GitHub CLI's package repository and Install it
	gpgURL := "https://cli.github.com/packages/githubcli-archive-keyring.gpg"
	command := fmt.Sprintf(`curl -fsSL %s | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg`, gpgURL)
	if _, keyRingInstallErr := gh.system().Exec("sh", []string{"-c", command}, false); keyRingInstallErr != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to Install GPG key: %v\n", keyRingInstallErr)
		if err != nil {
			return
		}
		return
	}

	// Step 3: Adds the GitHub CLI's package repository to aptitude's list of package sources
	command = `deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main`
	command = fmt.Sprintf(`echo "%s" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null`, command)
	if _, sourcesInstallErr := gh.system().Exec("sh", []string{"-c", command}, false); sourcesInstallErr != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to add GitHub CLI's package repository: %v\n", sourcesInstallErr)
		if err != nil {
			return
		}
		return
	}

	// Step 4: Update the apt package lists again
	if _, updateErr := gh.system().Exec("apt-get", []string{"update"}, true); updateErr != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to Update apt package list: %v\n", updateErr)
		if err != nil {
			return
		}
		return
	}

	// Step 5: Install the GitHub CLI
	if _, installErr := gh.system().Exec("apt-get", []string{"install", "gh", "-y"}, true); installErr != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to Install GitHub CLI: %v\n", installErr)
		if err != nil {
			return
		}
	}

	log.Println("Successfully installed GitHub CLI.")
}

// uninstallGh uninstalls the gh cli application.
func (gh *github) uninstallGH() {
	// Step 1: Uninstall gh
	if _, removeErr := gh.system().Exec("apt-get", []string{"remove", "gh", "-y"}, true); removeErr != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to Uninstall GitHub CLI: %v\n", removeErr)
		if err != nil {
			return
		}
		return
	}

	// Step 2: Remove the repository from the list of sources
	if _, sourcesRemoveErr := gh.system().Exec("rm", []string{"/etc/apt/sources.list.d/github-cli.list"}, true); sourcesRemoveErr != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to remove the repository from sources list: %v\n", sourcesRemoveErr)
		if err != nil {
			return
		}
		return
	}

	// Step 3: Remove the keyring
	if _, keyringRemoveErr := gh.system().Exec("rm", []string{"/usr/share/keyrings/githubcli-archive-keyring.gpg"}, true); keyringRemoveErr != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to remove the keyring: %v\n", keyringRemoveErr)
		if err != nil {
			return
		}
		return
	}

	// Step 4: Update the apt package lists after the changes
	if _, aptUpdateErr := gh.system().Exec("apt-get", []string{"update"}, true); aptUpdateErr != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to Update apt package list: %v\n", aptUpdateErr)
		if err != nil {
			return
		}
		return
	}

	log.Println("Successfully uninstalled GitHub CLI and cleaned up.")
}

// upgradeGH upgrades the gh cli application.
func (gh *github) upgradeGH() error {
	// Check if gh is installed.
	_, err := gh.system().Exec("type", []string{"-p", "gh"}, false)
	if err != nil {
		return fmt.Errorf("gh is not installed, cannot upgrade: %v", err)
	}

	// Update package lists for upgrades and installations
	_, err = gh.system().Exec("sudo", []string{"apt", "update"}, false)
	if err != nil {
		return fmt.Errorf("error running sudo apt Update: %v", err)
	}

	// Upgrade gh
	_, err = gh.system().Exec("sudo", []string{"apt", "upgrade", "gh", "-y"}, false)
	if err != nil {
		return fmt.Errorf("error running sudo apt upgrade gh -y: %v", err)
	}

	return nil
}

// generateSSHKeys will create new ssh keys for a given filename.
func generateSSHKeys(fileName string) (string, error) {
	keyPath := path.Join(userHomeDir(), ".ssh", fileName)
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return "", err
	}

	privateKeyBytes := privateKey.Seed()
	err = os.WriteFile(keyPath, privateKeyBytes, 0600)
	if err != nil {
		return "", err
	}

	pub, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	publicKeyBytes := ssh.MarshalAuthorizedKey(pub)
	err = os.WriteFile(keyPath+".pub", publicKeyBytes, 0600)
	if err != nil {
		return "", err
	}

	return keyPath, nil
}

// userHomeDir gets the home directory for the user.
func userHomeDir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}
