package bootstrap

func Root(flag string) {
	Register(flag)
	Migrate()
	Setup()
	r.Logger.Fatal(r.Start(":8443"))
}
