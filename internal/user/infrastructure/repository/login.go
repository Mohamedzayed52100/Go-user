package repository

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"aidanwoods.dev/go-paseto"
	"github.com/goplaceapp/goplace-common/pkg/auth"
	"github.com/goplaceapp/goplace-common/pkg/errorhelper"
	"github.com/goplaceapp/goplace-common/pkg/meta"
	"github.com/goplaceapp/goplace-common/pkg/utils"
	userProto "github.com/goplaceapp/goplace-user/api/v1"
	tenantDomain "github.com/goplaceapp/goplace-user/pkg/tenantservice/domain"
	"github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	extUserDomain "github.com/goplaceapp/goplace-user/pkg/userservice/domain"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (r *UserRepository) Login(ctx context.Context, req *userProto.LoginRequest) (*userProto.LoginResponse, error) {
	var (
		getUser *extUserDomain.User
		err     error
	)

	if req.GetParams().GetEmail() == "" || req.GetParams().GetPassword() == "" {
		res, err := r.loginWithClientCredentials(ctx, req)
		if err != nil {
			return nil, err
		}

		return &userProto.LoginResponse{
			Result: &userProto.LoginResult{
				AccessToken: res.GetResult().GetAccessToken(),
				ExpiresAt:   res.GetResult().GetExpiresAt(),
			},
		}, nil
	} else if req.GetParams().GetEmail() != "" && req.GetParams().GetPassword() != "" {
		getUser, err = r.GetUserByEmail(strings.ToLower(req.GetParams().GetEmail()))
		if err != nil {
			return nil, status.Error(http.StatusUnauthorized, errorhelper.ErrInvalidEmailOrPassword)
		}
	} else {
		return nil, status.Error(http.StatusUnauthorized, errorhelper.ErrInvalidEmailOrPassword)
	}

	tenant, err := r.GetTenantByID(getUser.TenantID)
	if err != nil {
		return nil, status.Error(http.StatusNotFound, errorhelper.ErrThisOrganizationDoesNotExistInOurSystem)
	}

	ctx = context.WithValue(ctx, meta.TenantDBNameContextKey.String(), tenant.DbName)

	password := getUser.Password
	if err := bcrypt.CompareHashAndPassword([]byte(password),
		[]byte(req.GetParams().GetPassword())); err != nil {
		return nil, status.Error(http.StatusUnauthorized, errorhelper.ErrInvalidEmailOrPassword)
	}

	sessionBytes := make([]byte, 32)
	rand.Read(sessionBytes)
	sessionId := hex.EncodeToString(sessionBytes)
	_, err = r.GetRolePermissionsByRoleID(ctx, getUser.RoleID)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	generateToken, err := generateTokenForUser(sessionId, getUser, tenant)
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	if err := r.GetTenantDBConnection(ctx).Create(&extUserDomain.UserSession{
		UserID:     getUser.ID,
		SessionID:  sessionId,
		LoginTime:  time.Now().UTC(),
		LogoutTime: nil,
	}).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if os.Getenv("ENVIRONMENT") == meta.ProdEnvironment {
		otp, _, err := r.SendUserOTP(
			ctx,
			getUser,
			"sms",
			"login",
			generateToken,
			4,
		)
		if err != nil {
			return nil, status.Error(http.StatusInternalServerError, err.Error())
		}

		if err := r.SendSMSOTP(getUser, otp); err != nil {
			return nil, err
		}
	} else {
		if err := r.GetSharedDB().Create(&domain.UserOtp{
			UserID:       getUser.ID,
			Code:         "1234",
			Token:        generateToken,
			Type:         "login",
			VerifyMethod: "sms",
			ExpiresAt:    time.Now().UTC().Add(15 * time.Minute),
		}).Error; err != nil {
			if err := r.GetSharedDB().
				Where("user_id = ?", getUser.ID).
				Updates(&domain.UserOtp{
					Code:      "1234",
					Token:     generateToken,
					ExpiresAt: time.Now().UTC().Add(15 * time.Minute),
				}).Error; err != nil {
				return nil, status.Error(http.StatusInternalServerError, err.Error())
			}
		}
	}

	res := &userProto.LoginResponse{
		Result: &userProto.LoginResult{
			AccessToken: generateToken,
			ExpiresAt:   timestamppb.New(time.Now().Add(24 * time.Hour)),
			Setup:       getUser.SetupCompleted,
		},
	}

	return res, nil
}

