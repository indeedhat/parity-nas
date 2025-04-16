package parinas

import (
	"net/http"

	"github.com/indeedhat/parity-nas/internal/auth"
	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/internal/servermux"
	"github.com/indeedhat/parity-nas/internal/sysmon"
	"github.com/indeedhat/parity-nas/internal/tty"
)

func BuildRoutes(serverCfg servermux.ServerConfig) *http.ServeMux {
	r := servermux.NewRouter(serverCfg)

	public := r.Group("/api", auth.IsGuestMiddleware)
	{
		public.Post("/auth/login", auth.LoginController)
	}

	privateAny := r.Group("/api", auth.IsLoggedInMiddleware)
	{
		privateAny.Get("/auth/verify", auth.VerifyLoginController)
		privateAny.Get("/system/monitor", sysmon.LiveMonitorController)
	}

	privateAdmin := r.Group("/api", auth.UserHasPermissionMiddleware(auth.PermissionAdmin))
	{
		privateAdmin.Get("/debug/config", config.ViewConfigController)
		privateAdmin.Get("/system/tty", tty.TtyController)
	}

	return r.ServerMux()
}
