//Package card
package card

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrCardNotFound   = errors.New("card not found")
	ErrNoTransactions = errors.New("no user transactions")
)

// Описание банковской карты"
type Card struct {
	Id           cardId
	Owner               // Владелец карты
	Issuer       string // Платежная система
	Balance      int    // Баланс карты
	Currency     string // Валюта
	Number       string // Номер карты в платежной системе
	Icon         string // Иконка платежной системы
	Transactions Transactions
}

// Идентификатор банковской карты
type cardId int64

// Инициалы владельца банковской карты
type Owner struct {
	FirstName string // Имя владельца карты
	LastName  string // Фамилия владельца карты
}

type Transaction struct {
	XMLName string `xml:"transaction"`
	Id      string `json:"id" xml:"id"`
	Bill    int64  `json:"bill" xml:"bill"`
	Time    int64  `json:"time" xml:"time"`
	MCC     string `json:"mcc" xml:"mcc"`
	Status  string `json:"status" xml:"status"`
}

type Transactions struct {
	XMLName      string `xml:"transactions"`
	Transactions []Transaction
}

// Метод добавления транзакции
func (c *Card) AddTransaction(transaction Transaction) {
	c.Transactions.Transactions = append(c.Transactions.Transactions, transaction)

}

// Метод генерации 2х транзакций с разными MCC
func (c *Card) MakeTransactions(count int) error {

	if c == nil {
		return ErrCardNotFound
	}

	if count <= 0 {
		log.Println("count must be > 0")
		return nil
	}

	for i := 0; i < count; i++ {
		c.AddTransaction(Transaction{
			Id: strconv.Itoa((i + 1) + i),

			Bill: int64(100_00),

			Time:   time.Date(2020, 9, 10, 12+i, 23+i, 21+i, 0, time.UTC).Unix(),
			MCC:    "5411",
			Status: "Done",
		})
		c.AddTransaction(Transaction{
			Id: strconv.Itoa((i + 2) + i),

			Bill: int64(102_00),

			Time:   time.Date(2020, 9, 10, 14+i, 15+i, 21+i, 0, time.UTC).Unix(),
			MCC:    "5812",
			Status: "Done",
		})

	}

	return nil

}

// Функция расчета суммы по категории
func SumByMCC(transactions []Transaction, mcc []string) int64 {
	var mmcSum int64

	for _, code := range mcc {
		for _, t := range transactions {
			if code == t.MCC {
				mmcSum += t.Bill
			}
		}
	}

	return mmcSum

}

// Функция преобразования кода в название категории
func TranslateMCC(code string) string {
	// представим, что mcc читается из файла (научимся позже)
	mcc := map[string]string{
		"5411": "Супермаркеты",
		"5812": "Рестораны",
	}

	const errCategoryUndef = "Категория не указана"

	if value, ok := mcc[code]; ok {
		return value
	}

	return errCategoryUndef

}

// Функция сложения сумм транзакций по категориям
func SumCategoryTransactions(transactions []Transaction) (map[string]int64, error) {

	if transactions == nil {
		return nil, ErrNoTransactions
	}

	m := make(map[string]int64)

	for i := range transactions {
		m[transactions[i].MCC] += transactions[i].Bill
	}

	return m, nil

}

// Функция сложения сумм транзакций по категориям с использованием goroutines и mutex
func SumCategoryTransactionsMutex(transactions []Transaction, goroutines int) (map[string]int64, error) {
	wg := sync.WaitGroup{}
	wg.Add(goroutines)

	mu := sync.Mutex{}

	if transactions == nil {
		return nil, ErrNoTransactions
	}

	m := make(map[string]int64)

	partSize := len(transactions) / goroutines

	for i := 0; i < goroutines; i++ {
		part := transactions[i*partSize : (i+1)*partSize]
		go func() {
			mapSum, err := SumCategoryTransactions(part)
			if err != nil {
				fmt.Println(err)
			}
			mu.Lock()
			for key, i := range mapSum {
				m[key] += i

			}
			mu.Unlock()
			wg.Done()
		}()

	}
	wg.Wait()

	return m, nil

}

// Функция сложения сум транзакций по категориям с использованием goroutines и каналов
func SumCategoryTransactionsChan(transactions []Transaction, goroutines int) (map[string]int64, error) {

	if transactions == nil {
		return nil, ErrNoTransactions
	}

	result := make(map[string]int64)
	ch := make(chan map[string]int64)
	partSize := len(transactions) / goroutines

	for i := 0; i < goroutines; i++ {
		part := transactions[i*partSize : (i+1)*partSize]
		go func(ch chan<- map[string]int64) {
			s, err := SumCategoryTransactions(part)
			if err != nil {
				log.Printf("failed to sum: %s\n", err)
			}
			ch <- s
		}(ch)
	}

	fin := 0

	for sum := range ch {
		for k, v := range sum {
			result[k] += v

		}
		fin++
		if fin == goroutines {
			close(ch)
			break
		}
	}

	return result, nil

}

// Функция с mutex'ом, который защищает любые операции с map, соответственно, её задача: разделить слайс транзакций на несколько кусков и в отдельных горутинах посчитать, но теперь горутины напрямую пишут в общий map с результатами. Важно: эта функция внутри себя не должна вызывать функцию из п.1
func SumCategoryTransactionsMutexWithoutFunc(transactions []Transaction, goroutines int) (map[string]int64, error) {
	wg := sync.WaitGroup{}
	wg.Add(goroutines)

	mu := sync.Mutex{}

	if transactions == nil {
		return nil, ErrNoTransactions
	}

	mapSum := make(map[string]int64)

	partSize := len(transactions) / goroutines

	for i := 0; i < goroutines; i++ {
		part := transactions[i*partSize : (i+1)*partSize]
		go func() {

			for i := range part {
				mu.Lock()
				mapSum[part[i].MCC] += part[i].Bill
				mu.Unlock()
			}
			wg.Done()

		}()

	}
	wg.Wait()

	return mapSum, nil

}

