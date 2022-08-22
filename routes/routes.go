package routes

import (
	"database/sql"
	"net/http"

	"github.com/code-golf/code-golf/middleware"
	"github.com/go-chi/chi/v5"
)

func Router(db *sql.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(
		middleware.Logger,
		errorMiddleware,
		middleware.Recoverer,
		middleware.RedirectHost,
		middleware.Static,
		middleware.RedirectSlashes,
		middleware.Compress(5),
		// middleware.Downtime,
		middleware.Database(db),
		middleware.Golfer,
	)

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	r.Get("/", indexGET)
	r.Get("/{hole}", holeGET)
	r.Get("/about", aboutGET)
	r.With(middleware.AdminArea).Route("/admin", func(r chi.Router) {
		r.Get("/", adminGET)
		r.Get("/solutions", adminSolutionsGET)
		r.Get("/solutions/run", adminSolutionsRunGET)
	})
	r.With(middleware.API).Route("/api", func(r chi.Router) {
		r.Get("/", apiGET)
		r.Get("/cheevos", apiCheevosGET)
		r.Get("/cheevos/{cheevo}", apiCheevoGET)
		r.Get("/holes", apiHolesGET)
		r.Get("/holes/{hole}", apiHoleGET)
		r.Get("/langs", apiLangsGET)
		r.Get("/langs/{lang}", apiLangGET)
		r.Get(
			"/mini-rankings/{hole}/{lang}/{scoring:bytes|chars}/{view:top|me|following}",
			apiMiniRankingsGET,
		)
		r.Get("/panic", apiPanicGET)
		r.Get("/suggestions/golfers", apiSuggestionsGolfersGET)
	})
	r.Get("/callback", callbackGET)
	r.Get("/callback/dev", callbackDevGET)
	r.Get("/feeds", feedsGET)
	r.Get("/feeds/{feed:atom|json|rss}", feedGET)
	r.With(middleware.GolferArea).Route("/golfer", func(r chi.Router) {
		r.Post("/cancel-delete", golferCancelDeletePOST)
		r.Get("/connect/{connection}", golferConnectGET)
		r.Post("/delete", golferDeletePOST)
		r.Post("/delete-solution", golferDeleteSolutionPOST)
		r.Get("/disconnect/{connection}", golferDisconnectGET)
		r.Get("/export", golferExportGET)
		r.Get("/settings", golferSettingsGET)
		r.Post("/settings", golferSettingsPOST)
	})
	r.With(middleware.GolferInfo).Route("/golfers/{name}", func(r chi.Router) {
		r.Get("/", golferGET)
		r.Post("/{action:follow|unfollow}", golferActionPOST)
		r.Get("/cheevos", golferCheevosGET)
		r.Get("/holes", golferHolesGET)
		r.Get("/holes/{scoring}", golferHolesGET)
		r.Get("/{hole}/{lang}/{scoring}", golferSolutionGET)
		// r.Post("/{hole}/{lang}/{scoring}", golferSolutionPOST)
	})
	r.Get("/healthz", healthzGET)
	r.Get("/ideas", ideasGET)
	r.Get("/log-out", logOutGET)
	r.Get("/random", randomGET)
	r.Route("/rankings", func(r chi.Router) {
		// Redirect some old URLs that got out.
		r.Get("/", redir("/rankings/holes/all/all/bytes"))
		r.Get("/holes", redir("/rankings/holes/all/all/bytes"))
		r.Get("/holes/all/all/all", redir("/rankings/holes/all/all/bytes"))
		r.Get("/langs/bytes", redir("/rankings/langs/all/bytes"))
		r.Get("/langs/chars", redir("/rankings/langs/all/chars"))
		r.Get("/medals", redir("/rankings/medals/all/all/all"))

		r.Get("/cheevos", rankingsCheevosGET)
		r.Get("/cheevos/all", redir("/rankings/cheevos"))
		r.Get("/cheevos/{cheevo}", rankingsCheevosGET)

		r.Get("/holes/{hole}/{lang}/{scoring}", rankingsHolesGET)
		r.Get("/recent-holes/{lang}/{scoring}", rankingsHolesGET)

		r.Get("/medals/{hole}/{lang}/{scoring}", rankingsMedalsGET)

		r.Get("/langs/{lang}/{scoring}", rankingsLangsGET)
		r.Get("/solutions", rankingsSolutionsGET)
	})
	r.Route("/recent", func(r chi.Router) {
		r.Get("/", redir("/recent/solutions/all/all/bytes"))
		r.Get("/{lang}", recentGET)

		r.Get("/golfers", recentGolfersGET)
		r.Get("/solutions/{hole}/{lang}/{scoring}", recentSolutionsGET)
	})
	r.Get("/scores/{hole}/{lang}", scoresGET)
	r.Get("/scores/{hole}/{lang}/all", scoresAllGET)
	r.Get("/scores/{hole}/{lang}/{scoring}", scoresGET)
	r.Get("/scores/{hole}/{lang}/{scoring}/{page}", scoresGET)
	r.Get("/sitemap.xml", sitemapGET)
	r.Post("/solution", solutionPOST)
	r.Get("/stats", statsGET)
	r.Get("/users/{name}", userGET)
	r.Get("/wiki", wikiGET)
	r.Get("/wiki/*", wikiPageGET)

	return r
}
