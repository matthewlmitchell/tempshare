package main

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// initializeServer defines the necessary settings for TLS in the tls.Config struct, configures
// the http.Server struct, then starts the server with support for graceful shutdown
func (app *application) initializeServer() error {

	tlsConfig := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		},
	}

	srv := &http.Server{
		Addr:     fmt.Sprintf(":%d", app.serverConfig.port),
		ErrorLog: app.errorLog,
		Handler: app.routes(),
		TLSConfig:      tlsConfig,
		TLSNextProto:   make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 520192, // 0.5MB minus 4096 bytes that Go adds on top automatically
	}

	// This channel will be used to receive any errors from inside the shutdown goroutine
	shutdownError := make(chan error)

	// Listen in the background for any signals to for server shutdown
	go func() {
		// Buffering the channel with size 1 prevents the channel from possibly missing
		// a signal if it is sent before the channel has allocated space to receive it.
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGKILL)

		// Block on the channel until a signal is received, then store it in a variable
		sig := <-quit

		app.infoLog.Println("Shutting down server with signal", sig.String())

		// This context will timeout after 20 seconds have passed
		ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
		defer cancel()

		// Attempt to gracefully shutdown the server within the 20 second timeout context
		// srv.Shutdown() will return any errors if necessary, which will be sent into our 
		// shutdownError channel outside of the goroutine
		shutdownError <- srv.Shutdown(ctx)

		// The goroutine will then exit successfully (status code 0)
	}()

	app.infoLog.Printf("Starting server at :%d in %s mode", app.serverConfig.port, app.serverConfig.env)

	// Start the server and look for any errors. If the error is not related to the server being shutdown, return it
	if err := srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem"); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Block on the error channel until a value is received, if it is non-nil return it
	err := <-shutdownError
	if err != nil {
		return err
	}

	// Since the value received from err was nil (if it was non-nil we wouldn't reach this line),
	// print to the console that the server shutdown was successful.
	app.infoLog.Printf("Stopped server at %s\n", srv.Addr)

	return nil
}
