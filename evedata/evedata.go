package evedata

import (
	"encoding/gob"
	"os"
	"strings"

	"log"
	"net/http"

	"github.com/antihax/evedata/appContext"
	"github.com/antihax/evedata/config"
	"github.com/antihax/evedata/discord"
	"github.com/antihax/evedata/internal/apicache"
	"github.com/antihax/evedata/internal/redisqueue"
	"github.com/antihax/evedata/internal/tokenstore"

	"github.com/antihax/evedata/models"
	"github.com/antihax/goesi"

	gsr "github.com/antihax/redistore"
	"github.com/gorilla/context"
	"golang.org/x/oauth2"
)

var ctx appContext.AppContext

// GetContext Returns the global appContext for EVEData Server
func GetContext() *appContext.AppContext {
	return &ctx
}

// GoServer Runs the EVEData Server
func GoServer() {
	var err error

	// Read configuation.
	if ctx.Conf, err = config.ReadConfig("config/config.conf"); err != nil {
		log.Fatalf("Error reading configuration: %v", err)
	}

	// Build the redis pool
	ctx.Cache = setupRedis(GetContext())

	ctx.OutQueue = redisqueue.NewRedisQueue(
		ctx.Cache,
		"evedata-hammer",
	)

	// Build a HTTP Client pool this client will be shared with APIs for:
	//   - ESI
	//   - ZKillboard
	//   - EVE SSO
	//   - EVE CREST and XML
	ctx.HTTPClient = apicache.CreateHTTPClientCache(ctx.Cache)
	if ctx.HTTPClient == nil {
		panic("http client is null")
	}

	// Build Connection Pool
	if ctx.Db, err = models.SetupDatabase(ctx.Conf.Database.Driver, ctx.Conf.Database.Spec); err != nil {
		log.Fatalf("Cannot build database pool: %v", err)
	}

	// Setup the SSO authenticator, this is the main login.
	ssoScopes := []string{}
	ctx.SSOAuthenticator = goesi.NewSSOAuthenticator(
		ctx.HTTPClient,
		ctx.Conf.CREST.SSO.ClientID,
		ctx.Conf.CREST.SSO.SecretKey,
		ctx.Conf.CREST.SSO.RedirectURL,
		ssoScopes)

	// Setup the Token authenticator, this handles sub characters.
	tokenScopes := models.GetCharacterScopes()
	ctx.TokenAuthenticator = goesi.NewSSOAuthenticator(
		ctx.HTTPClient,
		ctx.Conf.CREST.Token.ClientID,
		ctx.Conf.CREST.Token.SecretKey,
		ctx.Conf.CREST.Token.RedirectURL,
		tokenScopes)

	bootstrapScopes := strings.Split("esi-calendar.respond_calendar_events.v1 esi-calendar.read_calendar_events.v1 esi-mail.organize_mail.v1 esi-mail.read_mail.v1 esi-mail.send_mail.v1 esi-wallet.read_character_wallet.v1 esi-wallet.read_corporation_wallet.v1 esi-search.search_structures.v1 esi-universe.read_structures.v1 esi-corporations.read_corporation_membership.v1 esi-markets.structure_markets.v1 esi-characters.read_chat_channels.v1 esi-corporations.track_members.v1 esi-wallet.read_corporation_wallets.v1 esi-corporations.read_divisions.v1 esi-assets.read_corporation_assets.v1", " ")

	// Setup the bootstrap authenticator. Needed to update the site main token.
	ctx.ESIBootstrapAuthenticator = goesi.NewSSOAuthenticator(
		ctx.HTTPClient,
		ctx.Conf.CREST.ESIAccessToken.ClientID,
		ctx.Conf.CREST.ESIAccessToken.SecretKey,
		ctx.Conf.CREST.ESIAccessToken.RedirectURL,
		bootstrapScopes)

	// Get the token from config and build a TokenSource (refreshes the token if needed).
	token := &oauth2.Token{
		AccessToken:  ctx.Conf.CREST.ESIAccessToken.AccessToken,
		TokenType:    ctx.Conf.CREST.ESIAccessToken.TokenType,
		RefreshToken: ctx.Conf.CREST.ESIAccessToken.RefreshToken,
		Expiry:       ctx.Conf.CREST.ESIAccessToken.Expiry,
	}

	ctx.ESIPublicToken, err = ctx.ESIBootstrapAuthenticator.TokenSource(token)
	if err != nil {
		log.Fatalf("Error starting bootstrap ESI client: %v", err)
	}

	// Setup our token store for oauth2 optimizations
	ctx.TokenStore = tokenstore.NewTokenStore(ctx.Cache, ctx.Db, ctx.TokenAuthenticator)

	// Setup the EVE ESI Client
	ctx.ESI = goesi.NewAPIClient(ctx.HTTPClient, ctx.Conf.UserAgent)

	// Create a redis session store.
	ctx.Store, err = gsr.NewRediStoreWithPool(ctx.Cache, []byte(ctx.Conf.Store.Key))
	defer ctx.Store.Close()
	if err != nil {
		log.Fatalf("Cannot build database pool: %v", err)
	}

	// Set options for the store
	ctx.Store.SetMaxLength(1024 * 10)
	ctx.Store.Options.Domain = ctx.Conf.Store.Domain

	// Register structs for storage.
	gob.Register(oauth2.Token{})
	gob.Register(goesi.VerifyResponse{})

	// Set our logging flags.
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Run the discord bot
	if ctx.Conf.Discord.Enabled {
		go discord.GoDiscordBot(&ctx)
	}

	// Handle command line arguments
	if len(os.Args) > 1 {

		if os.Args[1] == "dumpdb" {
			// Dump the database to sql file.
			log.Printf("Dumping Database to ./sql/evedata.sql\n")
			err := models.DumpDatabase("./sql/evedata.sql", "evedata")
			if err != nil {
				log.Fatalln(err)
			}

		} else if os.Args[1] == "flushredis" {
			// Erase everything in redis for modified deployments
			log.Printf("Flushing Redis\n")
			r := ctx.Cache.Get()
			r.Do("FLUSHALL")
			r.Close()
		}
	}

	// Allocate the routes.
	rtr := NewRouter(&ctx)

	log.Printf("EveData Listening port 3000...\n")
	http.ListenAndServe(":3000", context.ClearHandler(rtr))
	log.Printf("EveData Quitting..\n")
}
