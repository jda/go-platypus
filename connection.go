package platypus // import "go.jona.me/platypus"

import (
	"bufio"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Platypus struct {
	host        string
	addr        *net.TCPAddr
	username    string
	password    string
	logintype   string
	ssl         bool
	sslinsecure bool
	Debug       bool
}

// New creates a new connection to the Platypus WOW API
func New(host string, user string, pass string) (Platypus, error) {
	p := Platypus{
		host:      host,
		username:  user,
		password:  pass,
		logintype: "staff",
	}

	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return p, errors.New(ERR_BAD_HOST)
	}
	p.addr = addr

	return p, nil
}

// NewSSL creates a new connection to the Platypus WOW API using SSL
func NewSSL(host string, user string, pass string, insecure bool) (Platypus, error) {
	p := Platypus{
		host:        host,
		username:    user,
		password:    pass,
		logintype:   "staff",
		ssl:         true,
		sslinsecure: insecure,
	}

	addr, err := net.ResolveTCPAddr("tcp", host)
	if err != nil {
		return p, errors.New(ERR_BAD_HOST)
	}
	p.addr = addr

	return p, nil
}

func (p Platypus) newDataBlock() DataBlock {
	db := DataBlock{
		Protocol:  "Plat",
		Object:    "addusr",
		Username:  p.username,
		Password:  p.password,
		Logintype: p.logintype,
	}

	return db
}

func (p Platypus) getConn() {

}

// Exec calls the WOW API method named in action with a parameter struct.
func (p Platypus) Exec(action string, params interface{}) (DataBlock, error) {
	reply := DataBlock{}

	// prep request
	db := p.newDataBlock()
	db.Action = action
	db.Parameters = params

	b := Body{
		Data: db,
	}

	c := RequestContainer{
		Body: b,
	}

	xmlcmd, err := xml.Marshal(c)
	if err != nil {
		return reply, err
	}

	// package request / calculate header size
	xmlcmd = append([]byte(xml.Header), xmlcmd...)
	prefix := []byte("Content-Length:" + strconv.Itoa(len(xmlcmd)) + "\r\n\r\n")
	rawout := append(prefix, xmlcmd...)

	if p.Debug {
		fmt.Fprintf(os.Stderr, "Request:\n%s\n", rawout)
	}

	var conn net.Conn
	if p.ssl == true {
		sslcfg := tls.Config{InsecureSkipVerify: p.sslinsecure}
		conn, err = tls.Dial("tcp", p.addr.String(), &sslcfg)
	} else {
		conn, err = net.DialTCP("tcp", nil, p.addr)

	}
	if err != nil {
		return reply, err
	}

	// send request
	_, err = conn.Write(rawout)
	if err != nil {
		return reply, err
	}

	// handle response
	header, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return reply, err
	}

	connlen, err := strconv.Atoi(strings.TrimSpace(header[15:]))
	if err != nil {
		return reply, err
	}

	buf := make([]byte, connlen)
	_, err = conn.Read(buf)
	if err != nil {
		return reply, err
	}

	conn.Close()

	if p.Debug {
		fmt.Fprintf(os.Stderr, "Response:\n%s\n", buf)
	}

	var res = ResponseContainer{}
	err = xml.Unmarshal(buf, &res)
	if err != nil {
		return reply, err
	}

	reply = res.Body.Data

	return reply, nil
}
