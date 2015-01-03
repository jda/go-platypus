package platypus

import (
	"bufio"
	"encoding/xml"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

type Platypus struct {
	host      string
	addr      *net.TCPAddr
	username  string
	password  string
	logintype string
	Debug     bool
}

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

	// send request
	conn, err := net.DialTCP("tcp", nil, p.addr)
	if err != nil {
		return reply, err
	}

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
