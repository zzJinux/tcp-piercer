package pnet

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/coreos/go-iptables/iptables"
	"go.uber.org/multierr"
)

type Dialer struct {
	net.Dialer
	ServicePort int // the port that some process is listening on
}

type cleanableConn struct {
	net.Conn
	cleanup func() error
}

func (c *cleanableConn) Close() error {
	return multierr.Append(
		c.cleanup(),
		c.Conn.Close(),
	)
}

func (d *Dialer) DialContext(ctx context.Context, network, address string) (conn net.Conn, err error) {
	// I decided to *shadow* the embedded method instead of defining a new one without a `network` argument
	if network != "tcp" {
		panic(fmt.Sprint("Not implemented for this network: ", network))
	}

	localPort, err := getAvailablePort()
	if err != nil {
		return nil, fmt.Errorf("no port available: %w", err)
	}
	serverAddr, err := net.ResolveTCPAddr("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("resolve error: %w", err)
	}

	// configure iptables
	//
	// # [Agent Client] ---> [Agent Server]
	// iptables -t nat -A POSTROUTING -p tcp --sport LOCAL_PORT -d SERVER_IP --dport SERVER_PORT -j SNAT --to-source :SERVICE_PORT

	// # [Agent Client] <--- [Agent Server]
	// iptables -t nat -A PREROUTING -p tcp --dport SERVICE_PORT -s SERVER_IP --sport SERVER_PORT -j DNAT --to-destination :LOCAL_PORT

	ipt, err := iptables.New()
	if err != nil {
		return nil, fmt.Errorf("iptables unavailable: %w", err)
	}

	type iptRuleSpec []string
	outgoingRuleSpec := iptRuleSpec{
		"-p", "tcp",
		"--sport", strconv.Itoa(localPort),
		"-d", serverAddr.IP.String(), "--dport", strconv.Itoa(serverAddr.Port),
		"-j", "SNAT", "--to-source", fmt.Sprint(":", d.ServicePort),
	}

	incomingRuleSpec := iptRuleSpec{
		"-p", "tcp",
		"--dport", strconv.Itoa(d.ServicePort),
		"-s", serverAddr.IP.String(), "--sport", strconv.Itoa(serverAddr.Port),
		"-j", "SNAT", "--to-source", fmt.Sprint(":", localPort),
	}

	// Append iptable rules
	if err = ipt.Append("nat", "POSTROUTING", outgoingRuleSpec...); err != nil {
		return nil, err
	}
	if err = ipt.Append("nat", "POSTROUTING", incomingRuleSpec...); err != nil {
		return nil, err
	}

	undoIPTActions := func() (err error) {
		// Delete iptable rules
		if err = ipt.Delete("nat", "POSTROUTING", outgoingRuleSpec...); err != nil {
			return
		}
		if err = ipt.Delete("nat", "POSTROUTING", incomingRuleSpec...); err != nil {
			return
		}
		return nil
	}

	// defer a clean-up
	defer func() {
		if err != nil {
			err = multierr.Append(err, undoIPTActions())
		}
	}()

	// This is the actual conection to communicate with the server
	d.LocalAddr, err = net.ResolveTCPAddr("tcp", fmt.Sprint(":", localPort))
	if err != nil {
		return nil, err
	}
	conn, err = d.Dialer.DialContext(ctx, network, address)
	if err != nil {
		return conn, err
	}

	return &cleanableConn{conn, undoIPTActions}, nil
}

func getAvailablePort() (port int, err error) {
	// let the OS determine the available port
	l, err := net.Listen("tcp", "")
	if err != nil {
		// extremely rare case
		panic(fmt.Errorf("port exhaustion?: %w", err))
	}
	defer func() {
		err = l.Close()
	}()
	return l.Addr().(*net.TCPAddr).Port, nil
}
