package dns

import (
	"context"
	"fmt"
	"github.com/alibaba/kt-connect/pkg/kt/cluster"
	"github.com/alibaba/kt-connect/pkg/kt/util"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	resolverDir = "/etc/resolver"
	ktResolverPrefix = "kt."
	resolverComment  = "# Generated by KtConnect"
	dnsDefaultPort   = "53"
)

// SetDnsServer set dns server records
func (s *Cli) SetDnsServer(k cluster.KubernetesInterface, dnsServers []string, isDebug bool) error {
	dnsSignal := make(chan error)
	util.CreateDirIfNotExist(resolverDir)
	go func() {
		namespaces, err := k.GetAllNamespaces(context.TODO())
		if err != nil {
			dnsSignal <-err
			return
		}

		preferedDnsInfo := strings.Split(dnsServers[0], ":")
		dnsIp := preferedDnsInfo[0]
		dnsPort := dnsDefaultPort
		if len(preferedDnsInfo) > 1 {
			dnsPort = preferedDnsInfo[1]
		}

		// TODO: read domain suffix from option
		createResolverFile("local", "cluster.local", dnsIp, dnsPort)
		for _, ns := range namespaces.Items {
			createResolverFile(fmt.Sprintf("%s.local", ns.Name), ns.Name, dnsIp, dnsPort)
		}
		dnsSignal <- nil

		defer s.RestoreDnsServer()
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		<-sigCh
	}()
	return <-dnsSignal
}

func createResolverFile(postfix, domain, dnsIp, dnsPort string) {
	resolverFile := fmt.Sprintf("%s/%s%s", resolverDir, ktResolverPrefix, postfix)
	if _, err := os.Stat(resolverFile); err == nil {
		_ = os.Remove(resolverFile)
	}
	resolverContent := fmt.Sprintf("%s\ndomain %s\nnameserver %s\nport %s\n",
		resolverComment, domain, dnsIp, dnsPort)
	if err := ioutil.WriteFile(resolverFile, []byte(resolverContent), 0644); err != nil {
		log.Warn().Err(err).Msgf("Failed to create resolver file of %s", domain)
	}
}

// RestoreDnsServer remove the nameservers added by ktctl
func (s *Cli) RestoreDnsServer() {
	rd, _ := ioutil.ReadDir(resolverDir)
	for _, f := range rd {
		if !f.IsDir() && strings.HasPrefix(f.Name(), ktResolverPrefix) {
			if err := os.Remove(fmt.Sprintf("%s/%s", resolverDir, f.Name())); err != nil {
				log.Warn().Err(err).Msgf("Failed to remove resolver file %s", f.Name())
			}
		}
	}
}