package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/goplaceapp/goplace-common/pkg/grpchelper"
	"github.com/goplaceapp/goplace-common/pkg/logger"
	"github.com/goplaceapp/goplace-user/database"
	departmentAdapters "github.com/goplaceapp/goplace-user/internal/department/adapters/grpc"
	roleAdapters "github.com/goplaceapp/goplace-user/internal/role/adapters/grpc"
	tenantAdapters "github.com/goplaceapp/goplace-user/internal/tenant/adapters/grpc"
	userAdapters "github.com/goplaceapp/goplace-user/internal/user/adapters/grpc"
	"go.uber.org/zap"
)

type Service struct {
	Log                     *zap.SugaredLogger
	BaseCfg                 *grpchelper.BaseConfig
	PostgresService         *database.PostgresService
	UserServiceServer       *userAdapters.UserServiceServer
	RoleServiceServer       *roleAdapters.RoleServiceServer
	DepartmentServiceServer *departmentAdapters.DepartmentServiceServer
	TenantServiceServer     *tenantAdapters.TenantServiceServer
}

func New() *Service {
	// initialize logger
	log, err := logger.New(os.Getenv("LOG_LEVEL"))
	if err != nil {
		panic(fmt.Errorf("failed to initialize the logger, %w", err))
	}

	postgresService := database.NewService()

	userServiceServer := userAdapters.NewUserService()

	roleServiceServer := roleAdapters.NewRoleService()

	departmentServiceServer := departmentAdapters.NewDepartmentService()

	tenantServiceServer := tenantAdapters.NewTenantService()

	return &Service{
		Log:                     log,
		BaseCfg:                 &grpchelper.BaseConfig{},
		PostgresService:         postgresService,
		UserServiceServer:       userServiceServer,
		RoleServiceServer:       roleServiceServer,
		DepartmentServiceServer: departmentServiceServer,
		TenantServiceServer:     tenantServiceServer,
	}
}

func (s *Service) SetBaseConfig(cfg *grpchelper.BaseConfig) {
	s.BaseCfg = cfg
}

func (s *Service) GetLog() *zap.SugaredLogger {
	return s.Log
}

func (s *Service) GetLivenessHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		// service specific health state definition goes here
		w.WriteHeader(http.StatusOK)
	}
}

func (s *Service) GetReadinessHandler() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		// service specific ready state definition goes here
		w.WriteHeader(http.StatusOK)
	}
}
