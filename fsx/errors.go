package fsx

import (
	"os"

	"github.com/boostgo/core/errorx"
)

var (
	ErrFailedChangeFilePermissions = errorx.New("fsx.file.change_permissions")
	ErrCreateFileDirectories       = errorx.New("fsx.file.create_directories")
	ErrReadFile                    = errorx.New("fsx.file.read")
	ErrOpenFile                    = errorx.New("fsx.file.open")
	ErrReadFileLines               = errorx.New("fsx.file.read.lines")
	ErrCreateFile                  = errorx.New("fsx.file.create")
	ErrCreateBackupFile            = errorx.New("fsx.file.create.backup")
	ErrAppendFile                  = errorx.New("fsx.file.append")
	ErrDeleteFile                  = errorx.New("fsx.file.delete")
	ErrStatFile                    = errorx.New("fsx.file.stat")
	ErrCopyFile                    = errorx.New("fsx.file.copy")
	ErrAtomicOperation             = errorx.New("fsx.file.atomic")
	ErrTempFile                    = errorx.New("fsx.file.temp")
	ErrFileLock                    = errorx.New("fsx.file.lock")
	ErrStreamOperation             = errorx.New("fsx.file.stream")
	ErrCompress                    = errorx.New("fsx.file.compress")
	ErrDecompress                  = errorx.New("fsx.file.decompress")
	ErrChecksum                    = errorx.New("fsx.file.checksum")
	ErrFileAlreadyLocked           = errorx.New("fsx.file.already_locked")
	ErrFileNotLocked               = errorx.New("fsx.file.not_locked")
	ErrInvalidArchive              = errorx.New("fsx.file.invalid_archive")

	ErrCreateDirectory            = errorx.New("fsx.file.create.directory")
	ErrCreateDirectories          = errorx.New("fsx.file.create.directories")
	ErrDeleteDirectory            = errorx.New("fsx.directory.delete")
	ErrDeleteDirectoryNotEmpty    = errorx.New("fsx.directory.delete.not_empty")
	ErrRenameDirectory            = errorx.New("fsx.directory.rename")
	ErrListDirectory              = errorx.New("fsx.directory.list")
	ErrReadDirectory              = errorx.New("fsx.directory.read")
	ErrStatDirectory              = errorx.New("fsx.directory.stat")
	ErrChangeDirectoryPermissions = errorx.New("fsx.directory.change_permissions")
	ErrDirectoryNotExist          = errorx.New("fsx.directory.not_exist")
	ErrNotDirectory               = errorx.New("fsx.directory.not_directory")
	ErrCopyDirectory              = errorx.New("fsx.directory.copy")
	ErrSyncDirectory              = errorx.New("fsx.directory.sync")
	ErrCompareDirectory           = errorx.New("fsx.directory.compare")
	ErrWalkDirectory              = errorx.New("fsx.directory.walk")
	ErrCalculateSize              = errorx.New("fsx.directory.calculate_size")
	ErrSourceNotDirectory         = errorx.New("fsx.directory.source_not_directory")
	ErrDestinationExists          = errorx.New("fsx.directory.destination_exists")

	ErrSearchFiles      = errorx.New("fsx.search.files")
	ErrSearchContent    = errorx.New("fsx.search.content")
	ErrInvalidPattern   = errorx.New("fsx.search.invalid_pattern")
	ErrInvalidRegex     = errorx.New("fsx.search.invalid_regex")
	ErrSearchDepthLimit = errorx.New("fsx.search.depth_limit")
)

type failedChangePermissionsContext struct {
	Path  string `json:"path"`
	Mode  string `json:"mode"`
	Error error  `json:"error"`
}

func newFailedChangePermissionsError(path string, mode os.FileMode, err error) error {
	return ErrFailedChangeFilePermissions.
		SetError(err).
		SetData(failedChangePermissionsContext{
			Path:  path,
			Mode:  mode.String(),
			Error: err,
		})
}

type pathErrorContext struct {
	Path  string `json:"path"`
	Error error  `json:"error"`
}

func newCreateFileDirectoriesError(path string, err error) error {
	return ErrCreateFileDirectories.
		SetError(err).
		SetData(pathErrorContext{
			Path:  path,
			Error: err,
		})
}

func newReadFileError(path string, err error) error {
	return ErrReadFile.SetError(err).SetData(pathErrorContext{
		Path:  path,
		Error: err,
	})
}

func newOpenFileError(path string, err error) error {
	return ErrOpenFile.
		SetError(err).
		SetData(pathErrorContext{
			Path:  path,
			Error: err,
		})
}

func newReadFileLinesError(path string, err error) error {
	return ErrReadFileLines.
		SetError(err).
		SetData(pathErrorContext{
			Path:  path,
			Error: err,
		})
}

func newCreateFile(path string, err error, mode os.FileMode) error {
	return ErrCreateFile.
		SetError(err).
		SetData(pathErrorContext{
			Path:  path,
			Error: err,
		})
}

func newCreateBackupFileError(path string, err error) error {
	return ErrCreateBackupFile.
		SetError(err).
		SetData(pathErrorContext{
			Path:  path,
			Error: err,
		})
}

func newCreateDirectory(path string, err error) error {
	return ErrCreateDirectory.
		SetError(err).
		SetData(pathErrorContext{
			Path:  path,
			Error: err,
		})
}

func newCreateDirectories(path string, err error) error {
	return ErrCreateDirectories.
		SetError(err).
		SetData(pathErrorContext{
			Path:  path,
			Error: err,
		})
}

func newAppendFile(path string, err error) error {
	return ErrAppendFile.
		SetError(err).
		SetData(pathErrorContext{
			Path:  path,
			Error: err,
		})
}

func newDeleteFile(path string, err error) error {
	return ErrDeleteFile.
		SetError(err).
		SetData(pathErrorContext{
			Path:  path,
			Error: err,
		})
}

func newStatFile(path string, err error) error {
	return ErrStatFile.
		SetError(err).
		SetData(pathErrorContext{
			Path:  path,
			Error: err,
		})
}

func newCopyFile(path string, err error) error {
	return ErrCopyFile.
		SetError(err).
		SetData(pathErrorContext{
			Path:  path,
			Error: err,
		})
}

type moveErrorContext struct {
	Source      string `json:"source"`
	Destination string `json:"destination"`
	Error       error  `json:"error"`
}
