/*
 * Copyright (c) Huawei Technologies Co., Ltd. 2020-2021. All rights reserved.
 * secPaver is licensed under the Mulan PSL v2.
 * You can use this software according to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND, EITHER EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR
 * PURPOSE.
 * See the Mulan PSL v2 for more details.
 */

package utils

import (
	"archive/zip"
	"compress/bzip2"
	"fmt"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"gitee.com/openeuler/secpaver/common/global"
	"strings"
	"syscall"
)

var (
	unixPathRegexp = regexp.MustCompile("^(/[^/]*)+$")
)

// PathExist method check path if exist or not
func PathExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, fmt.Errorf("fail to get file status")
}

// DirExist method check dir if exist or not
func DirExist(path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, fmt.Errorf("fail to get file status")
}

// CheckValidZipFile method check file is a valid zip or not
func CheckValidZipFile(file string) error {
	if err := CheckZipFileName(filepath.Base(file)); err != nil {
		return errors.Wrap(err, "fail to check zip filename")
	}

	rd, err := zip.OpenReader(file)
	if err != nil {
		return fmt.Errorf("fail to open zip file")
	}
	defer rd.Close()

	if len(rd.File) > maxFileNumberInZip {
		return fmt.Errorf(
			"invalid zip file, the number of files in zip should be less than %d", maxFileNumberInZip)
	}

	return nil
}

// SafeIoCopyN implements io.CopyN function with the copy size check
func SafeIoCopyN(dst io.Writer, src io.Reader, size int64) error {
	_, err := io.CopyN(dst, src, size)
	if err == io.EOF {
		return fmt.Errorf("fail to check copy size, copy size does not match file size")
	}

	return err
}

// CopyFile is a tool to copy file
func CopyFile(src, dst string) error {
	sfd, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("fail to open source file")
	}
	defer sfd.Close()

	dfd, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, global.DefaultFilePerm)
	if err != nil {
		return fmt.Errorf("fail to open target file")
	}
	defer dfd.Close()

	sInfo, err := sfd.Stat()
	if err != nil {
		return errors.Wrap(err, "fail to get source file information")
	}

	if err := SafeIoCopyN(dfd, sfd, sInfo.Size()); err != nil {
		return err
	}

	srcState, err := os.Stat(src)
	if err != nil {
		return errors.Wrap(err, "fail to get source file status")
	}

	if err := os.Chmod(dst, srcState.Mode()); err != nil {
		return errors.Wrap(err, "fail to chmod target file")
	}

	return nil
}

// ZipDir is a tool to compress a directory to a zip file
func ZipDir(dir, zipFile string) error {
	zipDir := filepath.Dir(zipFile)
	if err := os.MkdirAll(zipDir, global.DefaultDirPerm); err != nil {
		return fmt.Errorf("fail to create directory %s", filepath.Base(zipDir))
	}

	fz, err := os.OpenFile(zipFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, global.DefaultFilePerm)
	if err != nil {
		return fmt.Errorf("fail to create zip file %s", filepath.Base(zipFile))
	}
	defer fz.Close()

	w := zip.NewWriter(fz)
	defer w.Close()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fDest, err := w.Create(path[len(filepath.Dir(dir))+1:])
			if err != nil {
				return err
			}

			fSrc, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("fail to open %s", filepath.Base(path))
			}
			defer fSrc.Close()

			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				return err
			}
		}
		return nil
	})

	return err
}

// GetModTime get a file's mod time
func GetModTime(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return "N/A"
	}

	return info.ModTime().Format("2006-01-02 15:04:05")
}

// GetBodyFileName get the body name of a filename
// e.g. /dir/test.txt -> test
func GetBodyFileName(fileName string) string {
	ext := filepath.Ext(filepath.Base(fileName))
	return strings.TrimSuffix(filepath.Base(fileName), ext)
}

