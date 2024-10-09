package knotfree

import (
	"context"
	"net"
	"strings"

	"github.com/awootton/knotfreeiot/iot"
	"github.com/awootton/knotfreeiot/packets"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"

	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

// Define log to be a logger with the plugin name in it. This way we can just use log.Info and
// friends to log.
var log = clog.NewWithPlugin("knotfree")

// knotfree is an example plugin to show how to write a plugin.
type Knotfree struct {
	Next          plugin.Handler
	instanceCount int

	sc *iot.ServiceContactTcp
}

// init registers this plugin.
func init() {
	plugin.Register("knotfree", setup)
}

// turn this on after the knotfree plugin is ready aka SystenContact is up.
func (e Knotfree) Ready() bool {

	log.Info("knotfree is ready", e.instanceCount)

	return true
}

// Name implements the Handler interface.
func (e Knotfree) Name() string {
	// log.Info("knotfree Name")
	return "knotfree"
}

// ResponsePrinter wrap a dns.ResponseWriter and will write example to standard output when WriteMsg is called.
type ResponsePrinter struct {
	dns.ResponseWriter
}

// NewResponsePrinter returns ResponseWriter.
func NewResponsePrinter(w dns.ResponseWriter) *ResponsePrinter {
	return &ResponsePrinter{ResponseWriter: w}
}

// WriteMsg calls the underlying ResponseWriter's WriteMsg method and prints "example" to standard output.
func (r *ResponsePrinter) WriteMsg(res *dns.Msg) error {
	log.Info("knotfree ResponsePrinter WriteMsg")
	return r.ResponseWriter.WriteMsg(res)
}

// ServeDNS implements the plugin.Handler interface. This method gets called when knotfree is used
// in a Server.
func (e Knotfree) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {

	// log.Info("Received response")

	state := request.Request{W: w, Req: r}

	qname := state.Name()
	qtype := state.Type()

	// log.Info("Received request for ", qname, " ", qtype) // eg a-person-channel.iot. A
	subscriptionName := strings.TrimRight(qname, ".")
	// it could be a name like get.option.a.get-unix-time.iot
	// we remove the sub names
	parts := strings.Split(subscriptionName, ".")
	if len(parts) > 2 {
		parts = parts[len(parts)-2:]
	}
	subscriptionName = strings.Join(parts, "_")

	log.Info("Received request for ", subscriptionName, " ", qtype) // eg a-person-channel.iot

	// let's get it from the service contact
	// TODO: test add subtypes
	command := "get option " + strings.ToUpper(qtype) // eg get option A
	cmd := packets.Lookup{}
	cmd.Address.FromString(subscriptionName)
	cmd.SetOption("cmd", []byte(command))
	// send it
	replyPacket, err := e.sc.Get(&cmd)
	if err != nil {
		log.Error("knotfree failed to get from service contact", err)
		return dns.RcodeServerFailure, err
	}
	log.Info("knotfree returned from service contact", replyPacket.Sig())
	sendPacket, ok := replyPacket.(*packets.Send)
	if !ok {
		log.Error("knotfree failed to get 'send' from service contact", err, replyPacket.Sig())
		return dns.RcodeServerFailure, err
	}
	log.Info("knotfree returned message", string(sendPacket.Payload))

	// Create a new response message. We use the message from the request in our response.

	// TODO: do the rest of the types
	rr := new(dns.A)
	rr.Hdr = dns.RR_Header{Name: qname, Rrtype: dns.TypeA, Class: dns.ClassINET}
	str := string(sendPacket.Payload)
	rr.A = net.ParseIP(str).To4()

	answers := []dns.RR{}
	answers = append(answers, rr)

	m := new(dns.Msg)
	m.SetReply(r)
	m.Authoritative = true
	m.Answer = answers

	w.WriteMsg(m)
	return dns.RcodeSuccess, nil

	// Export metric with the server label set to the current server handling the request.
	// requestCount.WithLabelValues(metrics.WithServer(ctx)).Inc()

	// Call next plugin (if any).
	// return plugin.NextOrFailure(e.Name(), e.Next, ctx, w, r)
}
