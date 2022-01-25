package document_builder

import (
	"delivery/internal/services/datetime"
	eb "delivery/internal/services/excel_builder"
	"fmt"
	"github.com/LeKovr/num2word"
	"strconv"
)

const (
	deliveryPath = "c.Чалтырь - х.Ленина"
	fio          = "Срабионян Г.С"
	price        = 11000.00
)

type Builder struct {
	excelBuilder *eb.Builder
	dateInfo     *datetime.DateInfo
}

func New() *Builder {
	return &Builder{
		dateInfo:     datetime.New(),
		excelBuilder: eb.New(),
	}
}

func (b *Builder) Build() error {
	err := b.buildCheckPart()
	if err != nil {
		return err
	}
	err = b.buildDecipherPart()
	if err != nil {
		return err
	}
	return nil

}

func (b *Builder) Init() error {
	err := b.dateInfo.Init()
	if err != nil {
		return err
	}
	fmt.Println(b.excelBuilder.CheckFilePath)
	fmt.Println(b.excelBuilder.DecipherFilePath)
	return nil
}

func (b *Builder) Validate() bool {
	if !b.dateInfo.Validate() {
		return false
	}
	if !b.excelBuilder.Validate() {
		return false
	}
	return true
}

func (b *Builder) SetRawDates(dates string) {
	b.dateInfo.RawDates = dates
}

func (b *Builder) SetYear(year string) {
	b.dateInfo.Year = year
}

func (b *Builder) SetCheckFilePath(checkFilePath string) {
	b.excelBuilder.CheckFilePath = checkFilePath
}

func (b *Builder) SetDecipherFilePath(decipherFilePath string) {
	b.excelBuilder.DecipherFilePath = decipherFilePath
}

func (b *Builder) HasCheckFilePath() bool {
	return b.excelBuilder.CheckFilePath != ""
}

func (b *Builder) HasDecipherFilePath() bool {
	return b.excelBuilder.DecipherFilePath != ""
}

func (b *Builder) buildCheckPart() error {
	checkBuilder, f := b.excelBuilder.GetCheckBuilder()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	currentCheckNumber, err := strconv.Atoi(checkBuilder.GetCellValue(eb.CheckNumberLeftCell))
	if err != nil {
		return err
	}
	startDay, endDay, err := b.dateInfo.GetDateRange()
	if err != nil {
		return err
	}
	month, err := b.dateInfo.GetMonth()
	if err != nil {
		return err
	}
	totalPrice := b.dateInfo.CountDates() * price

	checkBuilder.ChangeCell(eb.CheckNumberLeftCell, currentCheckNumber+1)
	checkBuilder.ChangeCell(eb.CheckNumberRightACell, currentCheckNumber+1)
	checkBuilder.ChangeCell(eb.CheckDayLeftCell, endDay)
	checkBuilder.ChangeCell(eb.CheckDayRightCell, endDay)
	checkBuilder.ChangeCell(eb.CheckMonthLeftCell, month)
	checkBuilder.ChangeCell(eb.CheckMonthRightCell, month)
	checkBuilder.ChangeCell(eb.CheckYearLeftCell, b.dateInfo.Year)
	checkBuilder.ChangeCell(eb.CheckYearRightCell, b.dateInfo.Year)
	checkBuilder.ChangeCell(eb.CheckDeliveryDateStartLeftCell, startDay)
	checkBuilder.ChangeCell(eb.CheckDeliveryDateStartRightCell, startDay)
	checkBuilder.ChangeCell(eb.CheckDeliveryDateEndLeftCell, endDay)
	checkBuilder.ChangeCell(eb.CheckDeliveryDateEndRightCell, endDay)
	checkBuilder.ChangeCell(eb.CheckDeliveryMonthLeftCell, month)
	checkBuilder.ChangeCell(eb.CheckDeliveryMonthRightCell, month)
	checkBuilder.ChangeCell(eb.CheckDeliveryYearLeftCell, b.dateInfo.Year)
	checkBuilder.ChangeCell(eb.CheckDeliveryYearRightCell, b.dateInfo.Year)
	checkBuilder.ChangeCell(eb.CheckDeliveryFlightCountLeftCell, b.dateInfo.CountDates())
	checkBuilder.ChangeCell(eb.CheckDeliveryFlightCountRightCell, b.dateInfo.CountDates())
	checkBuilder.ChangeCell(eb.CheckPriceAsStringLeftStringCell, num2word.RuMoney(float64(b.dateInfo.CountDates()*price), true))
	checkBuilder.ChangeCell(eb.CheckPriceAsStringRightStringCell, num2word.RuMoney(float64(b.dateInfo.CountDates()*price), true))
	checkBuilder.ChangeCell(eb.CheckTotalPriceLeftCell, totalPrice)
	checkBuilder.ChangeCell(eb.CheckAllTotalPriceLeftCell, totalPrice)
	checkBuilder.ChangeCell(eb.CheckTotalPriceAfterTaxLeftCell, totalPrice)
	checkBuilder.ChangeCell(eb.CheckTotalPriceAfterTaxRightCell, totalPrice)
	checkBuilder.ChangeCell(eb.CheckTotalPriceRightCell, totalPrice)
	checkBuilder.ChangeCell(eb.CheckAllTotalPriceRightCell, totalPrice)
	checkBuilder.ChangeCell(eb.CheckRublePriceLeftCell, totalPrice)
	checkBuilder.ChangeCell(eb.CheckRublePriceRightCell, totalPrice)
	checkBuilder.ChangeCell(eb.CheckDeliveryPriceLeftCell, price)
	checkBuilder.ChangeCell(eb.CheckDeliveryPriceRightCell, price)

	err = checkBuilder.Save()
	if err != nil {
		return err
	}
	return nil
}

