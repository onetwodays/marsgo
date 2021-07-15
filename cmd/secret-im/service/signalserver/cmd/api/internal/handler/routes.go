// Code generated by goctl. DO NOT EDIT.
package handler

import (
	"net/http"

	accounts "secret-im/service/signalserver/cmd/api/internal/handler/accounts"
	certificate "secret-im/service/signalserver/cmd/api/internal/handler/certificate"
	channel "secret-im/service/signalserver/cmd/api/internal/handler/channel"
	device "secret-im/service/signalserver/cmd/api/internal/handler/device"
	directory "secret-im/service/signalserver/cmd/api/internal/handler/directory"
	group "secret-im/service/signalserver/cmd/api/internal/handler/group"
	keepalive "secret-im/service/signalserver/cmd/api/internal/handler/keepalive"
	keys "secret-im/service/signalserver/cmd/api/internal/handler/keys"
	messages "secret-im/service/signalserver/cmd/api/internal/handler/messages"
	profile "secret-im/service/signalserver/cmd/api/internal/handler/profile"
	profilekey "secret-im/service/signalserver/cmd/api/internal/handler/profilekey"
	provision "secret-im/service/signalserver/cmd/api/internal/handler/provision"
	textsecret "secret-im/service/signalserver/cmd/api/internal/handler/textsecret"
	"secret-im/service/signalserver/cmd/api/internal/svc"

	"github.com/tal-tech/go-zero/rest"
)

func RegisterHandlers(engine *rest.Server, serverCtx *svc.ServiceContext) {
	engine.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/index",
				Handler: IndexHandler(serverCtx),
			},
		},
	)

	engine.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/v1/accounts/:transport/code/:number",
				Handler: accounts.GetCodeReqHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/v1/accounts/code/:verificationCode",
				Handler: accounts.VerifyAccountHandler(serverCtx),
			},
		},
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.CheckBasicAuth},
			[]rest.Route{
				{
					Method:  http.MethodPut,
					Path:    "/v1/accounts/attributes",
					Handler: accounts.SetattributesHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/v1/accounts/pin",
					Handler: accounts.SetPinHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/v1/accounts/pin",
					Handler: accounts.DelPinHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/v1/accounts/registration_lock",
					Handler: accounts.SetRegLockHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/v1/accounts/registration_lock",
					Handler: accounts.DelRegLockHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/v1/accounts/name",
					Handler: accounts.SetDeviceNameHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/v1/accounts/signaling_key",
					Handler: accounts.DelSignlingKeyHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/accounts/whoami",
					Handler: accounts.GetWhoamiHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/v1/accounts/gcm",
					Handler: accounts.SetGcmRegistrationIDHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/v1/accounts/gcm",
					Handler: accounts.DelGcmRegistrationIDHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/v1/accounts/username",
					Handler: accounts.DeleteUserNameHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/v1/accounts/username/:username",
					Handler: accounts.SetUserNameHandler(serverCtx),
				},
			}...,
		),
	)

	engine.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/v1/textsecret/login",
				Handler: textsecret.AdxUserLoginHandler(serverCtx),
			},
		},
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.UserNameCheck},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/api/v1/textsecret/ws",
					Handler: textsecret.AdxUserWSHandler(serverCtx),
				},
			}...,
		),
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.CheckBasicAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/v1/messages",
					Handler: messages.GetMsgsHandler(serverCtx),
				},
			}...,
		),
	)

	engine.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPut,
				Path:    "/v1/messages/:destination",
				Handler: messages.PutMsgsHandler(serverCtx),
			},
		},
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.CheckBasicAuth},
			[]rest.Route{
				{
					Method:  http.MethodPut,
					Path:    "/v1/profile/:accountName",
					Handler: profilekey.PutProfileKeyHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/profile/:accountName",
					Handler: profilekey.GetProfileKeyHandler(serverCtx),
				},
			}...,
		),
	)

	engine.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/v1/keepalive",
				Handler: keepalive.GetKeepAliveHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/v1/keepalive/provisioning",
				Handler: keepalive.GetProvisioningKeepAliveHandler(serverCtx),
			},
		},
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.CheckBasicAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/v1/certificate/delivery",
					Handler: certificate.DeliveryHandler(serverCtx),
				},
			}...,
		),
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.CheckBasicAuth},
			[]rest.Route{
				{
					Method:  http.MethodPut,
					Path:    "/v2/keys",
					Handler: keys.PutKeysHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v2/keys",
					Handler: keys.GetKeyCountHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/v2/keys/signed",
					Handler: keys.SetSignedKeyHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v2/keys/signed",
					Handler: keys.GetSignedKeyHandler(serverCtx),
				},
			}...,
		),
	)

	engine.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/v2/keys/:identifier/:deviceId",
				Handler: keys.GetDeviceKeysHandler(serverCtx),
			},
		},
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.CheckBasicAuth},
			[]rest.Route{
				{
					Method:  http.MethodPut,
					Path:    "/v1/profile/name/:name",
					Handler: profile.SetNameHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/v1/profile",
					Handler: profile.SetProfileHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/profile/username/:username",
					Handler: profile.GetProfileByUserNameHandler(serverCtx),
				},
			}...,
		),
	)

	engine.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/v1/profile/:uuid/:version",
				Handler: profile.GetProfileByUuidHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/v1/profile/:uuid/:version/:credentialRequest",
				Handler: profile.GetProfileByUuidCredentiaHandler(serverCtx),
			},
		},
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.CheckBasicAuth},
			[]rest.Route{
				{
					Method:  http.MethodPut,
					Path:    "/v1/provisioning/:destination",
					Handler: provision.SendProvisioningMessageHandler(serverCtx),
				},
			}...,
		),
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.CheckBasicAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/v1/directory/auth",
					Handler: directory.GetAuthTokenHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/directory/:token",
					Handler: directory.GetTokenPresenceHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/v1/directory/tokens",
					Handler: directory.GetContactIntersectionHandler(serverCtx),
				},
			}...,
		),
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.CheckBasicAuth},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/v1/devices",
					Handler: device.GetDevicesHandler(serverCtx),
				},
				{
					Method:  http.MethodDelete,
					Path:    "/v1/devices/:device_id",
					Handler: device.DelDeviceHandler(serverCtx),
				},
				{
					Method:  http.MethodGet,
					Path:    "/v1/devices/provisioning/code",
					Handler: device.CreateDeviceTokenHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/v1/devices/unauthenticated_delivery",
					Handler: device.SetUnauthenticatedDeliveryHandler(serverCtx),
				},
				{
					Method:  http.MethodPut,
					Path:    "/v1/devices/:verification_code",
					Handler: device.VerifyDeviceTokenHandler(serverCtx),
				},
			}...,
		),
	)

	engine.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{serverCtx.CheckBasicAuth},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/v1/channels",
					Handler: channel.CreateChannelHandler(serverCtx),
				},
			}...,
		),
	)

	engine.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPut,
				Path:    "/v1/groups",
				Handler: group.CreateGroupHandler(serverCtx),
			},
		},
	)
}
