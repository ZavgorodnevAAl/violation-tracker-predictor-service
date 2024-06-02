package predictor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func (p *Predictor) makeFullFilePath(id uint64) string {
	fullFilePath := fmt.Sprintf("%s/%d", p.videoCaptureFolderPath, id)
	return fullFilePath
}

func (p *Predictor) uploadFile(id uint64, bytes []byte) (err error) {
	// Создаем файл
	file, err := os.Create(p.makeFullFilePath(id))
	if err != nil {
		err = errors.Wrap(err, "[os.Create(path)]")
		return err
	}
	defer file.Close()

	// Записываем декодированные данные в файл
	err = os.WriteFile(p.makeFullFilePath(id), bytes, 0644)
	if err != nil {
		err = errors.Wrap(err, "[os.WriteFile(fullFilePath, bytes, 0644)]")
		return err
	}

	return
}

func (p *Predictor) deleteFile(id uint64) (err error) {
	fullFilePath := fmt.Sprintf("%s/%d", p.videoCaptureFolderPath, id)
	if err = os.Remove(fullFilePath); err != nil {
		err = errors.Wrap(err, "[os.Remove(fullFilePath)]")
		return
	}

	return
}

func (p *Predictor) removeFilesInFolder() {
	os.Mkdir(p.videoCaptureFolderPath, 0755)

	// Чтение всех файлов в указанной папке
	fileList, err := filepath.Glob(filepath.Join(p.videoCaptureFolderPath, "*"))
	if err != nil {
		err = errors.Wrap(err, "[filepath.Glob(filepath.Join(p.videoCaptureFolderPath, *))]")
		panic(err)
	}

	// Удаление каждого файла в папке
	for _, file := range fileList {
		if err = os.Remove(file); err != nil {
			err = errors.Wrap(err, "[os.Remove(file)]")
			panic(err)
		}
		p.log.Info(fmt.Sprintf("delete %s file", file))
	}
}
