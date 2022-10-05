package node

import (
	"context"
	"log"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/go-playground/validator/v10"
	"github.com/guilhem/csi-runtime/node"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/guilhem/csi-driver-git/pkg/git"
)

func New(nodeID string, maxVolumesPerNode int64) (*node.Server, error) {
	ns := &Node{}
	return node.New(nodeID, maxVolumesPerNode, ns)
}

type Node struct {
}

var ErrBlockNotImplemented = status.Error(codes.Unimplemented, "block not implemented")

func (ns *Node) ValidateCapability(vCap *csi.VolumeCapability) error {
	if blk := vCap.GetBlock(); blk != nil {
		return ErrBlockNotImplemented
	}

	return nil
}

type VolumeContext struct {
	Repo   string `mapstructure:"repo" validate:"required"`
	Branch string `mapstructure:"repo"`
}

func (ns *Node) Populate(ctx context.Context, path string, volumeContext map[string]string, controllerContext map[string]string, secrets map[string]string) error {
	log.Print("populate")

	var result VolumeContext
	if err := mapstructure.Decode(volumeContext, &result); err != nil {
		return err
	}

	validate := validator.New()
	if err := validate.Struct(result); err != nil {
		return nil
	}

	if _, err := git.Clone(ctx, result.Repo, result.Branch, path); err != nil {
		return err
	}

	return nil
}
