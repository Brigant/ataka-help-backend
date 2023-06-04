package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

const reportFileName = "report.pdf"

type ReportService struct{}

func (s ReportService) ReturnReport() string {
	reportPath := uploadDirectory + reportFileName

	return reportPath
}

func (s ReportService) ChangeReport(form *multipart.Form) error {
	fileHeader := form.File["report"][0]

	filePath := uploadDirectory + reportFileName

	file, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("error happens while fileHeader.Open(): %w", err)
	}

	osFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, filePermition)
	if err != nil {
		return fmt.Errorf("error happens while os.OpenFile(): %w", err)
	}

	defer osFile.Close()

	written, err := io.Copy(osFile, file)
	if err != nil {
		return fmt.Errorf(" written bytes: %v, error happens while io.Copy(): %w", written, err)
	}

	return nil
}
