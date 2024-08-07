package routes

import (
	"net/http"

	"github.com/code-golf/code-golf/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func Router(db *sqlx.DB) http.Handler {
	r := chi.NewRouter()

	r.Use(
		middleware.Session,
		middleware.Logger,
		errorMiddleware,
		middleware.Recoverer,
		middleware.Static,
		middleware.RedirectSlashes,
		// middleware.Downtime,
		middleware.Database(db),
	)

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	// Simple routes that don't need middleware.Golfer.
	r.Get("/callback/dev", callbackDevGET)
	r.Get("/feeds/{feed:atom|json|rss}", feedGET)
	r.Get("/golfers/{name}/avatar", golferAvatarGET)
	r.Get(`/golfers/{name}/avatar/{size:\d+}`, golferAvatarGET)
	r.Get("/healthz", healthzGET)
	r.Post("/log-out", logOutPost)
	r.Get("/random", randomGET)
	r.Get("/sitemap.xml", sitemapGET)
	r.Get("/users/{name}", userGET)

	// HTML routes that need middleware.Golfer.
	r.Group(func(r chi.Router) {
		r.Use(middleware.Golfer)

		r.Get("/", homeGET)
		r.Get("/{hole}", holeGET)
		r.Get("/about", aboutGET)
		r.With(middleware.AdminArea).Route("/admin", func(r chi.Router) {
			r.Get("/", adminGET)
			r.Get("/banners", adminBannersGET)
			r.Post("/banners/{banner}", adminBannerPOST)
			r.Get("/solutions", adminSolutionsGET)
			r.Get("/solutions/run", adminSolutionsRunGET)
		})
		r.With(middleware.API).Route("/api", func(r chi.Router) {
			r.Get("/", apiGET)
			r.Get("/cheevos", apiCheevosGET)
			r.Get("/cheevos/{cheevo}", apiCheevoGET)
			r.Get("/golfers/{golfer}", apiGolferGET)
			r.Get("/holes", apiHolesGET)
			r.Get("/holes/{hole}", apiHoleGET)
			r.Get("/langs", apiLangsGET)
			r.Get("/langs/{lang}", apiLangGET)
			r.Get(
				"/mini-rankings/{hole}/{lang}/{scoring:bytes|chars}/{view:top|me|following}",
				apiMiniRankingsGET,
			)
			r.Get("/panic", apiPanicGET)
			r.Get("/solutions-log", apiSolutionsLogGET)
			r.Get("/suggestions/golfers", apiSuggestionsGolfersGET)
			r.Get("/wiki/*", apiWikiPageGET)

			// API routes that require a logged-in golfer.
			r.Group(func(r chi.Router) {
				r.With(middleware.GolferArea)

				r.Get("/notes", apiNotesGET)
				r.Delete("/notes/{hole}/{lang}", apiNoteDELETE)
				r.Get("/notes/{hole}/{lang}", apiNoteGET)
				r.Put("/notes/{hole}/{lang}", apiNotePUT)
			})
		})
		r.Get("/callback", callbackGET)
		r.Get("/feeds", feedsGET)
		r.With(middleware.GolferArea).Route("/golfer", func(r chi.Router) {
			r.Post("/cancel-delete", golferCancelDeletePOST)
			r.Get("/connect/{connection}", golferConnectGET)
			r.Post("/delete", golferDeletePOST)
			r.Post("/delete-solution", golferDeleteSolutionPOST)
			r.Get("/disconnect/{connection}", golferDisconnectGET)
			r.Get("/export", golferExportGET)
			r.Post("/hide-banner", golferHideBannerPOST)
			r.Get("/settings", golferSettingsGET)
			r.Post("/settings", golferSettingsPOST)
			r.Get("/settings/{page:delete-account}", golferSettingsDeleteAccountGET)
			r.Get("/settings/{page:export-data}", golferSettingsExportDataGET)
			r.Post("/settings/{page}", golferSettingsPagePOST)
			r.Post("/settings/{page}/reset", golferSettingsPageResetPOST)
		})
		r.With(middleware.GolferInfo).Route("/golfers/{name}", func(r chi.Router) {
			r.Get("/", golferGET)
			r.Post("/{action:follow|unfollow}", golferActionPOST)
			r.Get("/cheevos", golferCheevosGET)
			r.Get("/holes", redir("holes/rankings/lang/bytes"))
			r.Get("/holes/{scoring:bytes|chars}", redir("rankings/lang/{scoring}"))
			r.Get("/holes/{display:rankings|points}/{scope:lang|overall}/{scoring:bytes|chars}", golferHolesGET)
			r.Get("/{hole}/{lang}/{scoring}", golferSolutionGET)
			// r.Post("/{hole}/{lang}/{scoring}", golferSolutionPOST)
		})
		r.Get("/ideas", ideasGET)
		r.Route("/rankings", func(r chi.Router) {
			// Redirect some old URLs that got out.
			r.Get("/", redir("/rankings/holes/all/all/bytes"))
			r.Get("/cheevos", redir("/rankings/cheevos/all"))
			r.Get("/holes", redir("/rankings/holes/all/all/bytes"))
			r.Get("/holes/all/all/all", redir("/rankings/holes/all/all/bytes"))
			r.Get("/langs/bytes", redir("/rankings/langs/all/bytes"))
			r.Get("/langs/chars", redir("/rankings/langs/all/chars"))
			r.Get("/medals", redir("/rankings/medals/all/all/all"))
			r.Get("/solutions", redir("/rankings/misc/solutions"))

			r.Get("/holes/{hole}/{lang}/{scoring}", rankingsHolesGET)
			r.Get("/recent-holes/{lang}/{scoring}", rankingsHolesGET)

			r.Get("/cheevos/{cheevo}", rankingsCheevosGET)
			r.Get("/medals/{hole}/{lang}/{scoring}", rankingsMedalsGET)
			r.Get("/langs/{lang}/{scoring}", rankingsLangsGET)
			r.Get("/misc/{type}", rankingsMiscGET)
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
		r.Post("/solution", solutionPOST)
		r.Get("/stats", statsGET)
		r.Get("/stats/{page:countries}", statsCountriesGET)
		r.Get("/stats/{page:golfers}", statsGolfersGET)
		r.Get("/stats/{page:holes|langs}", statsTableGET)
		r.Get("/stats/{page:unsolved-holes}", statsUnsolvedHolesGET)
		r.Get("/wiki", wikiGET)
		r.Get("/wiki/*", wikiGET)
	})

	return r
}
