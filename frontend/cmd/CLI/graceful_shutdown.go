package main

// import (
// 	"context"
// 	"fmt"
// 	"log"
// 	"os"
// 	"os/signal"
// 	"sync"
// 	"syscall"
// 	"time"
// )

// // ShutdownHandler - обработчик graceful shutdown приложения
// type ShutdownHandler struct {
// 	mu            sync.RWMutex
// 	shutdownFuncs []ShutdownFunc
// 	isShutting    bool
// 	shutdownChan  chan struct{}
// }

// // ShutdownFunc - функция, выполняемая при shutdown
// type ShutdownFunc func(context.Context) error

// // NewShutdownHandler - создать обработчик graceful shutdown приложения
// func NewShutdownHandler() *ShutdownHandler {
// 	return &ShutdownHandler{
// 		shutdownFuncs: make([]ShutdownFunc, 0),
// 		shutdownChan:  make(chan struct{}, 1),
// 	}
// }

// // Register регистрирует функцию для выполнения при shutdown
// func (h *ShutdownHandler) Register(f ShutdownFunc) {
// 	h.mu.Lock()
// 	defer h.mu.Unlock()
// 	h.shutdownFuncs = append(h.shutdownFuncs, f)
// }

// // WaitForShutdown ждет сигналов завершения
// func (h *ShutdownHandler) WaitForShutdown(ctx context.Context) {
// 	sigChan := make(chan os.Signal, 1)
// 	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

// 	select {
// 	case <-sigChan:
// 		fmt.Println("\n📡 Received shutdown signal")
// 		h.InitiateShutdown(ctx)
// 	case <-ctx.Done():
// 		fmt.Println("\n📡 Context cancelled, initiating shutdown")
// 		h.InitiateShutdown(ctx)
// 	case <-h.shutdownChan:
// 		fmt.Println("\n📡 Internal shutdown requested")
// 		h.InitiateShutdown(ctx)
// 	}
// }

// // InitiateShutdown инициирует процесс graceful shutdown
// func (h *ShutdownHandler) InitiateShutdown(ctx context.Context) {
// 	h.mu.Lock()
// 	if h.isShutting {
// 		h.mu.Unlock()
// 		return
// 	}
// 	h.isShutting = true
// 	h.mu.Unlock()

// 	fmt.Println("🔄 Starting graceful shutdown...")

// 	// Создаем контекст с таймаутом для shutdown
// 	shutdownCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
// 	defer cancel()

// 	// Выполняем все зарегистрированные функции в обратном порядке (LIFO)
// 	var wg sync.WaitGroup
// 	h.mu.RLock()
// 	funcs := make([]ShutdownFunc, len(h.shutdownFuncs))
// 	copy(funcs, h.shutdownFuncs)
// 	h.mu.RUnlock()

// 	// Выполняем функции в обратном порядке (стек)
// 	for i := len(funcs) - 1; i >= 0; i-- {
// 		wg.Add(1)
// 		go func(f ShutdownFunc) {
// 			defer wg.Done()
// 			if err := f(shutdownCtx); err != nil {
// 				log.Printf("⚠️  Shutdown function failed: %v", err)
// 			}
// 		}(funcs[i])
// 	}

// 	// Ждем завершения всех функций или таймаута
// 	done := make(chan struct{})
// 	go func() {
// 		wg.Wait()
// 		close(done)
// 	}()

// 	select {
// 	case <-done:
// 		fmt.Println("✅ Graceful shutdown completed successfully")
// 	case <-shutdownCtx.Done():
// 		fmt.Println("⏰ Shutdown timeout reached, forcing exit")
// 	}

// 	// Гарантированный выход
// 	os.Exit(0)
// }

// // RequestShutdown запрашивает shutdown изнутри приложения
// func (h *ShutdownHandler) RequestShutdown() {
// 	select {
// 	case h.shutdownChan <- struct{}{}:
// 	default:
// 		// Уже запрошен shutdown
// 	}
// }

// // IsShutting проверяет, идет ли процесс shutdown
// func (h *ShutdownHandler) IsShutting() bool {
// 	h.mu.RLock()
// 	defer h.mu.RUnlock()
// 	return h.isShutting
// }
