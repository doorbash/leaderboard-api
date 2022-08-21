package domain

import "context"

type BannedIpsCache interface {
	IsThisIPBanned(ctx context.Context, ip string) (bool, error)
	BanThisIP(ctx context.Context, ip string) error
}
