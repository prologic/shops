package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

type URI struct {
	Type string
	User string
	Host string
	Port string
}

func (u URI) HostAddr() string {
	return fmt.Sprintf("%s:%s", u.Host, u.Port)
}

func (u URI) String() string {
	return fmt.Sprintf("%s://%s@%s:%s", u.Type, u.User, u.Host, u.Port)
}

func ParseURI(uri, user, port string) (u URI) {
	u.Type = "ssh"
	u.User = user
	u.Port = port

	tokens := strings.Split(uri, "://")
	if len(tokens) == 2 {
		u.Type = strings.ToLower(tokens[0])

		tokens := strings.Split(tokens[1], "@")
		if len(tokens) == 2 {
			u.User = tokens[0]
			tokens := strings.Split(tokens[1], ":")
			if len(tokens) == 2 {
				u.Host = tokens[0]
				u.Port = tokens[1]
			} else {
				u.Host = tokens[0]
			}
		} else {
			u.User = user
			tokens := strings.Split(tokens[0], ":")
			if len(tokens) == 2 {
				u.Host = tokens[0]
				u.Port = tokens[1]
			} else {
				u.Host = tokens[0]
			}
		}
	} else {
		u.Host = tokens[0]
	}

	return
}

func ParseURIs(args []string, user, port string) []URI {
	var uris []URI

	for _, arg := range args {
		uris = append(uris, ParseURI(arg, user, port))
	}

	return uris
}

func readLines(r io.Reader) (lines []string, err error) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	err = scanner.Err()

	return
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func CopyDirectory(scrDir, dest string) error {
	entries, err := ioutil.ReadDir(scrDir)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		sourcePath := filepath.Join(scrDir, entry.Name())
		destPath := filepath.Join(dest, entry.Name())

		fileInfo, err := os.Stat(sourcePath)
		if err != nil {
			return err
		}

		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("failed to get raw syscall.Stat_t data for '%s'", sourcePath)
		}

		switch fileInfo.Mode() & os.ModeType {
		case os.ModeDir:
			if err := CreateIfNotExists(destPath, 0755); err != nil {
				return err
			}
			if err := CopyDirectory(sourcePath, destPath); err != nil {
				return err
			}
		case os.ModeSymlink:
			if err := CopySymLink(sourcePath, destPath); err != nil {
				return err
			}
		default:
			if _, err := CopyFile(sourcePath, destPath); err != nil {
				return err
			}
		}

		if err := os.Lchown(destPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return err
		}

		isSymlink := entry.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if err := os.Chmod(destPath, entry.Mode()); err != nil {
				return err
			}
		}
	}
	return nil
}

func Exists(filePath string) bool {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateIfNotExists(dir string, perm os.FileMode) error {
	if Exists(dir) {
		return nil
	}

	if err := os.MkdirAll(dir, perm); err != nil {
		return fmt.Errorf("failed to create directory: '%s', error: '%s'", dir, err.Error())
	}

	return nil
}

func CopySymLink(source, dest string) error {
	link, err := os.Readlink(source)
	if err != nil {
		return err
	}
	return os.Symlink(link, dest)
}
