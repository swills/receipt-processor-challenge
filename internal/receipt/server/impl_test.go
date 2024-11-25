package receiptserver

import (
	"testing"
	"time"

	"github.com/oapi-codegen/runtime/types"
)

func Test_retailerPoints(t *testing.T) {
	type args struct {
		s string
	}

	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "example1/simple",
			args: args{s: "Target"},
			want: 6,
		},
		{
			name: "example2",
			args: args{s: "M&M Corner Market"},
			want: 14,
		},
		{
			name: "morning",
			args: args{s: "Walgreens"},
			want: 9,
		},
	}

	t.Parallel()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := retailerPoints(testCase.args.s)
			if got != testCase.want {
				t.Errorf("retailerPoints() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func Test_purchaseTotalPoints(t *testing.T) {
	type args struct {
		total string
	}

	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "example1",
			args: args{total: "35.35"},
			want: 0,
		},
		{
			name: "example2",
			args: args{total: "9.00"},
			want: 75,
		},
		{
			name: "simple",
			args: args{total: "1.25"},
			want: 25,
		},
		{
			name: "morning",
			args: args{total: "2.65"},
			want: 0,
		},
		{
			name: "badFloat",
			args: args{total: "NotANumber"}, // strconv.ParseFloat recognizes "NaN" and friends
			want: 0,
		},
	}

	t.Parallel()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := purchaseTotalPoints(testCase.args.total)
			if got != testCase.want {
				t.Errorf("purchaseTotalPoints() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func Test_itemCountPoints(t *testing.T) {
	type args struct {
		items []Item
	}

	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "example1",
			args: args{
				items: []Item{
					{
						Price:            "6.49",
						ShortDescription: "Mountain Dew 12PK",
					},
					{
						Price:            "12.25",
						ShortDescription: "Emils Cheese Pizza",
					},
					{
						Price:            "1.26",
						ShortDescription: "Knorr Creamy Chicken",
					},
					{
						Price:            "3.35",
						ShortDescription: "Doritos Nacho Cheese",
					},
					{
						Price:            "12.00",
						ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
					},
				},
			},
			want: 10,
		},
		{
			name: "example2",
			args: args{
				items: []Item{
					{
						Price:            "2.25",
						ShortDescription: "Gatorade",
					},
					{
						Price:            "2.25",
						ShortDescription: "Gatorade",
					},
					{
						Price:            "2.25",
						ShortDescription: "Gatorade",
					},
					{
						Price:            "2.25",
						ShortDescription: "Gatorade",
					},
				},
			},
			want: 10,
		},
		{
			name: "simple",
			args: args{
				items: []Item{
					{
						Price:            "1.25",
						ShortDescription: "Pepsi - 12-oz",
					},
				},
			},
			want: 0,
		},
		{
			name: "morning",
			args: args{
				items: []Item{
					{
						Price:            "1.25",
						ShortDescription: "Pepsi - 12-oz",
					},
					{
						Price:            "1.40",
						ShortDescription: "Dasani",
					},
				},
			},
			want: 5,
		},
	}

	t.Parallel()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := itemCountPoints(testCase.args.items)
			if got != testCase.want {
				t.Errorf("itemCountPoints() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func Test_descriptionPoints(t *testing.T) {
	type args struct {
		item Item
	}

	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "example1.item1",
			args: args{item: Item{
				Price:            "6.49",
				ShortDescription: "Mountain Dew 12PK"},
			},
			want: 0,
		},
		{
			name: "example1.item2",
			args: args{
				item: Item{
					Price:            "12.25",
					ShortDescription: "Emils Cheese Pizza",
				},
			},
			want: 3,
		},
		{
			name: "example1.item3",
			args: args{
				item: Item{
					Price:            "1.26",
					ShortDescription: "Knorr Creamy Chicken",
				},
			},
			want: 0,
		},
		{
			name: "example1.item4",
			args: args{
				item: Item{
					Price:            "3.35",
					ShortDescription: "Doritos Nacho Cheese",
				},
			},
			want: 0,
		},
		{
			name: "example1.item5",
			args: args{
				item: Item{
					Price:            "12.00",
					ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
				},
			},
			want: 3,
		},
		{
			name: "example2.item1",
			args: args{item: Item{
				Price:            "2.25",
				ShortDescription: "Gatorade"},
			},
			want: 0,
		},
		{
			name: "simple.item1",
			args: args{item: Item{
				Price:            "1.25",
				ShortDescription: "Pepsi - 12-oz"},
			},
			want: 0,
		},
		{
			name: "morning.item2",
			args: args{item: Item{
				Price:            "1.40",
				ShortDescription: "Dasani"},
			},
			want: 1,
		},
		{
			name: "badPrice",
			args: args{item: Item{
				Price:            "1.40.1",
				ShortDescription: "Dasani"},
			},
			want: 0,
		},
	}

	t.Parallel()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := descriptionPoints(testCase.args.item)
			if got != testCase.want {
				t.Errorf("descriptionPoints() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func Test_datePoints(t *testing.T) {
	t.Parallel()

	example1Time, err := time.Parse("2006-01-02", "2022-01-01")
	if err != nil {
		t.Error(err)
	}

	example2Time, err := time.Parse("2006-01-02", "2022-03-20")
	if err != nil {
		t.Error(err)
	}

	evenDayTime, err := time.Parse("2006-01-02", "2022-01-02")
	if err != nil {
		t.Error(err)
	}

	type args struct {
		date types.Date
	}

	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "example1",
			args: args{date: types.Date{Time: example1Time}},
			want: 6,
		},
		{
			name: "example2",
			args: args{date: types.Date{Time: example2Time}},
			want: 0,
		},
		{
			name: "evenDay",
			args: args{date: types.Date{Time: evenDayTime}},
			want: 0,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := datePoints(testCase.args.date)
			if got != testCase.want {
				t.Errorf("datePoints() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func Test_timeOfDayPoints(t *testing.T) {
	type args struct {
		time string
	}

	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "example1",
			args: args{time: "13:01"},
			want: 0,
		},
		{
			name: "example2",
			args: args{time: "14:33"},
			want: 10,
		},
		{
			name: "simple",
			args: args{time: "13:13"},
			want: 0,
		},
		{
			name: "morning",
			args: args{time: "08:13"},
			want: 0,
		},
		{
			name: "before2pm",
			args: args{time: "13:59"},
			want: 0,
		},
		{
			name: "twopm",
			args: args{time: "14:00"},
			want: 0,
		},
		{
			name: "after2pm",
			args: args{time: "14:01"},
			want: 10,
		},
		{
			name: "badTime",
			args: args{time: "99:99"},
			want: 0,
		},
	}

	t.Parallel()

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := timeOfDayPoints(testCase.args.time)
			if got != testCase.want {
				t.Errorf("timeOfDayPoints() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func Test_calculatePoints(t *testing.T) {
	t.Parallel()

	var err error

	var docExample1PurchaseDate time.Time

	var docExample2PurchaseDate time.Time

	var simplePurchaseDate time.Time

	var morningPurchaseDate time.Time

	docExample1PurchaseDate, err = time.Parse(time.DateOnly, "2022-01-01")
	if err != nil {
		t.Error(err)
	}

	docExample2PurchaseDate, err = time.Parse(time.DateOnly, "2022-03-20")
	if err != nil {
		t.Error(err)
	}

	simplePurchaseDate, err = time.Parse(time.DateOnly, "2022-01-02")
	if err != nil {
		t.Error(err)
	}

	morningPurchaseDate, err = time.Parse(time.DateOnly, "2022-01-02")
	if err != nil {
		t.Error(err)
	}

	type args struct {
		r Receipt
	}

	tests := []struct {
		name string
		args args
		want uint64
	}{
		{
			name: "docExample1",
			args: args{
				r: Receipt{
					Items: []Item{
						{
							Price:            "6.49",
							ShortDescription: "Mountain Dew 12PK",
						},
						{
							Price:            "12.25",
							ShortDescription: "Emils Cheese Pizza",
						},
						{
							Price:            "1.26",
							ShortDescription: "Knorr Creamy Chicken",
						},
						{
							Price:            "3.35",
							ShortDescription: "Doritos Nacho Cheese",
						},
						{
							Price:            "12.00",
							ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
						},
					},
					PurchaseDate: types.Date{
						Time: docExample1PurchaseDate,
					},
					PurchaseTime: "13:01",
					Retailer:     "Target",
					Total:        "35.35",
				},
			},
			want: 28,
		},
		{
			name: "docExample2",
			args: args{
				r: Receipt{
					Items: []Item{
						{
							Price:            "2.25",
							ShortDescription: "Gatorade",
						},
						{
							Price:            "2.25",
							ShortDescription: "Gatorade",
						},
						{
							Price:            "2.25",
							ShortDescription: "Gatorade",
						},
						{
							Price:            "2.25",
							ShortDescription: "Gatorade",
						},
					},
					PurchaseDate: types.Date{
						Time: docExample2PurchaseDate,
					},
					PurchaseTime: "14:33",
					Retailer:     "M&M Corner Market",
					Total:        "9.00",
				},
			},
			want: 109,
		},
		{
			name: "simple",
			args: args{
				r: Receipt{
					Items: []Item{
						{
							Price:            "1.25",
							ShortDescription: "Pepsi - 12-oz",
						},
					},
					PurchaseDate: types.Date{
						Time: simplePurchaseDate,
					},
					PurchaseTime: "13:13",
					Retailer:     "Target",
					Total:        "1.25",
				},
			},
			want: 31,
		},
		{
			name: "morning",
			args: args{
				r: Receipt{
					Items: []Item{
						{
							Price:            "1.25",
							ShortDescription: "Pepsi - 12-oz",
						},
						{
							Price:            "1.40",
							ShortDescription: "Dasani",
						},
					},
					PurchaseDate: types.Date{
						Time: morningPurchaseDate,
					},
					PurchaseTime: "08:13",
					Retailer:     "Walgreens",
					Total:        "2.65",
				},
			},
			want: 15,
		},
	}

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			got := calculatePoints(testCase.args.r)
			if got != testCase.want {
				t.Errorf("calculatePoints() = %v, want %v", got, testCase.want)
			}
		})
	}
}
