package modules

import (
	"bufio"
	"crypto/ed25519"
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/Linkinlog/gasible/internal/app"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v3"
)

type github struct {
	name        string
	Enabled     bool
	Settings    githubSettings
	application *app.App
}

type githubSettings struct {
	token       string
	TokenEnvKey string `yaml:"token-env-key"`
	SshKeyPath  string `yaml:"ssh-key-path"`
	SshKeyName  string `yaml:"ssh-key-name"`
}

func init() {
	ToBeRegistered = append(ToBeRegistered, &github{
		name:     "GitHub",
		Enabled:  true,
		Settings: githubSettings{},
	})
}

func (gh *github) SetApp(app *app.App) {
	gh.application = app
	// TODO handle this in init like ToBeRegistered
	installer, err := app.ModuleRegistry.GetModule("GenericPackageManager")
	if err != nil {
		panic(err)
	}
	installer.(*GenericPackageManager).AddToInstaller(PackageManagerMap{
		&Brew: {"gh"},
	})
}

func (gh *github) GetName() string { return gh.name }

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

func (gh *github) Config() app.ModuleConfig {
	return app.ModuleConfig{
		Enabled:  gh.Enabled,
		Settings: gh.Settings,
	}
}

func (gh *github) TearDown() error {
	sshKeyDelErr := gh.removeSSHKeys()
	if sshKeyDelErr != nil {
		return sshKeyDelErr
	}
	err := gh.authLogout()
	if err != nil {
		return err
	}
	uninstallGH(gh.application)
	return nil
}

func (gh *github) Setup() error {
	installGH(gh.application)
	if gh.Settings.SshKeyName == "" {
		gh.Settings.SshKeyName = "Gasible Created - " + time.DateTime
	}
	if gh.Settings.TokenEnvKey == "" {
		gh.getTokenFromUser()
	}
	// use token to login
	err := gh.authLogin()
	if err != nil {
		return err
	}
	// generate / prompt for the ssh key,
	// then add the ssh key to gh
	sshErr := gh.addSSHKey(gh.Settings.SshKeyName)
	if sshErr != nil {
		return sshErr
	}
	return nil
}

func (gh *github) Update() error {
	err := upgradeGH(gh.application)
	if err != nil {
		return err
	} else {
		return nil
	}
}

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

func (gh *github) authLogin() error {
	// use token to run `gh auth login --with-token`
	resp, err := gh.application.System.ExecCombinedWithInput(app.NormalRunner{}, "gh", []string{"auth", "login", "--with-token"}, gh.Settings.token)
	if err != nil {
		return fmt.Errorf("authLogin error: %w \n more details: %s", err, string(resp))
	} else {
		return nil
	}
}

func (gh *github) authLogout() error {
	resp, err := gh.application.System.ExecCombinedOutput(app.NormalRunner{}, "gh", []string{"auth", "logout", "--hostname", "github.com"})
	if err != nil {
		log.Fatal(resp, err)
		return err
	}
	return nil
	// return gh.application.System.ExecRun(app.NormalRunner{}, "gh", []string{"auth", "logout", "--hostname", "github.com"})
}

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
	out, err := gh.application.System.ExecCombinedOutput(
		app.NormalRunner{},
		"gh",
		[]string{"ssh-key", "add", gh.Settings.SshKeyPath, "--title", title},
	)
	if err != nil {
		return errors.Join(err, errors.New(string(out)))
	} else {
		return nil
	}
}

