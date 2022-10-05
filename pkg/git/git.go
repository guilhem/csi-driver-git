package git

import (
	"context"
	"fmt"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// type gitVolume struct {
// 	VolName string `json:"volName"`
// 	VolID   string `json:"volID"`
// 	// VolSize int64  `json:"volSize"`
// 	VolPath string `json:"volPath"`
// 	// VolGitPath    string            `json:"volGitPath"`
// 	repository *git.Repository `json:"-"`
// }

// // createVolume create the directory for the git volume.
// // It returns the volume path or err if one occurs.
// func ManageGitVolume(ctx context.Context, volID, name, repo, branch string, volAccessType accessType, ephemeral bool) (*gitVolume, error) {
// 	path := getVolumePath(volID)

// 	switch volAccessType {
// 	case mountAccess:
// 		err := os.MkdirAll(path, 0777)
// 		if err != nil {
// 			return nil, err
// 		}
// 	default:
// 		return nil, fmt.Errorf("unsupported access type %v", volAccessType)
// 	}

// 	repository, err := clone(ctx, repo, branch, path)
// 	if err != nil {
// 		return nil, fmt.Errorf("can't clone: %w", err)
// 	}

// 	gitVol := &gitVolume{
// 		VolID:   volID,
// 		VolName: name,
// 		VolPath: path,
// 		// VolGitPath:    gitPath,
// 		repository:    repository,
// 		VolAccessType: volAccessType,
// 		Ephemeral:     ephemeral,
// 	}
// 	gitVolumes[volID] = gitVol
// 	return gitVol, nil
// }

func Clone(ctx context.Context, repo string, branch string, path string) (*git.Repository, error) {
	// gitStorer := filesystem.NewStorage(osfs.New(path), cache.NewObjectLRUDefault())

	cloneOptions := &git.CloneOptions{
		URL:               repo,
		ReferenceName:     plumbing.NewBranchReferenceName(branch),
		SingleBranch:      true,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	if err := cloneOptions.Validate(); err != nil {
		return nil, fmt.Errorf("can't validate clone options: %v", err)
	}

	repository, err := git.PlainCloneContext(ctx, path, false, cloneOptions)
	if err != nil {
		return nil, fmt.Errorf("can't plain clone: %w", err)
	}

	return repository, nil
}
