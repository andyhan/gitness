// Code generated by Wire. DO NOT EDIT.

//go:build !wireinject && harness
// +build !wireinject,harness

package server

import (
	"context"
	"github.com/harness/gitness/events"
	"github.com/harness/gitness/gitrpc"
	server2 "github.com/harness/gitness/gitrpc/server"
	"github.com/harness/gitness/harness/auth/authn"
	"github.com/harness/gitness/harness/auth/authz"
	"github.com/harness/gitness/harness/bootstrap"
	"github.com/harness/gitness/harness/client"
	"github.com/harness/gitness/harness/router"
	"github.com/harness/gitness/harness/store"
	types2 "github.com/harness/gitness/harness/types"
	"github.com/harness/gitness/harness/types/check"
	"github.com/harness/gitness/internal/api/controller/githook"
	"github.com/harness/gitness/internal/api/controller/pullreq"
	"github.com/harness/gitness/internal/api/controller/repo"
	"github.com/harness/gitness/internal/api/controller/service"
	"github.com/harness/gitness/internal/api/controller/serviceaccount"
	"github.com/harness/gitness/internal/api/controller/space"
	"github.com/harness/gitness/internal/api/controller/user"
	webhook2 "github.com/harness/gitness/internal/api/controller/webhook"
	"github.com/harness/gitness/internal/cron"
	events3 "github.com/harness/gitness/internal/events/git"
	events2 "github.com/harness/gitness/internal/events/pullreq"
	router2 "github.com/harness/gitness/internal/router"
	"github.com/harness/gitness/internal/server"
	"github.com/harness/gitness/internal/services"
	pullreq2 "github.com/harness/gitness/internal/services/pullreq"
	"github.com/harness/gitness/internal/services/webhook"
	"github.com/harness/gitness/internal/store/cache"
	"github.com/harness/gitness/internal/store/database"
	"github.com/harness/gitness/internal/url"
	"github.com/harness/gitness/types"
)

// Injectors from harness.wire.go:

