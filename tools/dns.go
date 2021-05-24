package tools

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

var ErrorDnsProviderNotFound = errors.New("DNS provider cannot be found")

type DnsProviderType string

const (
	DnsCloudflare DnsProviderType = "cf"
)

type IDnsProvider interface {
	AddRecord(record string, ip string) error
	RemoveRecord(record string, ip string) error
}

type CloudflareProvider struct {
	ZoneId string
	API    *cloudflare.API
}

func (c CloudflareProvider) AddRecord(record string, ip string) error {
	r, err := c.getRecord(record)
	if r != nil || err != nil {
		return err
	}

	proxied := true
	rec := cloudflare.DNSRecord{Type: "A", Name: record, Content: ip, Proxied: &proxied, ZoneID: c.ZoneId}
	_, err = c.API.CreateDNSRecord(context.Background(), c.ZoneId, rec)

	return err
}

func (c CloudflareProvider) RemoveRecord(record string, ip string) error {
	r, err := c.getRecord(record)
	if r == nil || err != nil {
		return err
	}
	log.Printf("DNS record found for deletion %v", r)

	return c.API.DeleteDNSRecord(context.Background(), c.ZoneId, r.ID)
}

func (c CloudflareProvider) getRecord(record string) (r *cloudflare.DNSRecord, err error) {
	filter := cloudflare.DNSRecord{Name: record, ZoneID: c.ZoneId}
	recs, err := c.API.DNSRecords(context.Background(), c.ZoneId, filter)
	if len(recs) == 0 || err != nil {
		return nil, err
	}
	r = &recs[0]
	return r, err
}

func NewDnsProvider() (IDnsProvider, error) {
	cloudflareToken := os.Getenv("CF_TOKEN")
	cloudflareZone := os.Getenv("CF_ZONE_ID")
	if cloudflareToken != "" && cloudflareZone != "" {
		api, err := cloudflare.NewWithAPIToken(cloudflareToken)
		if err != nil {
			log.Fatal(err)
		}

		return CloudflareProvider{
			ZoneId: cloudflareZone,
			API:    api,
		}, nil
	}

	return nil, ErrorDnsProviderNotFound
}
