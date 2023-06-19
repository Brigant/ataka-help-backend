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
	SelectSlider(context.Context) ([]structs.Slide, error)
	InsertSlider(context.Context, structs.Slide, chan struct{}) error
	DelSlideByID(context.Context, string) (string, error)
}

type SliderService struct {
	Repo SliderRepo
}

func (s SliderService) ReturnSlider(ctx context.Context) ([]structs.Slide, error) {
	response, err := s.Repo.SelectSlider(ctx)
	if err != nil {
		return []structs.Slide{}, fmt.Errorf("error happens while slider returning: %w", err)
	}

	return response, nil
}

func (s SliderService) SaveSlider(ctx context.Context, form *multipart.Form, chWell chan struct{}) error {
	file := form.File["thumb"][0]

	slider := structs.Slide{
		Title: form.Value["title"][0],
		Thumb: uniqueFilePath(file.Filename, uploadDirectory),
		Alt:   form.Value["alt"][0],
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

	if err := s.Repo.InsertSlider(ctx, slider, chWell); err != nil {
		if err := os.Remove(slider.Thumb); err != nil {
			return fmt.Errorf("error happens while remove file: %w", err)
		}

		return fmt.Errorf("error happens while inserting: %w", err)
	}

	return nil
}

func (s SliderService) DeleteSlideByID(ctx context.Context, ID string) error {
	objectPath, err := s.Repo.DelSlideByID(ctx, ID)
	if err != nil {
		return fmt.Errorf("error while delete slide: %w", err)
	}

	if err := os.Remove(objectPath); err != nil {
		return fmt.Errorf("error happens while remove file: %w", err)
	}

	return nil
}
