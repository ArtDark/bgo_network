package card

import (
	"log"
	"reflect"
	"runtime"
	"testing"
)

// Unit-tests --------------------------------------------------------

func TestSumCategoryTransactions(t *testing.T) {
	type args struct {
		transactions []Transaction
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int64
		wantErr error
	}{
		{
			name: "Valid transactions",
			args: args{
				transactions: []Transaction{
					{Id: "0001", Bill: 100_00, Time: 1606192422, MCC: "5411", Status: "Done"},
					{Id: "0002", Bill: 200_00, Time: 1606192432, MCC: "5812", Status: "Done"},
					{Id: "0003", Bill: 400_00, Time: 1606192442, MCC: "5411", Status: "Done"},
					{Id: "0004", Bill: 300_00, Time: 1606192462, MCC: "5812", Status: "Done"},
				},
			},
			want: map[string]int64{
				"5411": 500_00,
				"5812": 500_00,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SumCategoryTransactions(tt.args.transactions)
			if err != nil {
				t.Errorf("SumCategoryTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SumCategoryTransactions() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSumCategoryTransactionsMutex(t *testing.T) {
	type args struct {
		transactions []Transaction
		goroutines   int
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int64
		wantErr error
	}{
		{
			name: "Valid transactions",
			args: args{
				transactions: []Transaction{
					{Id: "0001", Bill: 100_00, Time: 1606192422, MCC: "5411", Status: "Done"},
					{Id: "0002", Bill: 200_00, Time: 1606192432, MCC: "5812", Status: "Done"},
					{Id: "0003", Bill: 400_00, Time: 1606192442, MCC: "5411", Status: "Done"},
					{Id: "0004", Bill: 300_00, Time: 1606192462, MCC: "5812", Status: "Done"},
				},
				goroutines: 2,
			},
			want: map[string]int64{
				"5411": 500_00,
				"5812": 500_00,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SumCategoryTransactionsMutex(tt.args.transactions, tt.args.goroutines)
			if err != tt.wantErr {
				t.Errorf("SumCategoryTransactionsMutex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SumCategoryTransactionsMutex() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSumCategoryTransactionsChan(t *testing.T) {
	type args struct {
		transactions []Transaction
		goroutines   int
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int64
		wantErr error
	}{
		{
			name: "Valid transactions",
			args: args{
				transactions: []Transaction{
					{Id: "0001", Bill: 100_00, Time: 1606192422, MCC: "5411", Status: "Done"},
					{Id: "0002", Bill: 200_00, Time: 1606192432, MCC: "5812", Status: "Done"},
					{Id: "0003", Bill: 400_00, Time: 1606192442, MCC: "5411", Status: "Done"},
					{Id: "0004", Bill: 300_00, Time: 1606192462, MCC: "5812", Status: "Done"},
				},
				goroutines: 2,
			},
			want: map[string]int64{
				"5411": 500_00,
				"5812": 500_00,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SumCategoryTransactionsChan(tt.args.transactions, tt.args.goroutines)
			if err != tt.wantErr {
				t.Errorf("SumCategoryTransactionsChan() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SumCategoryTransactionsChan() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSumCategoryTransactionsMutexWithoutFunc(t *testing.T) {
	type args struct {
		transactions []Transaction
		goroutines   int
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]int64
		wantErr error
	}{
		{
			name: "Valid transactions",
			args: args{
				transactions: []Transaction{
					{Id: "0001", Bill: 100_00, Time: 1606192422, MCC: "5411", Status: "Done"},
					{Id: "0002", Bill: 200_00, Time: 1606192432, MCC: "5812", Status: "Done"},
					{Id: "0003", Bill: 400_00, Time: 1606192442, MCC: "5411", Status: "Done"},
					{Id: "0004", Bill: 300_00, Time: 1606192462, MCC: "5812", Status: "Done"},
				},
				goroutines: 2,
			},
			want: map[string]int64{
				"5411": 500_00,
				"5812": 500_00,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SumCategoryTransactionsMutexWithoutFunc(tt.args.transactions, tt.args.goroutines)
			if err != tt.wantErr {
				t.Errorf("SumCategoryTransactionsMutexWithoutFunc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SumCategoryTransactionsMutexWithoutFunc() got = %v, want %v", got, tt.want)
			}
		})
	}
}

// Benchmark tests --------------------------------------------------------

func BenchmarkSumCategoryTransactions(b *testing.B) {

	user := Card{
		Id: 1,
		Owner: Owner{
			FirstName: "User",
			LastName:  "User",
		},
		Issuer:       "Visa",
		Balance:      5000_00,
		Currency:     "RUR",
		Number:       "4619071400941155",
		Icon:         "https://cdn.visa.com/cdn/assets/images/logos/visa/logo.png",
		Transactions: []Transaction{},
	}

	err := user.MakeTransactions(1_000_000)
	if err != nil {
		log.Panicln(err)
	}

	want := map[string]int64{
		"5411": 5_099_995_000_00,
		"5812": 5_101_995_000_00,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка

	for i := 0; i < b.N; i++ {
		result, err := SumCategoryTransactions(user.Transactions)
		if err != nil {
			log.Println(err)
			break
		}
		b.StopTimer() // останавливаем таймер, чтобы время сравнения не учитывалось
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer() // продолжаем работу таймера
	}
}

func BenchmarkSumCategoryTransactionsMutex(b *testing.B) {

	user := Card{
		Id: 1,
		Owner: Owner{
			FirstName: "User",
			LastName:  "User",
		},
		Issuer:       "Visa",
		Balance:      5000_00,
		Currency:     "RUR",
		Number:       "4619071400941155",
		Icon:         "https://cdn.visa.com/cdn/assets/images/logos/visa/logo.png",
		Transactions: []Transaction{},
	}

	err := user.MakeTransactions(1_000_000)
	if err != nil {
		log.Panicln(err)
	}

	want := map[string]int64{
		"5411": 5_099_995_000_00,
		"5812": 5_101_995_000_00,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка

	for i := 0; i < b.N; i++ {
		result, err := SumCategoryTransactionsMutex(user.Transactions, runtime.NumCPU())
		if err != nil {
			log.Println(err)
			break
		}
		b.StopTimer() // останавливаем таймер, чтобы время сравнения не учитывалось
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer() // продолжаем работу таймера
	}
}

func BenchmarkSumCategoryTransactionsChan(b *testing.B) {

	user := Card{
		Id: 1,
		Owner: Owner{
			FirstName: "User",
			LastName:  "User",
		},
		Issuer:       "Visa",
		Balance:      5000_00,
		Currency:     "RUR",
		Number:       "4619071400941155",
		Icon:         "https://cdn.visa.com/cdn/assets/images/logos/visa/logo.png",
		Transactions: []Transaction{},
	}

	err := user.MakeTransactions(1_000_000)
	if err != nil {
		log.Panicln(err)
	}

	want := map[string]int64{
		"5411": 5_099_995_000_00,
		"5812": 5_101_995_000_00,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка

	for i := 0; i < b.N; i++ {
		result, err := SumCategoryTransactionsChan(user.Transactions, runtime.NumCPU())
		if err != nil {
			log.Println(err)
			break
		}
		b.StopTimer() // останавливаем таймер, чтобы время сравнения не учитывалось
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer() // продолжаем работу таймера
	}
}

func BenchmarkSumCategoryTransactionsMutexWithoutFunc(b *testing.B) {

	user := Card{
		Id: 1,
		Owner: Owner{
			FirstName: "User",
			LastName:  "User",
		},
		Issuer:       "Visa",
		Balance:      5000_00,
		Currency:     "RUR",
		Number:       "4619071400941155",
		Icon:         "https://cdn.visa.com/cdn/assets/images/logos/visa/logo.png",
		Transactions: []Transaction{},
	}

	err := user.MakeTransactions(1_000_000)
	if err != nil {
		log.Panicln(err)
	}

	want := map[string]int64{
		"5411": 5_099_995_000_00,
		"5812": 5_101_995_000_00,
	}
	b.ResetTimer() // сбрасываем таймер, т.к. сама генерация транзакций достаточно ресурсоёмка

	for i := 0; i < b.N; i++ {
		result, err := SumCategoryTransactionsMutexWithoutFunc(user.Transactions, runtime.NumCPU())
		if err != nil {
			log.Println(err)
			break
		}
		b.StopTimer() // останавливаем таймер, чтобы время сравнения не учитывалось
		if !reflect.DeepEqual(result, want) {
			b.Fatalf("invalid result, got %v, want %v", result, want)
		}
		b.StartTimer() // продолжаем работу таймера
	}
}
