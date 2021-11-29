/* Package declaration */

package main


/* Import(s) */

import (
  "context"
  "fmt"
  "net/http"
  "os"
  "github.com/gin-gonic/gin"
  "github.com/jackc/pgx/v4"
)


/* Constant(s) */

const ANY_IPV4_ADDRESS = "0.0.0.0"
const BIND_ADDRESS_TEMPLATE = "%s:%s"
const DATABASE_URL_KEY = "DATABASE_URL"
const DEFAULT_PORT = "5001"
const ENV_OR_PANIC_MESSAGE_TEMPLATE = `Key: "%s" was not found in the environment`;
const LOCALHOST = "localhost"
const PORT_KEY = "PORT"
const SERVER_PUBLIC_ADDRESS_KEY = "SERVER_PUBLIC_ADDRESS"
const SESSION_COOKIE = "_todo"
const DATE_FORMAT = "2006-01-02"
const INDEX_HTML_TEMPLATE = "index.html.tmpl"


/* Helper(s) */

func bindAddress() string {
  return fmt.Sprintf(
    BIND_ADDRESS_TEMPLATE,
    ANY_IPV4_ADDRESS,
    bindPort(),
  )
}

func bindInterface() string {
  if gin.Mode() == gin.ReleaseMode {
    return ANY_IPV4_ADDRESS
  }
  return LOCALHOST
}

func bindPort() string {
  return envOrDefault(PORT_KEY, DEFAULT_PORT)
}

func databaseConnection() *pgx.Conn {
  databaseConnection, databaseConnectionError := pgx.Connect(context.Background(), databaseUrl())
  if databaseConnectionError != nil {
    panic(databaseConnectionError)
  }
  return databaseConnection
}

func databaseUrl() string {
  return envOrPanic(DATABASE_URL_KEY)
}

func envOrDefault(key string, defaultValue string) string {
  result, found := os.LookupEnv(key)
  if found {
    return result
  }
  return defaultValue
}

func envOrPanic(key string) string {
  result, found := os.LookupEnv(key)
  if !found {
    panic(fmt.Sprintf(ENV_OR_PANIC_MESSAGE_TEMPLATE, key))
  }
  return result
}


/* Handler(s) */

func indexHandler(ginContext *gin.Context) {
  ginContext.HTML(
    http.StatusOK,
    INDEX_HTML_TEMPLATE,
    nil,
  )
}

/* Application entry-point */

func main() {
  router := gin.New()
  router.LoadHTMLGlob("templates/*.tmpl")
  router.Use(gin.Logger(), gin.Recovery())
  router.GET("/", indexHandler)
  router.StaticFile("/favicon.ico", "assets/favicon.ico")
  router.Static("/assets", "assets")
  router.Run(bindAddress())
}
