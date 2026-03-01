package ansible

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

type Executor struct {
	hostname      string
	port          int
	username      string
	sshKeyPath    string
	composeContent string
}

type InventoryData struct {
	Hostname string
	Port     int
	Username string
	SSHKey   string
}

func NewExecutor(hostname string, port int, username, sshKeyPath, composeContent string) *Executor {
	return &Executor{
		hostname:      hostname,
		port:          port,
		username:      username,
		sshKeyPath:    sshKeyPath,
		composeContent: composeContent,
	}
}

func (e *Executor) Run(onOutput func(string)) error {
	dir, err := os.MkdirTemp("", "daos-deploy-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(dir)

	inventoryPath := filepath.Join(dir, "inventory")
	if err := e.createInventory(inventoryPath); err != nil {
		return fmt.Errorf("failed to create inventory: %w", err)
	}

	composePath := filepath.Join(dir, "docker-compose.yml")
	if err := os.WriteFile(composePath, []byte(e.composeContent), 0644); err != nil {
		return fmt.Errorf("failed to write compose file: %w", err)
	}

	playbookPath := filepath.Join(dir, "deploy.yml")
	if err := e.createPlaybook(playbookPath, composePath); err != nil {
		return fmt.Errorf("failed to create playbook: %w", err)
	}

	cmd := exec.Command("ansible-playbook", "-i", inventoryPath, playbookPath, "-v")
	cmd.Dir = dir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		output := stderr.String()
		if output == "" {
			output = stdout.String()
		}
		onOutput("Error: " + err.Error())
		onOutput(output)
		return fmt.Errorf("ansible failed: %w", err)
	}

	onOutput(stdout.String())
	return nil
}

func (e *Executor) createInventory(path string) error {
	tmpl := `{{.Hostname}} ansible_port={{.Port}} ansible_user={{.Username}} ansible_private_key_file={{.SSHKey}}`

	t, err := template.New("inventory").Parse(tmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	data := InventoryData{
		Hostname: e.hostname,
		Port:     e.port,
		Username: e.username,
		SSHKey:   e.sshKeyPath,
	}

	if err := t.Execute(&buf, data); err != nil {
		return err
	}

	return os.WriteFile(path, buf.Bytes(), 0644)
}

func (e *Executor) createPlaybook(path, composePath string) error {
	playbook := `- name: Deploy Docker Compose application
  hosts: all
  gather_facts: false
  tasks:
    - name: Ensure Docker is installed
      apt:
        name: docker.io
        state: present
      become: yes

    - name: Ensure Docker Compose is installed
      get_url:
        url: https://github.com/docker/compose/releases/download/v2.24.0/docker-compose-linux-x86_64
        dest: /usr/local/bin/docker-compose
        mode: '0755'
      become: yes
      ignore_errors: yes

    - name: Create deployment directory
      file:
        path: /opt/daos-deploy
        state: directory
      become: yes

    - name: Copy docker-compose.yml
      copy:
        src: ` + composePath + `
        dest: /opt/daos-deploy/docker-compose.yml
      become: yes

    - name: Deploy Docker Compose services
      shell: |
        cd /opt/daos-deploy
        docker-compose up -d
      become: yes
`

	return os.WriteFile(path, []byte(playbook), 0644)
}
