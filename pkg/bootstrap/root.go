package bootstrap

func Root() {
	Register()
	Migrate()
	Setup()
	r.Logger.Fatal(r.Start(":8443"))
}
