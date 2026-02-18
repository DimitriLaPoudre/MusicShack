package utils

import (
	"fmt"
	"io"
	"os"
)

func RenameForce(src, dst string) error {
	//linux
	if os.Rename(src, dst) == nil {
		return nil
	}

	//windows
	if err := os.Remove(dst); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("utils.RenameForce: os.Remove(%q): %w", dst, err)
	}
	if err := os.Rename(src, dst); err != nil {
		return fmt.Errorf("utils.RenameForce: os.Rename(%q, %q): %w", src, dst, err)
	}
	return nil
}

func RenameSoft(src, dst string) error {
	if _, err := os.Stat(dst); err == nil {
		return fmt.Errorf("utils.RenameSoft: os.Stat(%q) file exist", dst)
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("os.Stat(%q): %w", dst, err)
	}

	// rename standard
	if err := os.Rename(src, dst); err != nil {
		return fmt.Errorf("utils.RenameSoft: os.Rename(%q, %q): %w", src, dst, err)
	}
	return nil
}

func CopyTemporary(reader io.Reader) (*os.File, error) {
	file, err := os.CreateTemp("", "upload-*")
	if err != nil {
		return nil, fmt.Errorf("utils.CopyTemporary: os.CreateTemp: %w", err)
	}
	if _, err := io.Copy(file, reader); err != nil {
		return nil, fmt.Errorf("utils.CopyTemporary: io.Copy: %w", err)
	}
	if err := file.Sync(); err != nil {
		return nil, fmt.Errorf("utils.CopyTemporary: file.Sync: %w", err)
	}
	return file, nil
}

func DuplicateFile(srcPath, dstPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("utils.DuplicateFile: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("utils.DuplicateFile: %w", err)
	}
	defer func() {
		dst.Close()
		if err != nil {
			os.Remove(dstPath)
		}
	}()

	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("utils.DuplicateFile: %w", err)
	}

	if err = dst.Sync(); err != nil {
		return fmt.Errorf("utils.DuplicateFile: %w", err)
	}

	return nil
}

func RootDuplicateFile(root *os.Root, srcPath, dstPath string) error {
	src, err := root.Open(srcPath)
	if err != nil {
		return fmt.Errorf("utils.RootDuplicateFile: %w", err)
	}
	defer src.Close()

	dst, err := root.Create(dstPath)
	if err != nil {
		return fmt.Errorf("utils.RootDuplicateFile: %w", err)
	}
	defer func() {
		dst.Close()
		if err != nil {
			root.Remove(dstPath)
		}
	}()

	if _, err = io.Copy(dst, src); err != nil {
		return fmt.Errorf("utils.RootDuplicateFile: %w", err)
	}

	if err = dst.Sync(); err != nil {
		return fmt.Errorf("utils.RootDuplicateFile: %w", err)
	}

	return nil
}
