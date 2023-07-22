package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

type dnsRecords struct {
	hasMX       bool
	hasSPF      bool
	hasDMARC    bool
	spfRecord   string
	dmarcRecord string
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Printf("Domain, hasMX, hasSPF, sprRecord, hasDMARC, dmarcRecord, isValidEmailSetup\n")

	for scanner.Scan() {
		checkDomain(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error: could not read from input: %v\n", err)
	}
}

func checkDomain(domain string) {
	var wg sync.WaitGroup
	record := dnsRecords{}
	wg.Add(3)

	go func() {
		defer wg.Done()
		record.checkMX(domain)
	}()

	go func() {
		defer wg.Done()
		record.checkSPF(domain)
	}()

	go func() {
		defer wg.Done()
		record.checkDMARC(domain)
	}()

	wg.Wait()

	validEmailSetup := "No"
	if record.hasMX && record.hasSPF && record.hasDMARC {
		validEmailSetup = "Yes"
	}

	fmt.Printf("\nDomain: %v\nHas MX Records: %v\nHas SPF Record: %v\nSPF Record: %v\nHas DMARC Record: %v\nDMARC Record: %v\nIs Valid Email Setup: %v\n", domain, record.hasMX, record.hasSPF, record.spfRecord, record.hasDMARC, record.dmarcRecord, validEmailSetup)
}

func (r *dnsRecords) checkMX(domain string) {
	mxRecords, err := net.LookupMX(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	if len(mxRecords) > 0 {
		r.hasMX = true
	}
}

func (r *dnsRecords) checkSPF(domain string) {
	txtRecord, err := net.LookupTXT(domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	for _, record := range txtRecord {
		if strings.HasPrefix(record, "v=spf1") {
			r.hasSPF = true
			r.spfRecord += record
			break
		}
	}
}

func (r *dnsRecords) checkDMARC(domain string) {
	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		log.Printf("Error: %v\n", err)
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			r.hasDMARC = true
			r.dmarcRecord = record
			break
		}
	}
}
