package donations

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"

	"github.com/dotswing/opn-challenges-go/pkg/cipher"
	"github.com/dotswing/opn-challenges-go/pkg/fileutils"
)

type Donation struct {
	Name           string
	AmountSubunits int
	CCNumber       string
	CVV            string
	ExpMonth       int
	ExpYear        int
}

func GetDecryptedCSVFromFile(filePath string) (*[]Donation, error) {
	donationsDataBuffer := &bytes.Buffer{}
	defer func() {
		donationsDataBuffer = nil
	}()

	data, err := fileutils.ReadFileToBytes(filePath)
	byteCount := len(data)
	if err != nil {
		fmt.Printf("Error ReadFileToBytes")
		return nil, err
	}

	writer, err := cipher.NewRot128Writer(donationsDataBuffer, byteCount)
	if err != nil {
		fmt.Printf("Error NewRot128Writer")
		return nil, err
	}
	_, err = writer.Write(data)
	if err != nil {
		fmt.Printf("Error NewRot128Writer write")
		return nil, err
	}

	reader := csv.NewReader(donationsDataBuffer)
	csvRecords, err := reader.ReadAll()
	defer func() {
		csvRecords = nil
	}()
	if err != nil {
		fmt.Printf("Error reading CSV: %v\n", err)
		return nil, err
	}
	donations := make([]Donation, len(csvRecords)-1)
	for i, record := range csvRecords[1:] {
		amountSubunits, _ := strconv.Atoi(record[1])
		expMonth, _ := strconv.Atoi(record[4])
		expYear, _ := strconv.Atoi(record[5])
		donations[i] = Donation{
			Name:           record[0],
			AmountSubunits: amountSubunits,
			CCNumber:       record[2],
			CVV:            record[3],
			ExpMonth:       expMonth,
			ExpYear:        expYear,
		}
		amountSubunits = 0
		expMonth = 0
		expYear = 0
	}
	return &donations, nil
}

func SumDonationsTHB(donations *[]Donation) float64 {
	var total float64
	for _, donation := range *donations {
		total += float64(donation.AmountSubunits)
	}
	return float64(total / 100)
}