func generateTokenForUser(sessionId string, getUser *extUserDomain.User, tenant *tenantDomain.Tenant) (string, error) {
	token := paseto.NewToken()
	token.SetString("session_id", sessionId)
	token.Set("role_id", getUser.RoleID)
	token.Set("account", auth.AccountMetaData{
		Email:        getUser.Email,
		ClientDBName: tenant.DbName,
		ClientID:     tenant.ID,
	})
	token.SetExpiration(time.Now().Add(24 * time.Hour))

	symmetricKey, err := utils.GetSymmetricKeyFromBase64(os.Getenv("TOKEN_SYMMETRIC_KEY"))
	if err != nil {
		return "", status.Error(http.StatusInternalServerError, err.Error())
	}

	encryptedToken := token.V4Encrypt(symmetricKey, nil)

	return encryptedToken, nil
}

func (r *UserRepository) loginWithClientCredentials(ctx context.Context, req *userProto.LoginRequest) (*userProto.LoginResponse, error) {
	var (
		token  string
		err    error
		tenant *tenantDomain.Tenant
	)

	// convert hex to binary
	bytes, err := hex.DecodeString(req.GetParams().GetClientSecret())
	if err != nil {
		return nil, status.Error(http.StatusUnauthorized, "Invalid Credentials")
	}

	// convert binary to base64
	base64String := base64.StdEncoding.EncodeToString(bytes)

	// decrypt base64String using AES_ENCRYPTION_KEY
	decryptedClientSecret := utils.Decrypt(base64String, os.Getenv("AES_ENCRYPTION_KEY"))
	if decryptedClientSecret == "" {
		return nil, status.Error(http.StatusUnauthorized, "Invalid Credentials")
	}

	// convert decryptedClientSecret to json
	var clientSecretJson map[string]interface{}
	err = json.Unmarshal([]byte(decryptedClientSecret), &clientSecretJson)
	if err != nil {
		return nil, status.Error(http.StatusUnauthorized, "Invalid Credentials")
	}

	var branchId int64

	if clientSecretJson["branch_id"] != nil {
		if err := r.SharedDbConnection.
			First(&tenant, "id = ?", clientSecretJson["tenant_id"]).Error; err != nil {
			return nil, status.Error(http.StatusUnauthorized, "Invalid Credentials")
		}

		ctx = context.WithValue(ctx, meta.TenantDBNameContextKey.String(), tenant.DbName)
		if err := r.GetTenantDBConnection(ctx).
			First(&extUserDomain.BranchCredential{}, "branch_id = ? AND client_id = ? AND client_secret = ?", clientSecretJson["branch_id"], req.GetParams().GetClientID(), req.GetParams().GetClientSecret()).Error; err != nil {
			return nil, status.Error(http.StatusUnauthorized, "Invalid Credentials")
		}

		branchId, _ = strconv.ParseInt(clientSecretJson["branch_id"].(string), 10, 32)
	} else {
		if err := r.SharedDbConnection.
			Model(&tenant).
			Joins("JOIN tenant_credentials ON tenants.id = tenant_credentials.tenant_id").
			First(&tenant, "client_id = ? AND client_secret = ?",
				req.GetParams().GetClientID(),
				req.GetParams().GetClientSecret(),
			).Error; err != nil {
			return nil, status.Error(http.StatusUnauthorized, "Invalid Credentials")
		}
	}

	token, err = generateTokenForClient(tenant.DbName, int32(branchId))
	if err != nil {
		return nil, status.Error(http.StatusInternalServerError, err.Error())
	}

	return &userProto.LoginResponse{
		Result: &userProto.LoginResult{
			AccessToken: token,
			ExpiresAt:   timestamppb.New(time.Now().Add(24 * time.Hour)),
		},
	}, nil
}

func generateTokenForClient(dbName string, branchId int32) (string, error) {
	sessionBytes := make([]byte, 32)
	sessionId := hex.EncodeToString(sessionBytes)
	token := paseto.NewToken()
	token.SetString("session_id", sessionId)
	token.SetExpiration(time.Now().Add(24 * time.Hour))
	token.Set("account", auth.AccountMetaData{
		ClientDBName: dbName,
	})

	if branchId != 0 {
		token.Set("branch_id", branchId)
	}

	symmetricKey, err := utils.GetSymmetricKeyFromBase64(os.Getenv("TOKEN_SYMMETRIC_KEY"))
	if err != nil {
		return "", status.Error(http.StatusInternalServerError, err.Error())
	}

	encryptedToken := token.V4Encrypt(symmetricKey, nil)

	return encryptedToken, nil
}


