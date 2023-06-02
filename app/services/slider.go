package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/baza-trainee/ataka-help-backend/app/structs"
	"github.com/gofiber/fiber/v2"
)

const (
	uploadDirectory = "static/uploads/"
	filePermition   = 0o666
)

type SliderRepo interface {
	SelectSlider() (string, error)
	InsertSlider(structs.Slider, context.Context) error
}

type SliderService struct {
	Repo SliderRepo
}

func (s SliderService) ReturnSlider() (string, error) {
	result, err := s.Repo.SelectSlider()
	if err != nil {
		return "", fmt.Errorf("error happens while selecting: %w", err)
	}

	return result, nil
}

func (s SliderService) SaveSlider(form *multipart.Form, ctx *fiber.Ctx) error {

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

	if err := s.Repo.InsertSlider(slider, ctx.Context()); err != nil {
		if err := os.Remove(slider.Thumb); err != nil {
			return fmt.Errorf("error happens while remove file: %w", err)
		}

		return fmt.Errorf("error happens while inserting: %w", err)
	}

	return nil
}
