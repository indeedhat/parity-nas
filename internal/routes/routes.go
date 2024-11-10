package routes

import (
	"net/http"

	"github.com/indeedhat/parity-nas/internal/routes/controllers"
)

func Build() *http.ServeMux {
	r := newRouter()

	public := r.Group("/api", isGuest)
	{
		public.Post("/auth/login", controllers.Login)
	}

	private := r.Group("/api", isLoggedIn)
	{
		private.Post("/auth/verify", controllers.VerifyLogin)
		private.Get("/system/monitor", controllers.LiveMonitor)
	}

	return r.mux
}
