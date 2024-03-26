package v1

import (
	"github.com/ptsypyshev/gb-golang-level3-new/pkg/pb"
)

type usersClient interface {
	pb.UserServiceClient
}

type linksClient interface {
	pb.LinkServiceClient
}
