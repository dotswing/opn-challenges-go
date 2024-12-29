package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"sort"

	"github.com/dotswing/opn-challenges-go/internal/donations"
	"github.com/dotswing/opn-challenges-go/internal/payments"
	"github.com/dotswing/opn-challenges-go/pkg/numberutils"
)

func main() {
	requestPerSec := flag.Int("requestPerSec", runtime.NumCPU(), "Rate limiting request/sec. Default to no. of CPU")
	file := flag.String("file", "fng.1000.csv.rot128", "File path")
	flag.Parse()
	fmt.Printf("Rate limiting %d request/sec\n", *requestPerSec)
	fmt.Printf("File: %s\n", *file)

	donationRecords, err := donations.GetDecryptedCSVFromFile(*file)
	defer func() {
		*donationRecords = nil
	}()
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	paymentCharger, err := payments.NewPaymentCharger(os.Getenv("OMISE_PUBLIC_KEY"), os.Getenv("OMISE_SECRET_KEY"), *requestPerSec)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return
	}
	fmt.Println("performing donations...")
	chargedResults := paymentCharger.CreateChargesFromDonations(donationRecords)
	defer func() {
		chargedResults = nil
	}()

	totalReceivedTHB := donations.SumDonationsTHB(donationRecords)
	successfullyDonatedTHB := payments.SumChargesTHB(chargedResults)
	averagePerPerson := totalReceivedTHB / float64(len(*donationRecords))

	sort.Slice(chargedResults, func(i, j int) bool {
		return chargedResults[i].Amount > chargedResults[j].Amount
	})
	fmt.Printf("\ntotal received:       THB %s\n", numberutils.FormatFloat(totalReceivedTHB))
	fmt.Printf("successfully donated: THB %s\n", numberutils.FormatFloat(successfullyDonatedTHB))
	fmt.Printf("faulty donation:      THB %s\n", numberutils.FormatFloat(totalReceivedTHB-successfullyDonatedTHB))
	fmt.Println()
	fmt.Printf("average per person:   THB %s\n", numberutils.FormatFloat(averagePerPerson))
	fmt.Println("\ntop donors (only successfully donated):")
	if len(chargedResults) > 0 && chargedResults[0].Status != "" {
		for _, charge := range chargedResults[0:3] {
			fmt.Printf("Name: %s, Amount: %s\n", charge.Card.Name, numberutils.FormatFloat(float64(charge.Amount)/100))
		}
	}
	fmt.Println("\ntop donors (all including failed):")
	sort.Slice(*donationRecords, func(i, j int) bool {
		return int64((*donationRecords)[i].AmountSubunits) > int64((*donationRecords)[j].AmountSubunits)
	})
	for _, donation := range (*donationRecords)[0:3] {
		fmt.Printf("Name: %s, Amount: %s\n", donation.Name, numberutils.FormatFloat(float64(donation.AmountSubunits)/100))
	}
}
