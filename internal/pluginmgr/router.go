package pluginmgr

import "github.com/go-chi/chi"

func (mgr *Manager) AddPluginsToRouter(router chi.Router) error {

	router.Route("/v1", func(r chi.Router) {
		r.Mount("/wallet", *mgr.walletRouter)
	})

	return nil
}
