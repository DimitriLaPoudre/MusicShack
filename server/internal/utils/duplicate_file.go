package utils

import (
	"fmt"
	"io"
	"os"
)

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

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("utils.DuplicateFile: %w", err)
	}

	err = dst.Sync()
	if err != nil {
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

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("utils.RootDuplicateFile: %w", err)
	}

	err = dst.Sync()
	if err != nil {
		return fmt.Errorf("utils.RootDuplicateFile: %w", err)
	}

	return nil
}
