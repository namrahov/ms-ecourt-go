package service

import (
	"bytes"
	"context"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
	"github.com/namrahov/ms-ecourt-go/mapper"
	"github.com/namrahov/ms-ecourt-go/model"
	"github.com/namrahov/ms-ecourt-go/repo"
	log "github.com/sirupsen/logrus"
	"html/template"
	"mime/multipart"
	"strconv"
)

type IDocumentService interface {
	GenerateAct(ctx context.Context, dto model.TodoPageData) ([]byte, error)
	GenerateReportOfLightApplication(ctx context.Context) (*excelize.File, error)
	UploadExcel(ctx context.Context, xlsx multipart.File) error
}

type DocumentService struct {
	ApplicationRepo repo.IApplicationRepo
}

func (s *DocumentService) GenerateAct(ctx context.Context, dto model.TodoPageData) ([]byte, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GenerateAct.start")

	tmpl := template.Must(template.ParseFiles("layout.html"))

	var data = make(map[string]interface{})
	data["PageTitle"] = dto.PageTitle
	data["Todos"] = dto.Todos

	file, err := GeneratePdf(ctx, data, tmpl)
	if err != nil {
		logger.Error("ActionLog.GenerateAct.error cant generate pdf", err)
		return nil, err
	}

	logger.Info("ActionLog.GenerateAct.success")
	return file, nil
}

func GeneratePdf(ctx context.Context, data map[string]interface{}, tmpl *template.Template) ([]byte, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GeneratePdf.start")
	var tmplBytes bytes.Buffer

	tmpl.Execute(&tmplBytes, data)

	pdfg, err := wkhtml.NewPDFGenerator()
	if err != nil {
		logger.Error("ActionLog.ConvertHTML2PDF.error ", err)
		return nil, err
	}

	pdfg.AddPage(wkhtml.NewPageReader(&tmplBytes))

	err = pdfg.Create()
	if err != nil {
		logger.Error("ActionLog.ConvertHTML2PDF.error ", err)
		return nil, err
	}

	logger.Info("ActionLog.GeneratePdf.end")
	return pdfg.Bytes(), nil
}

func (s *DocumentService) GenerateReportOfLightApplication(ctx context.Context) (*excelize.File, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GenerateReportOfLightApplication.start")

	applications, err := s.ApplicationRepo.GetApplications()
	if err != nil {
		return nil, err
	}

	var lightApplications = mapper.ApplicationsToLightApplications(applications)

	file := excelize.NewFile()
	mapper.FillExcelStaticColumnsForReport(file)

	for i := 0; i < len(lightApplications); i++ {
		file.SetCellValue(model.Sheet, "A"+strconv.Itoa(2+i), lightApplications[i].Id)
		file.SetCellValue(model.Sheet, "B"+strconv.Itoa(2+i), lightApplications[i].CourtName)
		file.SetCellValue(model.Sheet, "C"+strconv.Itoa(2+i), lightApplications[i].JudgeName)
	}

	logger.Info("ActionLog.GenerateReportOfLightApplication.END")
	return file, nil
}

func (s *DocumentService) UploadExcel(ctx context.Context, xlsx multipart.File) error {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.UploadExcel.start")

	file, err := excelize.OpenReader(xlsx)
	if err != nil {
		log.Error("ActionLog.UploadExcel.response can't open file ", err.Error())
		return err
	}

	rows := file.GetRows(model.FirstSheet)
	var lights []model.LightApplicationDto
	for k := 1; k < 3; k++ {
		var light model.LightApplicationDto

		light.Id, err = strconv.ParseInt(rows[k][0], 10, 64)
		light.CourtName = rows[k][1]
		light.JudgeName = rows[k][2]

		lights = append(lights, light)
	}

	fmt.Println("dto=", lights)

	logger.Info("ActionLog.UploadExcel.end")
	return nil
}
