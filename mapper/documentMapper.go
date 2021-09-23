package mapper

import (
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/namrahov/ms-ecourt-go/model"
)

func FillExcelStaticColumnsForReport(file *excelize.File) {
	file.SetColWidth(model.Sheet, "A", "D", 30)
	file.SetCellValue(model.Sheet, "A1", "ID")
	file.SetCellValue(model.Sheet, "B1", "COURT_NAME")
	file.SetCellValue(model.Sheet, "C1", "JUDGE_NAME")
}