// FindAllSubDir return all first-level subdirectories in a directory
func FindAllSubDir(dir string) ([]string, error) {
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("fail to read directory")
	}

	subDirs := make([]string, 0, len(infos))

	for _, info := range infos {
		if info.IsDir() {
			subDirs = append(subDirs, filepath.Join(dir, info.Name()))
		}
	}

	return subDirs, nil
}

// IsAbsolutePath determine if a path is an absolute path
func IsAbsolutePath(path string) bool {
	return filepath.IsAbs(path)
}

// IsUnixFilePath checks if a path is an unix absolute path
func IsUnixFilePath(path string) bool {
	return unixPathRegexp.MatchString(path)
}

// GetUIDOfFile returns the uid of file, if fail to get the uid, returns an error
func GetUIDOfFile(path string) (uint32, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("fail to get file state")
	}

	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return stat.Uid, nil
	}

	return 0, fmt.Errorf("fail to get file uid")
}

// ReadBzip2 read a bzip file
func ReadBzip2(path string) (string, error) {
	fd, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("fail to open file %s", filepath.Base(path))
	}
	defer fd.Close()

	info, err := fd.Stat()
	if err != nil {
		return "", err
	}

	if info.Size() > maxFileSizeInBzip2 {
		return "", fmt.Errorf(
			"invalid bzip2 file %s, the file size must be smaller than %d", path, maxFileSizeInBzip2)
	}

	data, err := ioutil.ReadAll(bzip2.NewReader(fd))
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// SearchFileFromZipByName searches path in zip
func SearchFileFromZipByName(reader *zip.Reader, filename string) (string, error) {
	if reader == nil {
		return "", nil
	}

	pathSearched := ""

	if len(reader.File) > maxFileNumberInZip {
		return "", fmt.Errorf(
			"invalid zip file, the number of files in zip should be less than %d", maxFileNumberInZip)
	}

	for _, file := range reader.File {
		if err := CheckUnsafePath(file.Name); err != nil {
			return "", errors.Wrap(err, "unsafe path found in zip")
		}

		if filepath.Base(file.Name) == filename {
			if file.FileInfo().IsDir() {
				continue
			}

			if pathSearched != "" {
				return "", fmt.Errorf("multiple files with the same name %s exist", filename)
			}

			pathSearched = file.Name;
		}
	}

	return pathSearched, nil
}

// ExtractFileFromZip searches and returns file content from zip
func ExtractFileFromZip(reader *zip.Reader, filename string) ([]byte, error) {
	if reader == nil {
		return nil, nil
	}

	if len(reader.File) > maxFileNumberInZip {
		return nil, fmt.Errorf(
			"invalid zip file, the number of files in zip should be less than %d", maxFileNumberInZip)
	}

	for _, file := range reader.File {
		if err := CheckUnsafePath(file.Name); err != nil {
			return nil, errors.Wrap(err, "unsafe path found in zip")
		}

		if file.Name == filename {
			info := file.FileInfo()
			if info.IsDir() {
				continue
			}

			if info.Size() > maxFileSizeInZip {
				return nil, fmt.Errorf("oversize file %s in zip, the size must be less than %d",
					info.Name(), maxFileSizeInZip)
			}

			rd, err := file.Open()
			if err != nil {
				return nil, err
			}

			data, err := ioutil.ReadAll(rd)
			_ = rd.Close()

			return data, err
		}
	}

	return nil, fmt.Errorf("can't find %s in zip", filename)
}

// CheckFileSize checks if the file size is smaller than MaxFileSize
func CheckFileSize(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("fail to get file status")
	}

	if info.Size() > maxFileSize {
		return fmt.Errorf(
			"the file size must be smaller than %d", maxFileSize)
	}

	return nil
}

// WriteFile is a warp function for creating directory and writing data to file
func WriteFile(path string, data []byte, perm os.FileMode) error {
	path = filepath.Clean(path)
	if err := os.MkdirAll(filepath.Dir(path), global.DefaultDirPerm); err != nil {
		return fmt.Errorf("fail to create directory")
	}

	if err := ioutil.WriteFile(path, data, perm); err != nil {
		return fmt.Errorf("fail to write data to file")
	}

	return nil
}