func initSystem(ctx context.Context, config *types.Config) (*system, error) {
	checkUser := check.ProvideUserCheck()
	typesConfig, err := types2.LoadConfig()
	if err != nil {
		return nil, err
	}
	serviceJWTProvider, err := client.ProvideServiceJWTProvider(typesConfig)
	if err != nil {
		return nil, err
	}
	aclClient, err := client.ProvideACLClient(serviceJWTProvider, typesConfig)
	if err != nil {
		return nil, err
	}
	authorizer := authz.ProvideAuthorizer(typesConfig, aclClient)
	db, err := database.ProvideDatabase(ctx, config)
	if err != nil {
		return nil, err
	}
	principalUIDTransformation := store.ProvidePrincipalUIDTransformation()
	principalStore := database.ProvidePrincipalStore(db, principalUIDTransformation)
	tokenStore := database.ProvideTokenStore(db)
	controller := user.NewController(checkUser, authorizer, principalStore, tokenStore)
	checkService := check.ProvideServiceCheck()
	serviceController := service.NewController(checkService, authorizer, principalStore)
	bootstrapBootstrap := bootstrap.ProvideBootstrap(config, controller, serviceController)
	tokenClient, err := client.ProvideTokenClient(serviceJWTProvider, typesConfig)
	if err != nil {
		return nil, err
	}
	userClient, err := client.ProvideUserClient(serviceJWTProvider, typesConfig)
	if err != nil {
		return nil, err
	}
	serviceAccountClient, err := client.ProvideServiceAccountClient(serviceJWTProvider, typesConfig)
	if err != nil {
		return nil, err
	}
	serviceAccount := check.ProvideServiceAccountCheck()
	pathTransformation := store.ProvidePathTransformation()
	pathStore := database.ProvidePathStore(db, pathTransformation)
	pathCache := cache.ProvidePathCache(pathStore, pathTransformation)
	spaceStore := database.ProvideSpaceStore(db, pathCache)
	repoStore := database.ProvideRepoStore(db, pathCache)
	serviceaccountController := serviceaccount.NewController(serviceAccount, authorizer, principalStore, spaceStore, repoStore, tokenStore)
	provider, err := url.ProvideURLProvider(config)
	if err != nil {
		return nil, err
	}
	pathUID := check.ProvidePathUIDCheck()
	spaceController := space.ProvideController(db, provider, pathUID, authorizer, pathStore, spaceStore, repoStore, principalStore)
	accountClient, err := client.ProvideAccountClient(serviceJWTProvider, typesConfig)
	if err != nil {
		return nil, err
	}
	authenticator, err := authn.ProvideAuthenticator(controller, tokenClient, userClient, typesConfig, serviceAccountClient, serviceaccountController, serviceController, spaceController, accountClient)
	if err != nil {
		return nil, err
	}
	gitrpcConfig := ProvideGitRPCClientConfig(config)
	gitrpcInterface, err := gitrpc.ProvideClient(gitrpcConfig)
	if err != nil {
		return nil, err
	}
	repoController := repo.ProvideController(config, db, provider, pathUID, authorizer, pathStore, repoStore, spaceStore, principalStore, gitrpcInterface)
	principalInfoView := database.ProvidePrincipalInfoView(db)
	principalInfoCache := cache.ProvidePrincipalInfoCache(principalInfoView)
	pullReqStore := database.ProvidePullReqStore(db, principalInfoCache)
	pullReqActivityStore := database.ProvidePullReqActivityStore(db, principalInfoCache)
	pullReqReviewStore := database.ProvidePullReqReviewStore(db)
	pullReqReviewerStore := database.ProvidePullReqReviewerStore(db, principalInfoCache)
	eventsConfig := ProvideEventsConfig(config)
	cmdable, err := ProvideRedis(config)
	if err != nil {
		return nil, err
	}
	eventsSystem, err := events.ProvideSystem(eventsConfig, cmdable)
	if err != nil {
		return nil, err
	}
	reporter, err := events2.ProvideReporter(eventsSystem)
	if err != nil {
		return nil, err
	}
	pullreqController := pullreq.ProvideController(db, provider, authorizer, pullReqStore, pullReqActivityStore, pullReqReviewStore, pullReqReviewerStore, repoStore, principalStore, gitrpcInterface, reporter)
	webhookStore := database.ProvideWebhookStore(db)
	webhookExecutionStore := database.ProvideWebhookExecutionStore(db)
	webhookConfig := ProvideWebhookConfig(config)
	readerFactory, err := events3.ProvideReaderFactory(eventsSystem)
	if err != nil {
		return nil, err
	}
	webhookService, err := webhook.ProvideService(ctx, webhookConfig, readerFactory, webhookStore, webhookExecutionStore, repoStore, provider, principalStore, gitrpcInterface)
	if err != nil {
		return nil, err
	}
	webhookController := webhook2.ProvideController(config, db, authorizer, webhookStore, webhookExecutionStore, repoStore, webhookService)
	eventsReporter, err := events3.ProvideReporter(eventsSystem)
	if err != nil {
		return nil, err
	}
	githookController := githook.ProvideController(db, authorizer, principalStore, repoStore, eventsReporter)
	apiHandler := router.ProvideAPIHandler(config, authenticator, accountClient, controller, spaceController, repoController, pullreqController, webhookController, githookController)
	gitHandler := router.ProvideGitHandler(config, provider, repoStore, authenticator, authorizer, gitrpcInterface)
	webHandler := router2.ProvideWebHandler(config)
	routerRouter := router2.ProvideRouter(config, apiHandler, gitHandler, webHandler)
	serverServer := server.ProvideServer(config, routerRouter)
	serverConfig := ProvideGitRPCServerConfig(config)
	server3, err := server2.ProvideServer(serverConfig)
	if err != nil {
		return nil, err
	}
	nightly := cron.NewNightly()
	eventsReaderFactory, err := events2.ProvideReaderFactory(eventsSystem)
	if err != nil {
		return nil, err
	}
	pullreqService, err := pullreq2.ProvideService(ctx, config, readerFactory, eventsReaderFactory, gitrpcInterface, db, repoStore, pullReqStore, pullReqActivityStore)
	if err != nil {
		return nil, err
	}
	servicesServices := services.ProvideServices(webhookService, pullreqService)
	serverSystem := newSystem(bootstrapBootstrap, serverServer, server3, nightly, servicesServices)
	return serverSystem, nil
}
