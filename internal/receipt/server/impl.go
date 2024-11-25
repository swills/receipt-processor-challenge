package receiptserver

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../../../tools/config.server.yaml ../../../api.yml

import (
	"encoding/json"
	"log/slog"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
	openapiTypes "github.com/oapi-codegen/runtime/types"
)

type Server struct{}

type PostReceiptsProcessResponse struct {
	ID string `json:"id"`
}

type GetReceiptsIDPoints struct {
	Points uint64 `json:"points"`
}

var pointsMap = map[uuid.UUID]uint64{}

func NewServer() Server {
	return Server{}
}

// GetReceiptsIdPoints handles HTTP 'GET' request for points with given ID
//
// the revive and stylecheck linters are disabled because they want "Id" to be "ID", but "Id" is used in generated code
//
//nolint:revive,stylecheck
func (s Server) GetReceiptsIdPoints(writer http.ResponseWriter, _ *http.Request, reqReceiptId string) {
	var receiptUUID uuid.UUID

	var respJSON []byte

	var err error

	receiptUUID, err = uuid.Parse(reqReceiptId)
	if err != nil {
		slog.Error("error parsing receiptUUID", "err", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	points, valid := pointsMap[receiptUUID]
	if !valid {
		slog.Error("receipt not found", "receiptUUID", receiptUUID.String())
		writer.WriteHeader(http.StatusNotFound)

		return
	}

	resp := GetReceiptsIDPoints{Points: points}

	respJSON, err = json.Marshal(resp)
	if err != nil {
		slog.Error("error encoding response", "err", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	respJSON = append(respJSON, '\n')

	writer.WriteHeader(http.StatusOK)

	_, err = writer.Write(respJSON)
	if err != nil {
		slog.Error("error sending response", "err", err)

		return
	}
}

func (Server) PostReceiptsProcess(writer http.ResponseWriter, request *http.Request) {
	var respJSON []byte

	var receipt Receipt

	var err error

	var receiptUUID uuid.UUID

	receiptDecoder := json.NewDecoder(request.Body)

	err = receiptDecoder.Decode(&receipt)
	if err != nil {
		slog.Error("error parsing request", "err", err)
		writer.WriteHeader(http.StatusBadRequest)

		return
	}

	points := calculatePoints(receipt)

	receiptUUID, err = uuid.NewRandom()

	if err != nil {
		slog.Error("error generating uuid", "err", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	pointsMap[receiptUUID] = points

	resp := PostReceiptsProcessResponse{ID: receiptUUID.String()}

	respJSON, err = json.Marshal(resp)
	if err != nil {
		slog.Error("error encoding response", "err", err)
		writer.WriteHeader(http.StatusInternalServerError)

		return
	}

	respJSON = append(respJSON, '\n')

	writer.WriteHeader(http.StatusOK)

	_, err = writer.Write(respJSON)
	if err != nil {
		slog.Error("error sending response", "err", err)

		return
	}
}

// calculatePoints calculates and returns the total number of points awarded for the receipt
func calculatePoints(receipt Receipt) uint64 {
	var points uint64

	// add retailer name points
	points += retailerPoints(receipt.Retailer)

	// add purchase total points
	points += purchaseTotalPoints(receipt.Total)

	// item count points
	points += itemCountPoints(receipt.Items)

	// item description points for each item
	for _, v := range receipt.Items {
		points += descriptionPoints(v)
	}

	// add date points
	points += datePoints(receipt.PurchaseDate)

	// add points for time of date
	points += timeOfDayPoints(receipt.PurchaseTime)

	return points
}

// retailerPoints calculates and returns the number of points awarded for the retailer name
//
// * One point for every alphanumeric character in the retailer name
func retailerPoints(retailer string) uint64 {
	var points uint64

	// values must be kept sorted
	myRT := &unicode.RangeTable{
		R16: []unicode.Range16{
			{0x0030, 0x0039, 1}, // numbers
			{0x0041, 0x005a, 1}, // upper case letters
			{0x0061, 0x007a, 1}, // lower case letters
		},
		LatinOffset: 0,
	}

	for _, i := range retailer {
		if unicode.In(i, myRT) {
			points++
		}
	}

	return points
}

// purchaseTotalPoints calculates and returns the number of points awarded for the purchase total
//
// * 50 points if the total is a round dollar amount with no cents.
// * 25 points if the total is a multiple of `0.25`.
func purchaseTotalPoints(total string) uint64 {
	var points uint64

	totalFloat, err := strconv.ParseFloat(total, 64)
	if err != nil {
		return points
	}

	if math.Mod(totalFloat, 0.50) == 0 {
		points += 50
	}

	if math.Mod(totalFloat, 0.25) == 0 {
		points += 25
	}

	return points
}

// itemCountPoints calculates and returns the number of points awarded for the number of items
//
// * 5 points for every two items on the receipt.
func itemCountPoints(items []Item) uint64 {
	return uint64((len(items) / 2) * 5)
}

// descriptionPoints calculates and returns the number of points awarded for the item description
//
//   - If the trimmed length of the item description is a multiple of 3, multiply the price by `0.2`
//     and round up to the nearest integer. The result is the number of points earned.
func descriptionPoints(item Item) uint64 {
	var points uint64

	tLen := len(strings.Trim(item.ShortDescription, " "))

	if tLen%3 == 0 {
		p, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			return points
		}

		points = +uint64(math.Ceil(p * 0.2))
	}

	return points
}

// datePoints calculates and returns the number of points awarded for the purchase date
//
// * 6 points if the day in the purchase date is odd.
func datePoints(date openapiTypes.Date) uint64 {
	var points uint64

	purchaseDay := date.Day()

	if purchaseDay%2 == 1 {
		points += 6
	}

	return points
}

// timeOfDayPoints calculates and returns the number of points awarded for the time of day of the purchase
//
// * 10 points if the time of purchase is after 2:00pm and before 4:00pm.
func timeOfDayPoints(purchaseTimeStr string) uint64 {
	var points uint64

	purchaseTime, err := time.Parse("15:04", purchaseTimeStr)
	if err != nil {
		return points
	}

	purchaseTimeHour := purchaseTime.Hour()
	purchaseTimeMin := purchaseTime.Minute()

	if purchaseTimeHour < 16 {
		// 2pm is not after 2pm
		if (purchaseTimeHour == 14 && purchaseTimeMin > 0) || purchaseTimeHour > 14 {
			points += 10
		}
	}

	return points
}
