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

	rec := cloudflare.DNSRecord{Name: record, Content: ip}
	_, err = c.API.CreateDNSRecord(context.Background(), c.ZoneId, rec)

	return err
}

func (c CloudflareProvider) RemoveRecord(record string, ip string) error {
	r, err := c.getRecord(record)
	if err != nil {
		log.Fatal(err)
	}

	return c.API.DeleteDNSRecord(context.Background(), c.ZoneId, r.ID)
}

func (c CloudflareProvider) getRecord(record string) (r *cloudflare.DNSRecord, err error) {
	filter := cloudflare.DNSRecord{Name: record}
	recs, err := c.API.DNSRecords(context.Background(), c.ZoneId, filter)
	if len(recs) == 0 || err != nil {
		return nil, err
	}
	r = &recs[0]
	return r, err
}

func NewDnsProvider() (IDnsProvider, error) {
	cloudflareToken := os.Getenv("CF_TOKEN")
	cloudflareAccount := os.Getenv("CF_ACCOUNT_EMAIL")
	cloudflareZone := os.Getenv("CF_ZONE_ID")
	if cloudflareToken != "" && cloudflareAccount != "" && cloudflareZone != "" {
		api, err := cloudflare.New(cloudflareToken, cloudflareAccount)
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
