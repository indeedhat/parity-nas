package parinas

import (
	"net/http"

	"github.com/indeedhat/parity-nas/internal/auth"
	"github.com/indeedhat/parity-nas/internal/config"
	"github.com/indeedhat/parity-nas/internal/env"
	"github.com/indeedhat/parity-nas/internal/servermux"
	"github.com/indeedhat/parity-nas/internal/sysmon"
)

func BuildRoutes(serverCfg servermux.ServerConfig) *http.ServeMux {
	r := servermux.NewRouter(serverCfg)

	public := r.Group("/api", auth.IsGuestMiddleware)
	{
		public.Post("/auth/login", auth.LoginController)
	}

	private := r.Group("/api", auth.IsLoggedInMiddleware)
	{
		private.Get("/auth/verify", auth.VerifyLoginController)
		private.Get("/system/monitor", sysmon.LiveMonitorController)
	}

	if env.DebugMode.Get() {
		debug := r.Group("/api")
		{
			debug.Get("/debug/config", config.ViewConfigController)
		}
	}

	return r.ServerMux()
}