func (gh *github) getSSHKeyIDs(sshKeyName string) ([]string, error) {
	out, err := gh.application.System.ExecCombinedOutput(
		app.NormalRunner{},
		"gh",
		[]string{"ssh-key", "list"},
	)
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

func (gh *github) removeSSHKeys() error {
	ids, err := gh.getSSHKeyIDs("Gasible-Created-Key")
	if err != nil {
		return err
	}
	for _, id := range ids {
		out, runErr := gh.application.System.ExecCombinedOutput(app.NormalRunner{}, "gh", []string{"ssh-key", "delete", id, "-y"})
		if runErr != nil {
			return errors.Join(runErr, errors.New(string(out)))
		}
	}
	return nil
}

func installGH(application *app.App) {
	// TODO use GenericPakcageManager
	if _, err := exec.LookPath("apt-get"); err != nil {
		if application.System.Name != "darwin" {
			panic("only apt-get works for now") // TODO
		}
	}
	runner := app.SudoRunner{}

	// Step 1: Check if curl is installed, if not install it
	if _, err := exec.LookPath("curl"); err != nil {
		// curl is not installed, install it
		if err := application.System.ExecRun(&runner, "apt-get", []string{"update"}); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Failed to update apt package list: %v\n", err)
			if err != nil {
				return
			}
			return
		}

		if err := application.System.ExecRun(&runner, "apt-get", []string{"install", "curl", "-y"}); err != nil {
			_, err := fmt.Fprintf(os.Stderr, "Failed to install curl: %v\n", err)
			if err != nil {
				return
			}
			return
		}
	}

	// Step 2: Fetch the GPG key for the GitHub CLI's package repository and install it
	gpgURL := "https://cli.github.com/packages/githubcli-archive-keyring.gpg"
	command := fmt.Sprintf(`curl -fsSL %s | sudo dd of=/usr/share/keyrings/githubcli-archive-keyring.gpg`, gpgURL)
	if err := application.System.ExecRun(&runner, "sh", []string{"-c", command}); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to install GPG key: %v\n", err)
		if err != nil {
			return
		}
		return
	}

	// Step 3: Adds the GitHub CLI's package repository to apt's list of package sources
	command = `deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main`
	command = fmt.Sprintf(`echo "%s" | sudo tee /etc/apt/sources.list.d/github-cli.list > /dev/null`, command)
	if err := application.System.ExecRun(&runner, "sh", []string{"-c", command}); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to add GitHub CLI's package repository: %v\n", err)
		if err != nil {
			return
		}
		return
	}

	// Step 4: Update the apt package lists again
	if err := application.System.ExecRun(&runner, "apt-get", []string{"update"}); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to update apt package list: %v\n", err)
		if err != nil {
			return
		}
		return
	}

	// Step 5: Install the GitHub CLI
	if err := application.System.ExecRun(&runner, "apt-get", []string{"install", "gh", "-y"}); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to install GitHub CLI: %v\n", err)
		if err != nil {
			return
		}
	}

	fmt.Println("Successfully installed GitHub CLI.")
}

func uninstallGH(application *app.App) {
	runner := app.SudoRunner{}
	// Step 1: Uninstall gh
	if err := application.System.ExecRun(&runner, "apt-get", []string{"remove", "gh", "-y"}); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to uninstall GitHub CLI: %v\n", err)
		if err != nil {
			return
		}
		return
	}

	// Step 2: Remove the repository from the list of sources
	if err := application.System.ExecRun(&runner, "rm", []string{"/etc/apt/sources.list.d/github-cli.list"}); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to remove the repository from sources list: %v\n", err)
		if err != nil {
			return
		}
		return
	}

	// Step 3: Remove the keyring
	if err := application.System.ExecRun(&runner, "rm", []string{"/usr/share/keyrings/githubcli-archive-keyring.gpg"}); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to remove the keyring: %v\n", err)
		if err != nil {
			return
		}
		return
	}

	// Step 4: Update the apt package lists after the changes
	if err := application.System.ExecRun(&runner, "apt-get", []string{"update"}); err != nil {
		_, err := fmt.Fprintf(os.Stderr, "Failed to update apt package list: %v\n", err)
		if err != nil {
			return
		}
		return
	}

	fmt.Println("Successfully uninstalled GitHub CLI and cleaned up.")

}

func upgradeGH(application *app.App) error {
	// Check if gh is installed.
	_, err := application.System.ExecCombinedOutput(application.Executor, "type", []string{"-p", "gh"})
	if err != nil {
		return fmt.Errorf("gh is not installed, cannot upgrade: %v", err)
	}

	// Update package lists for upgrades and installations
	err = application.System.ExecRun(application.Executor, "sudo", []string{"apt", "update"})
	if err != nil {
		return fmt.Errorf("error running sudo apt update: %v", err)
	}

	// Upgrade gh
	err = application.System.ExecRun(application.Executor, "sudo", []string{"apt", "upgrade", "gh", "-y"})
	if err != nil {
		return fmt.Errorf("error running sudo apt upgrade gh -y: %v", err)
	}

	return nil
}

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
	err = os.WriteFile(keyPath+".pub", publicKeyBytes, 0644)
	if err != nil {
		return "", err
	}

	return keyPath, nil
}

func userHomeDir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}
