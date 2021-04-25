package operator

import (
	"fmt"
	"time"

	"github.com/appleboy/easyssh-proxy"
	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	log "github.com/sirupsen/logrus"
)

// type SSHConfig struct {
// 	User    string
// 	Server  string
// 	Port    string
// 	KeyPath string
// }
// type SSHOperator struct {
// 	Config *easyssh.MakeConfig
// 	// ProxyConfig SSHConfig
// }

type SSHOperator struct {
	Config *easyssh.MakeConfig
	// ProxyConfig SSHConfig
}

// NewSshConnection New SSH Connection
func (r *SSHOperator) NewSSHOperator(bastion *k3sv1alpha1.BastionNode) {
	r.Config = &easyssh.MakeConfig{
		User: bastion.User,
		// User:   "ubuntu",
		// Server: bastion.Address,
		// Optional key or Password without either we try to contact your agent SOCKET
		// Password: "password",
		// Paste your source content of private key
		// Key: `-----BEGIN RSA PRIVATE KEY-----
		// .........................
		// -----END RSA PRIVATE KEY-----
		// `,
		KeyPath: "/home/jbo/.ssh/id_rsa",
		Port:    fmt.Sprintf("%d", bastion.SshPort),
		// KeyPath: bastion.SSHAuthorizedKey,
		Timeout: 60 * time.Second,

		// Parse PrivateKey With Passphrase
		// Passphrase: "XXXX",

		// Optional fingerprint SHA256 verification
		// Get Fingerprint: ssh.FingerprintSHA256(key)
		// Fingerprint: "SHA256:................E"

		// Enable the use of insecure ciphers and key exchange methods.
		// This enables the use of the the following insecure ciphers and key exchange methods:
		// - aes128-cbc
		// - aes192-cbc
		// - aes256-cbc
		// - 3des-cbc
		// - diffie-hellman-group-exchange-sha256
		// - diffie-hellman-group-exchange-sha1
		// Those algorithms are insecure and may allow plaintext data to be recovered by an attacker.
		// UseInsecureCipher: true,
	}
	r.Config.Server = bastion.Address
}

func NewSshConnection(bastion *k3sv1alpha1.BastionNode) (ssh *easyssh.MakeConfig) {
	// https://github.com/appleboy/easyssh-proxy
	ssh = &easyssh.MakeConfig{
		User: bastion.User,
		// User:   "ubuntu",
		// Server: bastion.Address,
		// Optional key or Password without either we try to contact your agent SOCKET
		// Password: "password",
		// Paste your source content of private key
		// Key: `-----BEGIN RSA PRIVATE KEY-----
		// .........................
		// -----END RSA PRIVATE KEY-----
		// `,
		Port: fmt.Sprintf("%d", bastion.SshPort),
		// KeyPath: bastion.SSHAuthorizedKey,
		Timeout: 60 * time.Second,

		// Parse PrivateKey With Passphrase
		// Passphrase: "XXXX",

		// Optional fingerprint SHA256 verification
		// Get Fingerprint: ssh.FingerprintSHA256(key)
		// Fingerprint: "SHA256:................E"

		// Enable the use of insecure ciphers and key exchange methods.
		// This enables the use of the the following insecure ciphers and key exchange methods:
		// - aes128-cbc
		// - aes192-cbc
		// - aes256-cbc
		// - 3des-cbc
		// - diffie-hellman-group-exchange-sha256
		// - diffie-hellman-group-exchange-sha1
		// Those algorithms are insecure and may allow plaintext data to be recovered by an attacker.
		// UseInsecureCipher: true,
	}
	ssh.Server = bastion.Address
	ssh.KeyPath = "/home/jbo/.ssh/id_rsa"

	log.Debugf("ssh -i %s %s@%s:%s", ssh.KeyPath, ssh.User, ssh.Server, ssh.Port)
	return ssh
}

