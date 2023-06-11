package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/baza-trainee/ataka-help-backend/app/structs"
)

const reportFileName = "report.pdf"

type ReportService struct{}

func (s ReportService) ReturnReport() (string, error) {
	reportPath := uploadDirectory + reportFileName

	if _, err := os.Stat(reportPath); err != nil {
		if os.IsNotExist(err) {
			return "", structs.ErrNotFound
		}

		return "", fmt.Errorf("error in os.Stat(): %w", err)
	}

	return reportPath, nil
}

func (s ReportService) ChangeReport(form *multipart.Form) error {
	fileHeader := form.File["thumb"][0]

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

func (s ReportService) DeleteReport() error {
	reportPath := uploadDirectory + reportFileName

	if _, err := os.Stat(reportPath); err != nil {
		if os.IsNotExist(err) {
			return structs.ErrNotFound
		}

		return fmt.Errorf("error in os.Stat(): %w", err)
	}

	if err := os.Remove(reportPath); err != nil {
		return fmt.Errorf("error while deleting report: %w", err)
	}

	return nil
}
