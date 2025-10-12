package middleware

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/samber/do/v2"
	"go.uber.org/zap"
)

// ProductiveFileSystemChecker provides real filesystem checking for production use
type ProductiveFileSystemChecker struct {
	logger *zap.Logger
}

// NewProductiveFileSystemChecker creates a new filesystem checker for DI
func NewProductiveFileSystemChecker(i do.Injector) (FileSystemChecker, error) {
	logger := do.MustInvoke[*zap.Logger](i)
	return &ProductiveFileSystemChecker{
		logger: logger,
	}, nil
}

// FileExists checks if a file exists at the given path
func (pfsc *ProductiveFileSystemChecker) FileExists(path string) bool {
	if path == "" {
		return false
	}
	
	info, err := os.Stat(path)
	if err != nil {
		pfsc.logger.Debug("File does not exist", 
			zap.String("path", path),
			zap.Error(err))
		return false
	}
	
	exists := !info.IsDir()
	pfsc.logger.Debug("File existence check", 
		zap.String("path", path),
		zap.Bool("exists", exists),
		zap.Bool("is_directory", info.IsDir()))
	
	return exists
}

// IsDirectory checks if the given path is a directory
func (pfsc *ProductiveFileSystemChecker) IsDirectory(path string) bool {
	if path == "" {
		return false
	}
	
	info, err := os.Stat(path)
	if err != nil {
		pfsc.logger.Debug("Directory does not exist", 
			zap.String("path", path),
			zap.Error(err))
		return false
	}
	
	isDir := info.IsDir()
	pfsc.logger.Debug("Directory check", 
		zap.String("path", path),
		zap.Bool("is_directory", isDir))
	
	return isDir
}

// WalkDirectory walks through a directory tree and calls walkFn for each file/directory
func (pfsc *ProductiveFileSystemChecker) WalkDirectory(root string, walkFn func(path string, isDir bool, err error) error) error {
	if root == "" {
		return walkFn("", false, fmt.Errorf("root path cannot be empty"))
	}
	
	pfsc.logger.Debug("Walking directory", zap.String("root", root))
	
	return filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			pfsc.logger.Debug("Walk error", 
				zap.String("path", path),
				zap.Error(err))
			return walkFn(path, false, err)
		}
		
		isDir := d.IsDir()
		pfsc.logger.Debug("Walking path", 
			zap.String("path", path),
			zap.Bool("is_directory", isDir))
		
		return walkFn(path, isDir, nil)
	})
}