package cert

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"net"
	"os"
	"strconv"
	"text/template"
	"time"
)

var timeout = 5 * time.Second
var expires = 7 * 24 * time.Hour

type HostInfo struct {
	Host  string
	Port  int
	Certs []*x509.Certificate
}

func (h *HostInfo) GetCerts(timeout time.Duration) error {
	// log.Printf("connecting to %s:%d", h.Host, h.Port)
	dialer := &net.Dialer{Timeout: timeout}
	conn, err := tls.DialWithDialer(
		dialer,
		"tcp",
		h.Host+":"+strconv.Itoa(h.Port),
		&tls.Config{
			InsecureSkipVerify: true,
		})
	if err != nil {
		return err
	}

	defer conn.Close()

	if err := conn.Handshake(); err != nil {
		return err
	}

	pc := conn.ConnectionState().PeerCertificates
	h.Certs = make([]*x509.Certificate, 0, len(pc))
	for _, cert := range pc {
		if cert.IsCA {
			continue
		}
		h.Certs = append(h.Certs, cert)
	}

	return nil
}

// show the certificates of the host
func ShowCerts(hosts []HostInfo, mode string) (err error) {

	switch mode {
	case "json":
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		enc.SetEscapeHTML(false)
		err = enc.Encode(&hosts)
		// errs.Push(err)

	case "text":
		t := template.Must(template.New("").Parse(`
{{- range . -}}
Host: {{ .Host }}:{{ .Port }}
Certs:
	{{ range .Certs -}}
	Issuer: {{ .Issuer.CommonName }}
	Subject: {{ .Subject.CommonName }}
	Not Before: {{ .NotBefore.Format "Jan 2, 2006 3:04 PM" }}
	Not After: {{ .NotAfter.Format "Jan 2, 2006 3:04 PM" }}
	DNS names: {{ range .DNSNames }}{{ . }} {{ end }}
{{ end }}
{{ end -}}
			`))
		err = t.Execute(os.Stdout, &hosts)
		// errs.Push(err)

	case "none":
	}
	return err
}
