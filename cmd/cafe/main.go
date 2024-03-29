package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudflare/cloudflare-go"
	"github.com/google/go-jsonnet"
)

func main() {
	// Construct a new API object using a global API key
	//api, err := cloudflare.New(os.Getenv("CLOUDFLARE_API_KEY"), os.Getenv("CLOUDFLARE_API_EMAIL"))
	// alternatively, you can use a scoped API token

	api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}

	// Most API calls require a Context
	ctx := context.Background()

	zones, err := api.ListZones(ctx)
	if err != nil {
		log.Fatalf("get zones: %s", err)
	}

	zoneIds := make(map[string]string, 0)

	cfRecords := make([]cloudflare.DNSRecord, 0)

	for _, zone := range zones {
		zoneIds[zone.Name] = zone.ID

		records, err := api.DNSRecords(ctx, zone.ID, cloudflare.DNSRecord{})

		if err != nil {
			log.Fatalf("get record for zone %s: %s", zone.Name, err)
		}

		cfRecords = append(cfRecords, records...) // Get All Records
	}

	vm := jsonnet.MakeVM()

	dfRecords := make([]cloudflare.DNSRecord, 0)

	zoneRoot := os.Getenv("ZONES_DIR")
	if zoneRoot == "" {
		zoneRoot = "./zones"
	}
	err = filepath.Walk(zoneRoot, func(path string, info fs.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, ".jsonnet") {
			jsonData, err := vm.EvaluateFile(path)
			if err != nil {
				return nil
			}

			records := make([]cloudflare.DNSRecord, 0)

			err = json.Unmarshal([]byte(jsonData), &records)

			if err != nil {
				return err
			}

			for _, record := range records {
				if _, ok := zoneIds[record.ZoneName]; !ok {
					fmt.Printf("Error at %s\nZone `%s` is not managed by cafe! Is it corrected?\n", path, record.ZoneName)
					os.Exit(1)
				}
			}

			dfRecords = append(dfRecords, records...)

		}

		return nil
	})

	if err != nil {
		log.Fatalf("walk dir error %s", err)
	}

	cnames := SetOf[string]()

	managed := SetOf[string]()
	stored := SetOf[string]()

	deleting := make([]cloudflare.DNSRecord, 0)
	adding := make([]cloudflare.DNSRecord, 0)

	defined := SetOf[string]()
	deprecate := SetOf[string]()

	for _, record := range dfRecords {
		managed.Add(nt(record))
	}

	for _, record := range dfRecords {
		if record.Type == "CNAME" {
			if cnames.Has(record.Name) {
				log.Fatalf("ERROR: cname is duplicated: %s", record.Name)
			}
			cnames.Add(record.Name)
			defined.Add(hash(record))
		} else if record.Type == "DEPRECATED" {
			deprecate.Add(record.Name)
		} else {
			defined.Add(hash(record))
		}
	}

	for _, record := range cfRecords {
		if managed.Has(nt(record)) {
			h := hash(record)
			if !defined.Has(h) || deprecate.Has(record.Name) {
				deleting = append(deleting, record)
			}
			stored.Add(h)
		}
	}

	for _, record := range dfRecords {
		if record.Type != "DEPRECATED" {
			if !stored.Has(hash(record)) {
				adding = append(adding, record)
			}
		}
	}

	if len(deleting) > 0 {
		fmt.Printf("Those records will be deleted:\n")

		fmt.Printf("%-12s%-8s%-6s%-24s%-20s%s\n", "ZONE", "TYPE", "TTL", "NAME", "CONTENT", "PROXY")
		for _, record := range deleting {
			fmt.Printf("%-12s%-8s%-6d%-24s%-20s%v\n", record.ZoneName, record.Type, record.TTL, record.Name, record.Content, *record.Proxied)
		}
	}

	if len(adding) > 0 {
		fmt.Printf("\nThose records will be created:\n")

		fmt.Printf("%-12s%-8s%-6s%-24s%s-10s%s\n", "ZONE", "TYPE", "TTL", "NAME", "CONTENT", "PROXY")
		for _, record := range adding {
			fmt.Printf("%-12s%-8s%-6d%-24s%-20s%v\n", record.ZoneName, record.Type, record.TTL, record.Name, record.Content, *record.Proxied)
		}
	}

	if os.Getenv("CAFE_CONFIRM") == "yes" {

		if len(deleting) > 0 || len(adding) > 0 {
			fmt.Printf("\nDo update records...\n\n")
		}

		for _, record := range deleting {
			fmt.Printf("deleting %-12s%-8s%-6d%-24s%s...\n", record.ZoneName, record.Type, record.TTL, record.Name, record.Content)
			err := api.DeleteDNSRecord(ctx, zoneIds[record.ZoneName], record.ID)
			if err != nil {
				log.Fatal(err)
			}
		}

		for _, record := range adding {
			fmt.Printf("creating %-12s%-8s%-6d%-24s%s... ", record.ZoneName, record.Type, record.TTL, record.Name, record.Content)

			res, err := api.CreateDNSRecord(ctx, zoneIds[record.ZoneName], record)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%v\n", res.Success)
		}
	}

}

func hash(r cloudflare.DNSRecord) string {

	var s [16]byte

	if r.Type == "MX" {
		s = md5.Sum([]byte(fmt.Sprintf("%s-%s-%d-%s-%s-%d", r.ZoneName, r.Type, r.TTL, r.Name, r.Content, *r.Priority)))
	} else if r.Type == "A" || r.Type == "CNAME" || r.Type == "DEPRECATED" {
		s = md5.Sum([]byte(fmt.Sprintf("%s-%s-%d-%s-%s-%v", r.ZoneName, "X", r.TTL, r.Name, r.Content, *r.Proxied)))
	} else {
		s = md5.Sum([]byte(fmt.Sprintf("%s-%s-%d-%s-%s", r.ZoneName, r.Type, r.TTL, r.Name, r.Content)))
	}

	return hex.EncodeToString(s[:])
}

func nt(r cloudflare.DNSRecord) string {
	if r.Type == "A" || r.Type == "CNAME" || r.Type == "DEPRECATED" {
		return fmt.Sprintf("%s-X", r.Name)
	}
	return fmt.Sprintf("%s-%s", r.Name, r.Type)
}

type Set[T comparable] map[T]struct{}

func SetOf[T comparable]() Set[T] {
	return make(Set[T], 0)
}

func (s Set[T]) Add(item T) bool {
	if _, ok := s[item]; ok {
		return false
	}
	s[item] = struct{}{}
	return true
}

func (s Set[T]) Has(item T) bool {
	_, ok := s[item]
	return ok
}
