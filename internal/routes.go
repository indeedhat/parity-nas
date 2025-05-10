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

func BuildRoutes(serverCfg servermux.ServerConfig, proxyCfg *config.WebProxyCfg) *http.ServeMux {
	logger := logging.New("router")

	r := servermux.NewRouter(
		serverCfg,
		logging.LoggingMiddleware(logger),
	)

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
