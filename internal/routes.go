package parinas

import (
	"net/http"

	"github.com/indeedhat/parity-nas/internal/auth"
	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/internal/logging"
	"github.com/indeedhat/parity-nas/internal/sysmon"
	"github.com/indeedhat/parity-nas/internal/tty"
	webproxy "github.com/indeedhat/parity-nas/internal/web_proxy"
	"github.com/indeedhat/parity-nas/pkg/server_mux"
)

const (
	PermissionPublic uint8 = iota
	PermissionGuest
	PermissionUser
	PermissionAdmin
)

func BuildRoutes(r servermux.Router, proxyCfg *config.WebProxyCfg) *http.ServeMux {
	logger := logging.New("router")
	r = r.Group("", logging.LoggingMiddleware(logger))

	r.HandleFunc("/"+proxyCfg.Prefix+"/", webproxy.WebProxyController)

	public := r.Group("/api", auth.IsGuestMiddleware)
	{
		public.HandleFunc("POST /auth/login", auth.LoginController)
	}

	privateAny := r.Group("/api", auth.IsLoggedInMiddleware)
	{
		privateAny.HandleFunc("GET /auth/verify", auth.VerifyLoginController)
		privateAny.HandleFunc("GET /system/monitor", sysmon.LiveMonitorController)
	}

	privateAdmin := r.Group("/api", auth.UserHasPermissionMiddleware(auth.PermissionAdmin))
	{
		privateAdmin.HandleFunc("GET /debug/config", config.ViewConfigController)
		privateAdmin.HandleFunc("GET /system/tty", tty.TtyController)
		privateAdmin.HandleFunc("GET /system/logs", logging.LiveMonitorLogsController)
	}

	return r.ServerMux()
}

func PluginRouter(r servermux.Router, permission uint8, pluginName string) servermux.Router {
	logger := logging.New("router").WithAttr("plugin", pluginName)

	var middleware servermux.Middleware

	switch permission {
	case PermissionPublic:
		middleware = nil
	case PermissionGuest:
		middleware = auth.IsGuestMiddleware
	case PermissionUser:
		middleware = auth.IsLoggedInMiddleware
	case PermissionAdmin:
		middleware = auth.UserHasPermissionMiddleware(auth.PermissionAdmin)
	default:
		logger.Fatalf("requested plugin router with invalid permission: %d", permission)
	}

	return r.Group("", logging.LoggingMiddleware(logger), middleware)
}
