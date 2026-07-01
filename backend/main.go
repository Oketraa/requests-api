// @title Requests API
// @version 1.0
// @description Backend для управления заявками
// @host localhost
// @BasePath /
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"study/database"
	_ "study/docs"
	"study/handlers"
	"study/metrics"
	"study/middleware"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger"
)

const JSON_Log = true

func main() {
	if JSON_Log {
		// Р’РєР»СЋС‡Р°РµРј РіР»РѕР±Р°Р»СЊРЅС‹Р№ JSON-СЂРµР¶РёРј РґР»СЏ slog (РџСѓРЅРєС‚ 13.6)
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
		slog.SetDefault(logger)
	} else {
		// РќР°СЃС‚СЂРѕР№РєРё РґРµС„РѕР»С‚РЅРѕРіРѕ С‚РµРєСЃС‚РѕРІРѕРіРѕ Р»РѕРіРіРµСЂР° РґР»СЏ 13.4 (РµСЃР»Рё РЅСѓР¶РЅС‹ РєР°СЃС‚РѕРјРЅС‹Рµ РїСЂРµС„РёРєСЃС‹)
		log.SetFlags(log.LstdFlags)
	}
	// 1. Р�РЅРёС†РёР°Р»РёР·РёСЂСѓРµРј Р±Р°Р·Сѓ РґР°РЅРЅС‹С…
	err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}

	// 2. РЎРѕР·РґР°РµРј РѕРґРёРЅ РёР·РѕР»РёСЂРѕРІР°РЅРЅС‹Р№ СЂРѕСѓС‚РµСЂ (mux) РґР»СЏ РІСЃРµРіРѕ РїСЂРёР»РѕР¶РµРЅРёСЏ
	mux := http.NewServeMux()

	// 3. Р РµРіРёСЃС‚СЂРёСЂСѓРµРј СЃР»СѓР¶РµР±РЅС‹Рµ СЂСѓС‡РєРё Р±СЌРєРµРЅРґРµСЂР° (Healthchecks)
	mux.HandleFunc("/health", handlers.Health)
	mux.HandleFunc("/live", handlers.Live)
	mux.HandleFunc("/ready", handlers.Ready)

	// 4. Р РµРіРёСЃС‚СЂРёСЂСѓРµРј СЂСѓС‡РєРё Р±РёР·РЅРµСЃ-Р»РѕРіРёРєРё (API Р·Р°СЏРІРѕРє)
	mux.HandleFunc("/api/requests", handlers.GetRequests)
	mux.HandleFunc("/api/requests/", handlers.RequestByID)
	mux.HandleFunc("/api/requests/create", handlers.CreateRequest)

	// 5. Р РµРіРёСЃС‚СЂРёСЂСѓРµРј С‚РІРѕСЋ СЂСѓС‡РєСѓ СЃР±РѕСЂР° РјРµС‚СЂРёРє Prometheus
	mux.Handle("/metrics", promhttp.Handler())

	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	// 6. РќР°РєР°С‚С‹РІР°РµРј РёРЅС„СЂР°СЃС‚СЂСѓРєС‚СѓСЂРЅС‹Рµ РјРёРґР»РІР°СЂРё Р±СЌРєРµРЅРґРµСЂР° РЅР° СЂРѕСѓС‚РµСЂ
	// Р¦РµРїРѕС‡РєР° РёРґРµС‚ СЃРЅРёР·Сѓ РІРІРµСЂС…: Cors -> Recovery -> Logger -> RequestID -> Р РѕСѓС‚РµСЂ
	handler := middleware.Logger(JSON_Log)(mux)
	handler = middleware.RequestID(handler)
	handler = middleware.Recovery(handler)
	handler = middleware.Cors(handler)
	handler = middleware.UserAgent(handler)

	// 7. РћР±РѕСЂР°С‡РёРІР°РµРј РїРѕР»СѓС‡РёРІС€РёР№СЃСЏ РїРёСЂРѕРі РІ С‚РІРѕСЋ РјРµС‚СЂРёС‡РµСЃРєСѓСЋ РјРёРґР»РІР°СЂСЊ,
	// С‡С‚РѕР±С‹ Prometheus СЃС‡РёС‚Р°Р» РґР»РёС‚РµР»СЊРЅРѕСЃС‚СЊ Рё СЃС‚Р°С‚СѓСЃС‹ РѕС‚РІРµС‚РѕРІ СЃ СѓС‡РµС‚РѕРј РІСЃРµС… РїСЂР°РІРёР»
	finalHandler := metrics.MetricsMiddleware(handler)

	// 8. Р—Р°РїСѓСЃРєР°РµРј РћР”Р�Рќ РµРґРёРЅСЃС‚РІРµРЅРЅС‹Р№ СЃРµСЂРІРµСЂ РЅР° РїРѕСЂС‚Сѓ 8080
	server := &http.Server{
		Addr:    ":8080",
		Handler: finalHandler,
	}

	go func() {
		fmt.Println("Server started on port 8080")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop, os.Interrupt, syscall.SIGTERM,
	)

	<-stop

	fmt.Println("Shutting down server..")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	err = server.Shutdown(ctx)

	if err != nil {
		log.Println("Shutdown error:", err)
	}

	err = database.DB.Close()

	if err != nil {
		log.Println("Database close error:", err)
	}

	fmt.Println("Server stopped")
}
