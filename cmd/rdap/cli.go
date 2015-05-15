package main

import (
	"io"
	"net"
	"net/http"
	"strconv"

	"github.com/registrobr/rdap-client/bootstrap"
	"github.com/registrobr/rdap-client/client"
	"github.com/registrobr/rdap-client/cmd/rdap/Godeps/_workspace/src/github.com/davecgh/go-spew/spew"
	"github.com/registrobr/rdap-client/output"
)

type CLI struct {
	uris       []string
	httpClient *http.Client
	bootstrap  *bootstrap.Client
	wr         io.Writer
}

type handler func(object string) (bool, error)

func (c *CLI) asn() handler {
	return func(object string) (bool, error) {
		uris := c.uris

		if asn, err := strconv.ParseUint(object, 10, 32); err == nil {
			if c.bootstrap != nil {
				var err error
				uris, err = c.bootstrap.ASN(asn)

				if err != nil {
					return true, err
				}
			}

			r, err := client.NewClient(uris, c.httpClient).ASN(asn)

			if err != nil {
				return true, err
			}

			spew.Dump(r)

			return true, nil
		}

		return false, nil
	}
}

func (c *CLI) ipnetwork() handler {
	return func(object string) (bool, error) {
		uris := c.uris

		if _, cidr, err := net.ParseCIDR(object); err == nil {
			if c.bootstrap != nil {
				var err error
				uris, err = c.bootstrap.IPNetwork(cidr)

				if err != nil {
					return true, err
				}
			}

			r, err := client.NewClient(uris, c.httpClient).IPNetwork(cidr)

			if err != nil {
				return true, err
			}

			spew.Dump(r)

			return true, nil
		}

		return false, nil
	}
}

func (c *CLI) domain() handler {
	return func(object string) (bool, error) {
		uris := c.uris

		if c.bootstrap != nil {
			var err error
			uris, err = c.bootstrap.Domain(object)

			if err != nil {
				return true, err
			}
		}

		r, err := client.NewClient(uris, c.httpClient).Domain(object)

		if err != nil {
			return true, err
		}

		if err := output.PrintDomain(r, c.wr); err != nil {
			return true, err
		}

		return false, nil
	}
}
