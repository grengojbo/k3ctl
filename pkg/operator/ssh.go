package operator

import (
	"fmt"
	"time"

	"github.com/appleboy/easyssh-proxy"
	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
)

type SSHOperator struct {
	Config *easyssh.MakeConfig
}

// NewSshConnection New SSH Connection
func (r *SSHOperator) NewSSHOperator(bastion *k3sv1alpha1.BastionNode) {
	r.Config = &easyssh.MakeConfig{
		User: bastion.User,
		// Optional key or Password without either we try to contact your agent SOCKET
		// Password: "password",
		// Paste your source content of private key
		// Key: `-----BEGIN RSA PRIVATE KEY-----
		// .........................
		// -----END RSA PRIVATE KEY-----
		// `,
		Port:    fmt.Sprintf("%d", bastion.SshPort),
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
	r.Config.KeyPath = util.ExpandPath("/home/jbo/.ssh/id_ed25519")
	if len(bastion.SSHAuthorizedKey) > 0 {
		r.Config.KeyPath = util.ExpandPath(bastion.SSHAuthorizedKey)
		log.Debugf("sshKeyPath: %s", r.Config.KeyPath)
	}
	log.Debugf("ssh -i %s %s@%s -p %s", r.Config.KeyPath, r.Config.User, r.Config.Server, r.Config.Port)
}

// Run command on remote machine
//   Example:
func (r *SSHOperator) Run(command string) (done bool, err error) {
	stdOut, stdErr, done, err := r.Config.Run(command, 60*time.Second)
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

// выполнить комманду на удаленном компьютере и вернуть результат как строка
func (r *SSHOperator) Execute(command string) (stdOut string, stdErr string, err error) {
	stdOut, stdErr, _, err = r.Config.Run(command, 60*time.Second)
	return stdOut, stdErr, err
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
					// fmt.Println("out:", outline)
					fmt.Println(outline)
				} else if len(outline) > 0 {
					log.Infoln(outline)
				}
			case errline := <-stderrChan:
				if isPrint && len(errline) > 0 {
					// fmt.Println("err:", errline)
					fmt.Println(errline)
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
