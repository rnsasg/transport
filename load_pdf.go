package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
)

const (
	columnWidth float64 = 6
)

type Entry struct {
	Name       string
	Place      string
	Quantity   int
	Paid       int
	ToPaid     int
	Collection int
	Comments   string
}

type Expense struct {
	SerialNumber int
	ExpenseName  string
	Amount       int
	Comments     string
}

type Income struct {
	SerialNumber int
	IncomeName   string
	Amount       int
	Comments     string
}

var entrySerialNumber int
var incomeSerialNumber int
var expenseSerialNumber int

// NextSerial returns the next sequential entry serial number
func NextEntrySerial() int {
	entrySerialNumber++
	return entrySerialNumber // Return the incremented value
}

// NextIncomeSerial returns the next sequential income serial number
func NextIncomeSerial() int {
	incomeSerialNumber++
	return incomeSerialNumber // Return the incremented value
}

// NextExpenseSerial returns the next sequential expense serial number
func NextExpenseSerial() int {
	expenseSerialNumber++
	return expenseSerialNumber // Return the incremented value
}

func readEntries(filename string) ([]Entry, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var entries []Entry
	for _, line := range lines {
		quantity, _ := strconv.Atoi(line[2])
		paid, _ := strconv.Atoi(line[3])
		toPaid, _ := strconv.Atoi(line[4])
		collection, _ := strconv.Atoi(line[5])

		entry := Entry{
			Name:       line[0],
			Place:      line[1],
			Quantity:   quantity,
			Paid:       paid,
			ToPaid:     toPaid,
			Collection: collection,
			Comments:   line[6],
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func readExpenses(filename string) ([]Expense, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var expenses []Expense
	for _, line := range lines {
		amount, _ := strconv.Atoi(line[1])

		expense := Expense{
			SerialNumber: NextExpenseSerial(),
			ExpenseName:  line[0],
			Amount:       amount,
			Comments:     line[2],
		}
		expenses = append(expenses, expense)
	}

	return expenses, nil
}

func readIncomes(filename string) ([]Income, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var incomes []Income
	for _, line := range lines {
		amount, _ := strconv.Atoi(line[1])

		income := Income{
			SerialNumber: NextIncomeSerial(),
			IncomeName:   line[0],
			Amount:       amount,
			Comments:     line[2],
		}
		incomes = append(incomes, income)
	}

	return incomes, nil
}

func GetCurrentDate() string {
	// Get the current date
	currentDate := time.Now()

	// Format the date as DD-MM-YYYY
	formattedDate := currentDate.Format("02-01-2006")

	return formattedDate
}

func GetCurrentTime() string {
	// Get the current time
	currentTime := time.Now()

	// Format the time as HH:MM
	formattedTime := currentTime.Format("15:04")

	return formattedTime
}

func GetCurrentDay() string {
	// Get the current date
	currentTime := time.Now()

	// Get the day of the week
	day := currentTime.Weekday()

	// Convert the day of the week to a string
	dayString := day.String()

	return dayString
}

func generatePDF(entries []Entry, expenses []Expense, incomes []Income) {

	var (
		totalQty        int
		totalPaid       int
		totalToPaid     int
		totalCollection int
		totalExpense    int
		totatIncome     int
		total           string = "Total"
	)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Set margins (left, top, right, bottom)
	leftMargin := 5.0
	topMargin := 2.0
	rightMargin := 5.0
	//bottomMargin := 10.0
	pdf.SetMargins(leftMargin, topMargin, rightMargin)
	pdf.SetFont("Arial", "", 8)

	// Calculate the width of the page and the table
	// pageWidth, _ := pdf.GetPageSize()
	// tableWidth := 100.0 // Adjust this value based on your table width

	// Calculate the x-coordinate to center the table horizontally
	// x := (pageWidth - leftMargin - rightMargin - tableWidth) / 2

	// Set the position to start drawing the table
	// pdf.SetX(x)

	// Add entries table
	addTable(pdf, "Load", []string{"SN", "Name", "Place", "Quantity", "Paid", "To Paid", "Collection", "Comments"}, func() [][]string {
		var data [][]string

		for _, entry := range entries {
			totalQty = totalQty + entry.Quantity
			totalPaid = totalPaid + entry.Paid
			totalToPaid = totalToPaid + entry.ToPaid
			totalCollection = totalCollection + entry.Collection

			data = append(data, []string{
				strconv.Itoa(NextEntrySerial()),
				entry.Name,
				entry.Place,
				strconv.Itoa(entry.Quantity),
				strconv.Itoa(entry.Paid),
				strconv.Itoa(entry.ToPaid),
				strconv.Itoa(entry.Collection),
				entry.Comments,
			})
		}
		totalEntry := []string{
			total,
			GetCurrentDate(),
			GetCurrentTime(),
			strconv.Itoa(totalQty),
			strconv.Itoa(totalPaid),
			strconv.Itoa(totalToPaid),
			strconv.Itoa(totalCollection),
			"",
		}
		data = append(data, totalEntry)
		return data
	})

	// Add expenses table
	addTable(pdf, "Expenses", []string{"SN", "Expense Name", "Amount", "Notes"}, func() [][]string {
		var data [][]string
		for _, expense := range expenses {
			totalExpense = totalExpense + expense.Amount
			data = append(data, []string{
				strconv.Itoa(expense.SerialNumber),
				expense.ExpenseName,
				strconv.Itoa(expense.Amount),
				expense.Comments,
			})
		}
		totalEntry := []string{
			total,
			GetCurrentDate(),
			strconv.Itoa(totalExpense),
			GetCurrentTime(),
		}
		data = append(data, totalEntry)
		return data
	})

	// Add incomes table
	addTable(pdf, "Incomes", []string{"SN", "Income Name", "Amount", "Notes"}, func() [][]string {
		var data [][]string
		for _, income := range incomes {
			totatIncome = totatIncome + income.Amount
			data = append(data, []string{
				strconv.Itoa(income.SerialNumber),
				income.IncomeName,
				strconv.Itoa(income.Amount),
				income.Comments,
			})
		}

		totalEntry := []string{
			total,
			GetCurrentDate(),
			strconv.Itoa(totatIncome),
			GetCurrentTime(),
		}
		data = append(data, totalEntry)
		return data
	})

	addTable(pdf, "Account", []string{"Date", "Time", "Day", "Paid + Collection", "Total Income", "Total Expense", "Profit", "Not paid"}, func() [][]string {
		var data [][]string
		totalEntry := []string{
			GetCurrentDate(),
			GetCurrentTime(),
			GetCurrentDay(),
			strconv.Itoa(totalPaid + totalCollection),
			strconv.Itoa(totatIncome),
			strconv.Itoa(totalExpense),
			strconv.Itoa((totalPaid + totalCollection + totatIncome) - totalExpense),
			strconv.Itoa(totalToPaid - totalCollection),
		}
		data = append(data, totalEntry)
		return data
	})

	err := pdf.OutputFileAndClose("report.pdf")
	if err != nil {
		fmt.Println("Error generating PDF:", err)
	}
}

func addTable(pdf *gofpdf.Fpdf, title string, header []string, getData func() [][]string) {
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(190, 2, title, "", 1, "L", false, 0, "")
	pdf.Ln(10)

	// Calculate column widths
	colWidths := make([]float64, len(header))
	for i, col := range header {
		switch col {
		case "Name":
			colWidths[i] = pdf.GetStringWidth(col) + 25
		case "Comments":
			colWidths[i] = pdf.GetStringWidth(col) + 40
		case "Notes":
			colWidths[i] = pdf.GetStringWidth(col) + 25
		case "Place":
			colWidths[i] = pdf.GetStringWidth(col) + 10
		case "Income Name":
			colWidths[i] = pdf.GetStringWidth(col) + 10
		case "Expense Name":
			colWidths[i] = pdf.GetStringWidth(col) + 10
		case "Date":
			colWidths[i] = pdf.GetStringWidth(col) + 5
		case "Time":
			colWidths[i] = pdf.GetStringWidth(col) + 5
		case "Day":
			colWidths[i] = pdf.GetStringWidth(col) + 5
		default:
			colWidths[i] = pdf.GetStringWidth(col) // Add padding
		}

	}

	// Add header
	pdf.SetFont("Arial", "", 8)

	for i, col := range header {

		pdf.CellFormat(colWidths[i], columnWidth, col, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	pdf.SetFont("Arial", "", 7)
	// Add data
	data := getData()
	for _, row := range data {
		for i, cell := range row {
			pdf.CellFormat(colWidths[i], columnWidth, cell, "1", 0, "C", false, 0, "")
		}
		pdf.Ln(-1)
	}
	pdf.Ln(10)
}

func main() {
	entries, err := readEntries("entry.csv")
	if err != nil {
		fmt.Println("Error reading entries:", err)
		return
	}

	expenses, err := readExpenses("expense.csv")
	if err != nil {
		fmt.Println("Error reading expenses:", err)
		return
	}

	incomes, err := readIncomes("income.csv")
	if err != nil {
		fmt.Println("Error reading incomes:", err)
		return
	}

	generatePDF(entries, expenses, incomes)
}