// Сервис банка
type Service struct {
	BankName string
	Cards    []*Card
}

// Конструктор сервиса
func New(bankName string) *Service {
	return &Service{BankName: bankName}
}

// Метод создания экземпляра банковской карты
func (s *Service) CardIssue(
	id cardId,
	fistName,
	lastName,
	issuer string,
	balance int,
	currency string,
	number string,
) *Card {
	var card = &Card{
		Id: id,
		Owner: Owner{
			FirstName: fistName,
			LastName:  lastName,
		},
		Issuer:   issuer,
		Balance:  balance,
		Currency: currency,
		Number:   number,
		Icon:     "https://.../logo.png",
	}
	s.Cards = append(s.Cards, card)
	return card
}

const prefix = "5106 21" //Первые 6 цифр нашего банка

// Метод поиска банковской карты по номеру платежной системы
func (s *Service) Card() (*Card, error) {

	for _, c := range s.Cards {
		if strings.HasPrefix(c.Number, prefix) == true {
			return c, nil
		}
	}
	return nil, ErrCardNotFound
}

// Функция экспорта пользовательских транзакций в .csv
func ExporterToCsv(user *Card) error {

	file, err := os.Create("export.csv")

	if err != nil {
		log.Println(err)
		return err
	}

	defer func(c io.Closer) {
		if cerr := c.Close(); cerr != nil {
			log.Println(cerr)
		}
	}(file)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"ID", "Bill", "Time", "MCC", "Status"})
	if err != nil {
		log.Println(err)
		return err
	}

	for _, value := range user.Transactions.Transactions {
		err = writer.Write(transactionToSlice(value))
		if err != nil {
			log.Println(err)
			return err
		}

	}

	return nil
}

// Функция импорта пользовательских транзакций из .csv
func ImporterFromCsv(us *Card, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		log.Println("Cannot open file", err)
	}
	defer func(c io.Closer) {
		if cerr := c.Close(); cerr != nil {
			log.Println("Cannot close file", cerr)
		}
	}(file)

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Println("Cannot read data", err)
	}

	err = us.MapRowToTransaction(records)
	if err != nil {
		return err
	}

	return nil
}

// Функция экспорта пользовательских транзакций в .json
func ExporterToJson(user *Card, fileName string) error {
	file, err := json.MarshalIndent(user.Transactions, "", " ")
	if err != nil {
		log.Println(err)
	}

	err = ioutil.WriteFile(fileName, file, 0644)
	if err != nil {
		log.Println(err)
	}
	return nil
}

// Функция импорта пользовательских транзакций из .json
func ImporterFromJson(user *Card, fileName string) error {

	file, err := os.Open(fileName)
	if err != nil {
		log.Println("Cannot open file", err)
	}
	defer func(c io.Closer) {
		if cerr := c.Close(); cerr != nil {
			log.Println("Cannot close file", cerr)
		}
	}(file)

	reader, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}
	var decoded Transactions
	err = json.Unmarshal(reader, &decoded)
	if err != nil {
		return err
	}

	for _, value := range decoded.Transactions {
		user.Transactions.Transactions = append(user.Transactions.Transactions, value)
	}

	return nil
}

// Функция экспорта пользовательских транзакций в .xml
func ExporterToXml(user *Card, fileName string) error {
	file, err := xml.MarshalIndent(user.Transactions, "", " ")
	if err != nil {
		return err
	}

	file = append([]byte(xml.Header), file...)

	err = ioutil.WriteFile("export.xml", file, 0644)
	if err != nil {
		return err
	}

	return nil
}

// Функция импорта пользовательских транзакций из .xml
func ImporterFromXml(user *Card, fileName string) error {

	file, err := os.Open(fileName)
	if err != nil {
		log.Println("Cannot open file", err)
	}
	defer func(c io.Closer) {
		if cerr := c.Close(); cerr != nil {
			log.Println("Cannot close file", cerr)
		}
	}(file)

	reader, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	var decoded Transactions
	err = xml.Unmarshal(reader, &decoded)
	if err != nil {
		return err
	}

	for _, value := range decoded.Transactions {
		user.Transactions.Transactions = append(user.Transactions.Transactions, value)
	}

	//TODO: Написать логику импорта xml

	return nil
}

func (c *Card) MapRowToTransaction(transactions [][]string) error {

	for _, i := range transactions {
		if i[0] == "ID" {
			continue
		}

		b, err := strconv.Atoi(i[1])
		if err != nil {
			return err
		}
		t, err := strconv.Atoi(i[2])
		if err != nil {
			return err
		}

		transaction := Transaction{
			Id:     i[0],
			Bill:   int64(b),
			Time:   int64(t),
			MCC:    i[3],
			Status: i[4],
		}

		c.AddTransaction(transaction)
	}
	return nil
}

// Функция преобразования пользовательских транзакций в slice
func transactionToSlice(transaction Transaction) []string {

	var data []string

	data = append(data, transaction.Id)
	data = append(data, strconv.Itoa(int(transaction.Bill)))
	data = append(data, strconv.Itoa(int(transaction.Time)))
	data = append(data, transaction.MCC)
	data = append(data, transaction.Status)

	return data

}
