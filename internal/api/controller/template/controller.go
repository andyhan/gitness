// Copyright 2022 Harness Inc. All rights reserved.
// Use of this source code is governed by the Polyform Free Trial License
// that can be found in the LICENSE.md file for this repository.

package template

import (
	"github.com/harness/gitness/internal/auth/authz"
	"github.com/harness/gitness/internal/store"
	"github.com/harness/gitness/types/check"

	"github.com/jmoiron/sqlx"
)

type Controller struct {
	db            *sqlx.DB
	uidCheck      check.PathUID
	pathStore     store.PathStore
	templateStore store.TemplateStore
	authorizer    authz.Authorizer
	spaceStore    store.SpaceStore
}

func NewController(
	db *sqlx.DB,
	uidCheck check.PathUID,
	authorizer authz.Authorizer,
	templateStore store.TemplateStore,
	spaceStore store.SpaceStore,
) *Controller {
	return &Controller{
		db:            db,
		uidCheck:      uidCheck,
		templateStore: templateStore,
		authorizer:    authorizer,
		spaceStore:    spaceStore,
	}
}