// Run command on remote machine
//   Example:
func (r *SSHOperator) Run(command string) (done bool, err error) {
	stdOut, stdErr, done, err := r.Config.Run(command, 60*time.Second)
	// stdout, stderr, done, err := ssh.Run("ls -al", 60*time.Second)
	// // Handle errors
	// if err != nil {
	// 	log.Fatalln("Can't run remote command: " + err.Error())
	// } else {
	// 	log.Infoln("don is :", done, "stdout is :", stdout, ";   stderr is :", stderr)
	// }
	if len(stdOut) > 0 {
		log.Debugln("===== stdOut ======")
		log.Debugf("%v", stdOut)
		log.Debugln("===================")
	}
	if len(stdErr) > 0 {
		log.Errorln("===== stdErr ======")
		log.Errorf("%v", stdErr)
		log.Errorln("===================")
	}
	return done, err
}

// Stream returns one channel that combines the stdout and stderr of the command
// as it is run on the remote machine, and another that sends true when the
// command is done. The sessions and channels will then be closed.
//  isPrint - выводить результат на экран или в лог
func (r *SSHOperator) Stream(command string, isPrint bool) {
	// Call Run method with command you want to run on remote server.
	stdoutChan, stderrChan, doneChan, errChan, err := r.Config.Stream(command, 60*time.Second)
	// Handle errors
	if err != nil {
		log.Fatalln("Can't run remote command: " + err.Error())
	} else {
		// read from the output channel until the done signal is passed
		isTimeout := true
	loop:
		for {
			select {
			case isTimeout = <-doneChan:
				break loop
			case outline := <-stdoutChan:
				if isPrint && len(outline) > 0 {
					fmt.Println("out:", outline)
				} else if len(outline) > 0 {
					log.Infoln(outline)
				}
			case errline := <-stderrChan:
				if isPrint && len(errline) > 0 {
					fmt.Println("err:", errline)
				} else if len(errline) > 0 {
					log.Errorln(errline)
				}
			case err = <-errChan:
			}
		}

		// get exit code or command error.
		if err != nil {
			log.Errorln("Error: " + err.Error())
		}

		// command time out
		if !isTimeout {
			log.Errorln("Error: command timeout")
		}
	}
}

func MakeSsshConfig(config *easyssh.MakeConfig) {

	// ssh := &easyssh.MakeConfig{
	// 	User:   "appleboy",
	// 	Server: "example.com",
	// 	// Optional key or Password without either we try to contact your agent SOCKET
	// 	// Password: "password",
	// 	// Paste your source content of private key
	// 	// Key: `-----BEGIN RSA PRIVATE KEY-----
	// 	// MIIEpAIBAAKCAQEA4e2D/qPN08pzTac+a8ZmlP1ziJOXk45CynMPtva0rtK/RB26
	// 	// 7XC9wlRna4b3Ln8ew3q1ZcBjXwD4ppbTlmwAfQIaZTGJUgQbdsO9YA==
	// 	// -----END RSA PRIVATE KEY-----
	// 	// `,
	// 	KeyPath: "/Users/username/.ssh/id_rsa",
	// 	Port:    "22",
	// 	Timeout: 60 * time.Second,

	// 	// Parse PrivateKey With Passphrase
	// 	Passphrase: "1234",

	// 	// Optional fingerprint SHA256 verification
	// 	// Get Fingerprint: ssh.FingerprintSHA256(key)
	// 	// Fingerprint: "SHA256:mVPwvezndPv/ARoIadVY98vAC0g+P/5633yTC4d/wXE"

	// 	// Enable the use of insecure ciphers and key exchange methods.
	// 	// This enables the use of the the following insecure ciphers and key exchange methods:
	// 	// - aes128-cbc
	// 	// - aes192-cbc
	// 	// - aes256-cbc
	// 	// - 3des-cbc
	// 	// - diffie-hellman-group-exchange-sha256
	// 	// - diffie-hellman-group-exchange-sha1
	// 	// Those algorithms are insecure and may allow plaintext data to be recovered by an attacker.
	// 	// UseInsecureCipher: true,
	// }
	// log.Debugln("ssh: ", ssh)
	// return nil, nil
}
