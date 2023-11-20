package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Logger interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Err(message string)
}

type Processor struct {
	logger Logger
}

func NewProcessor(logger Logger) *Processor {
	return &Processor{
		logger: logger,
	}
}

func (p *Processor) CopyDirectory(sourceDir, destDir string) error {
	fInfoSource, err := p.getFileInfo(sourceDir)
	if err != nil {
		return p.logAndReturnError(fmt.Sprintf("failed to access source directory %q", sourceDir), err)
	}

	if !fInfoSource.IsDir() {
		return p.logAndReturnError("Source directory is not a directory", fmt.Errorf("source %q is not a directory", sourceDir))
	}

	err = os.MkdirAll(destDir, fInfoSource.Mode())
	if err != nil {
		return p.logAndReturnError(fmt.Sprintf("failed to create destination directory %q", destDir), err)
	} else {
		p.logger.Info(fmt.Sprintf("Created destination directory %q", destDir))
	}

	err = filepath.WalkDir(sourceDir, p.walkDirFunc(sourceDir, destDir))

	if err != nil {
		return p.logAndReturnError("failed to copy directory", err)
	}

	return nil
}

func (p *Processor) CopyFile(sourceFile, destFile string) error {
	fInfoSource, err := p.getFileInfo(sourceFile)
	if err != nil {
		return p.logAndReturnError(fmt.Sprintf("failed to access source file %q", sourceFile), err)
	}

	src, err := os.Open(sourceFile)
	if err != nil {
		return p.logAndReturnError(fmt.Sprintf("failed to open source file %q", sourceFile), err)
	}
	defer src.Close()

	dst, err := os.OpenFile(destFile, os.O_WRONLY|os.O_CREATE, fInfoSource.Mode())
	if err != nil {
		return p.logAndReturnError(fmt.Sprintf("failed to create destination file %q", destFile), err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return p.logAndReturnError(fmt.Sprintf("failed to copy file %q to %q", sourceFile, destFile), err)
	}

	return nil
}

func (p *Processor) logAndReturnError(message string, err error) error {
	errorMessage := fmt.Sprintf("%s: %v", message, err)
	p.logger.Err(errorMessage)

	return err
}

func (p *Processor) getFileInfo(src string) (os.FileInfo, error) {
	fInfo, err := os.Stat(src)
	if err != nil {
		if os.IsNotExist(err) {
			return fInfo, fmt.Errorf("source %q does not exist", src)
		}

		if os.IsPermission(err) {
			return fInfo, fmt.Errorf("source %q is not accessible", src)
		}

		return fInfo, fmt.Errorf("failed to access source %q", src)
	}

	return fInfo, nil
}

func (p *Processor) walkDirFunc(sourceDir, destDir string) func(path string, de os.DirEntry, err error) error {
	return func(path string, de os.DirEntry, err error) error {
		if err != nil {
			return p.logAndReturnError(fmt.Sprintf("failed to access path %s", path), err)
		}

		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return p.logAndReturnError(fmt.Sprintf("failed to get relative path for %q", path), err)
		}

		destPath := filepath.Join(destDir, relPath)
		if de.IsDir() {
			info, err := de.Info()
			if err != nil {
				return p.logAndReturnError(fmt.Sprintf("failed to get info for directory %q", path), err)
			}

			err = os.MkdirAll(destPath, info.Mode())
			if err != nil {
				return p.logAndReturnError(fmt.Sprintf("failed to create directory %q", destPath), err)
			} else {
				p.logger.Info(fmt.Sprintf("Created directory %q", destPath))
			}
		} else {
			err = p.CopyFile(path, destPath)
			if err != nil {
				return p.logAndReturnError(fmt.Sprintf("failed to copy file %q", path), err)
			} else {
				p.logger.Info(fmt.Sprintf("Copied file %q", path))
			}
		}

		return nil
	}
}
