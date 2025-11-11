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
	r.Get("/callback", callbackGET)
	r.Get("/callback/dev", callbackDevGET)
	r.Get("/feeds/{feed:(?:atom|json|rss)}", feedGET)
	r.Get("/golfers/{name}/avatar", golferAvatarGET)
	r.Get(`/golfers/{name}/avatar/{size:\d+}`, golferAvatarGET)
	r.Get("/healthz", healthzGET)
	r.Post("/log-out", logOutPOST)
	r.Get("/random", randomGET)
	r.Get("/sitemap.xml", sitemapGET)

	// Legacy redirects.
	r.Get("/rankings", redir("/rankings/holes/all/all/bytes"))
	r.Get("/rankings/cheevos", redir("/rankings/cheevos/all"))
	r.Get("/rankings/holes", redir("/rankings/holes/all/all/bytes"))
	r.Get("/rankings/holes/all/all/all", redir("/rankings/holes/all/all/bytes"))
	r.Get("/rankings/langs/bytes", redir("/rankings/langs/all/bytes"))
	r.Get("/rankings/langs/chars", redir("/rankings/langs/all/chars"))
	r.Get("/rankings/medals", redir("/rankings/medals/all/all/all"))
	r.Get("/rankings/solutions", redir("/rankings/misc/solutions"))
	r.Get("/recent", redir("/recent/solutions/all/all/bytes"))
	r.Get("/users/{name}", redir("/golfers/{name}"))

	// HTML routes that need middleware.Golfer.
	r.With(middleware.RedirHolesLangs, middleware.Golfer).Group(func(r chi.Router) {
		r.Get("/", homeGET)
		r.Get("/{hole}", holeGET)
		r.Get("/about", aboutGET)
		r.Get("/feeds", feedsGET)
		r.Get("/ideas", ideasGET)
		r.Get("/rankings/cheevos/{cheevo}", rankingsCheevosGET)
		r.Get("/rankings/holes/{hole}/{lang}/{scoring}", rankingsHolesGET)
		r.Get("/rankings/langs/{lang}/{scoring}", rankingsLangsGET)
		r.Get("/rankings/medals/{hole}/{lang}/{scoring}", rankingsMedalsGET)
		r.Get("/rankings/misc/{type}", rankingsMiscGET)
		r.Get("/rankings/recent-holes/{lang}/{scoring}", rankingsHolesGET)
		r.Get("/recent/{lang}", recentGET)
		r.Get("/recent/golfers", recentGolfersGET)
		r.Get("/recent/solutions/{hole}/{lang}/{scoring}", recentSolutionsGET)
		r.Get("/scores/{hole}/{lang}", scoresGET)
		r.Get("/scores/{hole}/{lang}/all", scoresAllGET)
		r.Get("/scores/{hole}/{lang}/{scoring}", scoresGET)
		r.Get("/scores/{hole}/{lang}/{scoring}/{page}", scoresGET)
		r.Post("/solution", solutionPOST)
		r.Get("/stats", statsGET)
		r.Get("/stats/{page:cheevos}", statsCheevosGET)
		r.Get("/stats/{page:countries}", statsCountriesGET)
		r.Get("/stats/{page:golfers}", statsGolfersGET)
		r.Get("/stats/{page:(?:holes|langs)}", statsTableGET)
		r.Get("/stats/{page:unsolved-holes}", statsUnsolvedHolesGET)
		r.Get("/wiki", wikiGET)
		r.Get("/wiki/*", wikiGET)

		r.With(middleware.AdminArea).Route("/admin", func(r chi.Router) {
			r.Get("/", adminGET)
			r.Get("/banners", adminBannersGET)
			r.Post("/banners/{banner}", adminBannerPOST)
			r.Get("/solutions", adminSolutionsGET)
			r.Get("/solutions/run", adminSolutionsRunGET)
			r.Get("/solutions/{hole}/{lang}/{golferID}", adminSolutionGET)
		})

		r.With(middleware.API).Route("/api", func(r chi.Router) {
			r.Get("/", apiGET)
			r.Get("/cheevos", apiCheevosGET)
			r.Get("/cheevos/{cheevo}", apiCheevoGET)
			r.Get("/golfers/{golfer}", apiGolferGET)
			r.Get("/holes", apiHolesGET)
			r.Get("/langs", apiLangsGET)
			r.Get("/panic", apiPanicGET)
			r.Get("/solutions-log", apiSolutionsLogGET)
			r.Get("/suggestions/golfers", apiSuggestionsGolfersGET)
			r.Get("/wiki/*", apiWikiPageGET)

			// API routes that use {hole} or {lang}.
			r.With(middleware.RedirHolesLangs).Group(func(r chi.Router) {
				r.Get("/holes/{hole}", apiHoleGET)
				r.Get("/langs/{lang}", apiLangGET)
				r.Get("/mini-rankings/{hole}/{lang}/{scoring:(?:bytes|chars)}"+
					"/{view:(?:top|me|following)}", apiMiniRankingsGET)

				// API routes that require a logged-in golfer.
				r.With(middleware.GolferArea).Group(func(r chi.Router) {
					r.Get("/notes", apiNotesGET)
					r.Delete("/notes/{hole}/{lang}", apiNoteDELETE)
					r.Get("/notes/{hole}/{lang}", apiNoteGET)
					r.Put("/notes/{hole}/{lang}", apiNotePUT)
					r.Get("/solutions-search", apiSolutionsSearchGET)
				})
			})
		})

		r.With(middleware.GolferArea).Route("/golfer", func(r chi.Router) {
			r.Post("/cancel-delete", golferCancelDeletePOST)
			r.Get("/connect/{connection}", golferConnectGET)
			r.Post("/delete", golferDeletePOST)
			r.Post("/delete-solution", golferDeleteSolutionPOST)
			r.Get("/disconnect/{connection}", golferDisconnectGET)
			r.Get("/export", golferExportGET)
			r.Post("/hide-banner", golferHideBannerPOST)
			r.Get("/code-search", golferSearchGET)
			r.Get("/settings", golferSettingsGET)
			r.Post("/settings", golferSettingsPOST)
			r.Get("/settings/{page:delete-account}", golferSettingsDeleteAccountGET)
			r.Get("/settings/{page:export-data}", golferSettingsExportDataGET)
			r.Post("/settings/reset", golferSettingsResetPOST)
			r.Post("/settings/save", golferSettingsSavePOST)
		})

		r.With(middleware.GolferInfo).Route("/golfers/{name}", func(r chi.Router) {
			r.Get("/", golferGET)
			r.Post("/{action:(?:follow|unfollow)}", golferActionPOST)
			r.Get("/cheevos", golferCheevosGET)
			r.Get("/holes", redir("holes/rankings/lang/bytes"))
			r.Get("/holes/{scoring:(?:bytes|chars)}", redir("rankings/lang/{scoring}"))
			r.Get("/holes/{display:(?:rankings|points)}/{scope:(?:lang|overall)}"+
				"/{scoring:(?:bytes|chars)}", golferHolesGET)

			// Golfer info routes that use {hole} or {lang}.
			r.With(middleware.RedirHolesLangs).Group(func(r chi.Router) {
				r.Get("/{hole}/{lang}/{scoring}", golferSolutionGET)
				r.Post("/{hole}/{lang}/{scoring}", golferSolutionPOST)
			})
		})
	})

	return r
}
