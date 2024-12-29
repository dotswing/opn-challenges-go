package payments

import (
	"fmt"
	"sync"
	"time"

	"github.com/dotswing/opn-challenges-go/internal/donations"
	"github.com/omise/omise-go"
	"github.com/omise/omise-go/operations"
)

type PaymentCharger struct {
	omiseClient   *omise.Client
	requestPerSec int
}

type ChargeProgress struct {
	completed int
	failed    int
}

func NewPaymentCharger(omisePublicKey string, omiseSecretKet string, requestPerSec int) (*PaymentCharger, error) {
	omiseClient, err := omise.NewClient(omisePublicKey, omiseSecretKet)
	if err != nil {
		return nil, err
	}
	return &PaymentCharger{
		omiseClient:   omiseClient,
		requestPerSec: requestPerSec,
	}, nil
}

func (o *PaymentCharger) CreateChargesFromDonations(donationRecords *[]donations.Donation) []omise.Charge {
	var wg sync.WaitGroup
	var chargeResults []omise.Charge
	var mu sync.Mutex
	rateLimiter := time.NewTicker(time.Second / time.Duration(o.requestPerSec))
	defer rateLimiter.Stop()

	totalTasks := 5 //len(*donationRecords)
	progress := make(chan ChargeProgress, totalTasks)
	go func() {
		completed := 0
		failed := 0
		for p := range progress {
			completed += p.completed
			failed += p.failed
			fmt.Printf("\rProgress: %d/%d completed, failed: %d", completed, totalTasks, failed)
		}
	}()

	for _, donationRecord := range (*donationRecords)[0:totalTasks] {
		wg.Add(1)
		<-rateLimiter.C
		go func(donationRecord donations.Donation) {
			defer wg.Done()
			chargeResult, err := o.createChargeFromSingleDonation(donationRecord)
			if err != nil {
				chargeResult = &omise.Charge{}
				progress <- ChargeProgress{failed: 1}
			}
			mu.Lock()
			chargeResults = append(chargeResults, *chargeResult)
			mu.Unlock()
			progress <- ChargeProgress{completed: 1}
		}(donationRecord)
	}

	wg.Wait()

	return chargeResults
}

func SumChargesTHB(chargeResults []omise.Charge) float64 {
	var total float64
	for _, chargeResult := range chargeResults {
		total += float64(chargeResult.Amount)
	}
	return float64(total / 100)
}

func (o *PaymentCharger) createChargeFromSingleDonation(donation donations.Donation) (*omise.Charge, error) {
	defer func() {
		donation = donations.Donation{}
	}()
	token := &omise.Token{}
	err := o.omiseClient.Do(token, &operations.CreateToken{
		Name:            donation.Name,
		Number:          donation.CCNumber,
		ExpirationMonth: time.Month(donation.ExpMonth),
		ExpirationYear:  2026, // Hardcode due to omise api need expire year in the future
		SecurityCode:    donation.CVV,
	})
	if err != nil {
		return nil, err
	}
	charge := &omise.Charge{}
	err = o.omiseClient.Do(charge, &operations.CreateCharge{
		Amount:      int64(donation.AmountSubunits),
		Currency:    "thb",
		Card:        token.ID,
		Description: "Donation for Song-pah-pa",
	})
	if err != nil {
		return nil, err
	}
	return charge, nil
}
