package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/baza-trainee/ataka-help-backend/app/structs"
)

type SliderRepo interface {
	SelectSlider() ([]structs.Slider, error)
	InsertSlider(context.Context, structs.Slider) error
}

type SliderService struct {
	Repo SliderRepo
}

func (s SliderService) ReturnSlider() ([]structs.Slider, error) {
	response, err := s.Repo.SelectSlider()
	if err != nil {
		return nil, fmt.Errorf("error happens while slider returning: %w", err)
	}

	return response, nil
}

func (s SliderService) SaveSlider(ctx context.Context, form *multipart.Form) error {

	file := form.File["thumb"][0]

	slider := structs.Slider{
		Title: form.Value["title"][0],
		Thumb: uniqueFilePath(file.Filename, uploadDirectory),
	}

	fileOpened, err := file.Open()
	if err != nil {
		return fmt.Errorf("error happens while file.Open(): %w", err)
	}

	osFile, err := os.OpenFile(slider.Thumb, os.O_WRONLY|os.O_CREATE, filePermition)
	if err != nil {
		return fmt.Errorf("error happens while os.OpenFile(): %w", err)
	}

	defer osFile.Close()

	written, err := io.Copy(osFile, fileOpened)
	if err != nil {
		return fmt.Errorf(" written bytes: %v, error happens while io.Copy(): %w", written, err)
	}

	if err := s.Repo.InsertSlider(ctx, slider); err != nil {
		if err := os.Remove(slider.Thumb); err != nil {
			return fmt.Errorf("error happens while remove file: %w", err)
		}

		return fmt.Errorf("error happens while inserting: %w", err)
	}

	return nil
}
