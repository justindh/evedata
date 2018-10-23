package tokenserver

/*
var tokenserver *TokenServer

func TestMain(m *testing.M) {
	sql := sqlhelper.NewTestDatabase()
	redis := redigohelper.ConnectRedisTestPool()

	tokenserver = NewTokenServer(redis, sql, "notreal", "reallynotreal")

	go func() {
		if err := tokenserver.Run(); err != nil {
			panic(err)
		}
	}()
	retCode := m.Run()
	time.Sleep(time.Second * 5)
	tokenserver.Close()

	redis.Close()
	sql.Close()

	time.Sleep(time.Second)

	os.Exit(retCode)
}

func TestTokens(t *testing.T) {
	tokenserver.tokenStore.SetToken(555, 555, &oauth2.Token{
		RefreshToken: "fake",
		AccessToken:  "really fake",
		TokenType:    "Bearer",
		Expiry:       time.Now().Add(time.Hour),
	})

	r, err := grpc.Dial("localhost:3002", grpc.WithInsecure(), grpc.WithCodec(&msgpackcodec.MsgPackCodec{}))
	assert.Nil(t, err)
	assert.NotNil(t, r)

	token := oauth2.Token{}
	err = r.Invoke(context.Background(), "/TokenStore/GetToken", &TokenRequest{CharacterID: 555, TokenCharacterID: 555}, &token)
	assert.Nil(t, err)
}
*/
