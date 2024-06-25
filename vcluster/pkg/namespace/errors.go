package namespace

import (
	"fmt"

	"github.com/loft-sh/workshop-kcd-munich-2024/vcluster/pkg/auth"
)

var ErrUserNotFound = fmt.Errorf("namespace: %w", auth.ErrNotFound)
