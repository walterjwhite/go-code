package action

import (
	"context"

	"github.com/chromedp/cdproto/browser"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

func Grant(ctx context.Context, permissions []browser.PermissionType) {

	for _, permission := range permissions {
		descriptor := &browser.PermissionDescriptor{
			Name: string(permission),
		}
		logging.Error(browser.SetPermission(descriptor, browser.PermissionSettingGranted).Do(ctx), "grant permissions")
	}
}
