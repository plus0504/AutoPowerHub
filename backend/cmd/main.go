package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"autopowerhub/api/handler"
	"autopowerhub/config"
	"autopowerhub/database"
	"autopowerhub/repository"
	"autopowerhub/router"
	authsvc "autopowerhub/service/auth"
	blesvc "autopowerhub/service/ble"
	devicesvc "autopowerhub/service/device"
)

func main() {
	cfgPath := flag.String("config", resolveConfig(), "path to config.yaml")
	flag.Parse()

	// baseDir is the directory that contains config.yaml.
	// All relative paths (web/, data.db) are resolved from here,
	// so the binary works regardless of which subdirectory it is launched from.
	baseDir, err := filepath.Abs(filepath.Dir(*cfgPath))
	if err != nil {
		log.Fatalf("resolve base dir: %v", err)
	}

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	// Resolve sqlite path relative to baseDir if it is not absolute.
	if !filepath.IsAbs(cfg.SQLite.Path) {
		cfg.SQLite.Path = filepath.Join(baseDir, cfg.SQLite.Path)
	}

	db, err := database.Init(cfg)
	if err != nil {
		log.Fatalf("init database: %v", err)
	}

	bleMgr, err := blesvc.NewManager()
	if err != nil {
		log.Fatalf("init BLE: %v", err)
	}
	defer bleMgr.Close()

	userRepo := repository.NewUserRepository(db)
	deviceRepo := repository.NewDeviceRepository(db)
	logRepo := repository.NewLogRepository(db)

	authService := authsvc.NewService(userRepo, cfg)
	deviceService := devicesvc.NewService(deviceRepo, logRepo, bleMgr)

	authHandler := handler.NewAuthHandler(authService)
	deviceHandler := handler.NewDeviceHandler(deviceService)
	debugHandler := handler.NewDebugHandler(deviceRepo, bleMgr)

	r := router.Setup(authHandler, deviceHandler, debugHandler, authService, baseDir)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("AutoPowerHub starting on %s (base: %s)", addr, baseDir)
	if err := r.Run(addr); err != nil {
		log.Fatalf("run server: %v", err)
	}
}

// resolveConfig searches for config.yaml in the current directory first,
// then one level up (handles IDE debug runs launched from backend/cmd/).
func resolveConfig() string {
	for _, p := range []string{"config.yaml", "../config.yaml"} {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}
	return "config.yaml"
}
