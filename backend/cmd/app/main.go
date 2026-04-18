package main

import (
	"github.com/AniketSrivastava1/recruit/backend/internal"
	"github.com/AniketSrivastava1/recruit/backend/internal/boot"
	"github.com/AniketSrivastava1/recruit/backend/internal/controllers"
	"github.com/AniketSrivastava1/recruit/backend/internal/db"
	"github.com/AniketSrivastava1/recruit/backend/internal/routes"
	"github.com/AniketSrivastava1/recruit/backend/internal/s3"
	"github.com/AniketSrivastava1/recruit/backend/internal/services"
	"github.com/AniketSrivastava1/recruit/backend/internal/stores"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	boot.LoadEnv()

	fx.New(
		fx.Provide(
			zap.NewDevelopment,
			boot.NewFirebaseAuth,
			db.NewDBConn,
			s3.NewS3Client,
			stores.NewUserStore,
			stores.NewAdminStore,
			stores.NewContestStore,
			stores.NewSubmissionStore,
			services.NewUserService,
			services.NewContestService,
			services.NewSubmissionService,
			controllers.NewUserController,
			controllers.NewContestController,
			controllers.NewSubmissionController,
			internal.NewEchoServer,
		),
		fx.Invoke(
			routes.AddUserRoutes,
			routes.AddContestRoutes,
			routes.AddAdminRoutes,
			routes.AddSubmissionRoutes,
			internal.StartEchoServer,
		),
	).Run()
}