func (b *Builder) buildDecipherPart() error {
	decipherBuilder, f := b.excelBuilder.GetDecipherBuilder()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	for i := 0; i < 26; i++ {
		decipherBuilder.ClearCell(eb.DecipherNumber + strconv.Itoa(eb.GlobalStartCell+i))
		decipherBuilder.ClearCell(eb.DecipherPath + strconv.Itoa(eb.GlobalStartCell+i))
		decipherBuilder.ClearCell(eb.DecipherPrice + strconv.Itoa(eb.GlobalStartCell+i))
		decipherBuilder.ClearCell(eb.DecipherFIO + strconv.Itoa(eb.GlobalStartCell+i))
		decipherBuilder.ClearCell(eb.DecipherDate + strconv.Itoa(eb.GlobalStartCell+i))
		decipherBuilder.ClearCell(eb.DecipherDeliveryCount + strconv.Itoa(eb.GlobalStartCell+i))
		decipherBuilder.ClearCell(eb.DecipherTotalPrice + strconv.Itoa(eb.GlobalStartCell+i))
	}
	dateRange, err := b.dateInfo.GetDateRangeString()
	if err != nil {
		return err
	}
	number := 0
	for _, v := range b.dateInfo.GetDates() {
		decipherBuilder.ChangeCell(eb.DecipherNumber+strconv.Itoa(eb.GlobalStartCell+number), number+1)
		decipherBuilder.ChangeCell(eb.DecipherPath+strconv.Itoa(eb.GlobalStartCell+number), deliveryPath)
		decipherBuilder.ChangeCell(eb.DecipherPrice+strconv.Itoa(eb.GlobalStartCell+number), price)
		decipherBuilder.ChangeCell(eb.DecipherFIO+strconv.Itoa(eb.GlobalStartCell+number), fio)
		decipherBuilder.ChangeCell(eb.DecipherDate+strconv.Itoa(eb.GlobalStartCell+number), v)
		decipherBuilder.ChangeCell(eb.DecipherDeliveryCount+strconv.Itoa(eb.GlobalStartCell+number), 1)
		decipherBuilder.ChangeCell(eb.DecipherTotalPrice+strconv.Itoa(eb.GlobalStartCell+number), price)
		decipherBuilder.ChangeCell(eb.DecipherDateRange, dateRange)
		decipherBuilder.ChangeCell(eb.DecipherYear, b.dateInfo.GetYear())
		decipherBuilder.ChangeCell(eb.DecipherAllTotalPrice, price*b.dateInfo.CountDates())
		decipherBuilder.ChangeCell(eb.DecipherTotalDeliveryCount, b.dateInfo.CountDates())

		number++
	}
	err = decipherBuilder.Save()
	if err != nil {
		return err
	}
	return nil
}
