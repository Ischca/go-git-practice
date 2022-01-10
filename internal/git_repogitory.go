package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type GitRepository struct {
	worktree string
	gitdir   string
	conf     ini.File
}

func NewGitRepostitory(path string, force bool) (repo *GitRepository, err error) {
	worktree := path
	gitdir, err := filepath.Abs(filepath.Join(path, ".git"))
	if err != nil {
		return
	}

	if !(force || isDir(&gitdir)) {
		err = fmt.Errorf("Not a Git repository %s", path)
		return
	}

	repo = &GitRepository{
		worktree: worktree,
		gitdir:   gitdir,
	}

	// Read configuration file in .git/config
	confFileName := "config"
	cf := repo.repoFile(&confFileName, false)

	if result, _ := isExists(cf); result {
		bytes, err := os.ReadFile(*cf)
		if err != nil {
			return nil, err
		}
		conf, err := ini.Load(bytes)
		if err != nil {
			return nil, err
		}
		repo.conf = *conf
	} else if !force {
		err = fmt.Errorf("Configuration file missing")
		return
	}

	if !force {
		vers := repo.conf.Section("core").Key("repositoryformatversion").MustInt()
		if vers != 0 {
			err = fmt.Errorf("Unsupported repositoryformatversion %d", vers)
			return
		}
	}

	return
}

/*
 Same as repo_path, but create dirname(*path) if absent.

 For example, repo_file(r, "refs", "remotes", "origin", "HEAD") will create .git/refs/remotes/origin.
*/
func (repo GitRepository) repoFile(path *string, mkdir bool) *string {
	if err := repo.repoDir(path, mkdir); err != nil {
		return repo.repoPath(path)
	} else {
		return new(string)
	}
}

/*
 Same as repo_path, but mkdir *path if absent if mkdir.
*/
func (repo GitRepository) repoDir(path *string, mkdir bool) error {

	p := repo.repoPath(path)

	if result, _ := isExists(p); result {
		if isDir(p) {
			return nil
		} else {
			return fmt.Errorf("Not a directory %s", *p)
		}
	}

	if mkdir {
		err := os.MkdirAll(*p, 0777)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("")
	}
}

/*
 Compute path under repo's gitdir.
*/
func (repo GitRepository) repoPath(path *string) *string {
	str := filepath.Join(repo.gitdir, *path)
	return &str
}

func isDir(path *string) bool {
	result, stat := isExists(path)
	return result && stat.IsDir()
}

func isExists(path *string) (bool, os.FileInfo) {
	stat, err := os.Stat(*path)
	return !os.IsNotExist(err), stat
}

// Create a new repository at path.
func RepoCreate(path string) error {
	repo, err := NewGitRepostitory(path, true)
	if err != nil {
		return err
	}

	// First, we make sure the path either doesn't exist or is an empty dir.
	if result, stat := isExists(&repo.worktree); result {
		if !stat.IsDir() {
			err = fmt.Errorf("%s is not a directory!", path)
			return nil
		}
		if dirs, err := os.ReadDir(repo.worktree); err != nil || len(dirs) != 0 {
			fmt.Println(err)
			err = fmt.Errorf("%s is not empty!", path)
			return err
		}
	} else {
		err := os.MkdirAll(repo.worktree, 0777)
		if err != nil {
			return err
		}
	}

	branches := "branchs"
	assert(repo.repoDir(&branches, true))
	objects := "objects"
	assert(repo.repoDir(&objects, true))
	refs := "refs"
	tags := refs + "/" + "tags"
	assert(repo.repoDir(&tags, true))
	heads := refs + "/" + "heads"
	assert(repo.repoDir(&heads, true))

	// .git/description
	description := "description"
	repo.writeFile(&description, func(f *os.File) {
		text := []byte("Unnamed repository; edit this file 'description' to name the repository.\n")
		f.Write(text)
	})
	// .git/HEAD
	HEAD := "HEAD"
	repo.writeFile(&HEAD, func(f *os.File) {
		text := []byte("ref: refs/heads/master\n")
		f.Write(text)
	})
	config := "config"
	repo.writeFile(&config, func(f *os.File) {
		config, err2 := defaultConfig()
		if err != nil {
			err = err2
		}
		config.WriteTo(f)
	})
	return err
}

func defaultConfig() (*ini.File, error) {
	config := ini.Empty()
	sec, err := config.NewSection("core")
	if err != nil {
		return config, err
	}
	sec.NewKey("repositoryformatversion", "0")
	sec.NewKey("filemode", "false")
	sec.NewKey("bare", "false")

	return config, nil
}

func assert(err error) {
	if err != nil {
		panic(err)
	}
}

func (repo GitRepository) writeFile(fileName *string, ioFunc func(f *os.File)) (err error) {
	f, err := os.Create(*repo.repoFile(fileName, false))
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()
	ioFunc(f)
	return err
}
