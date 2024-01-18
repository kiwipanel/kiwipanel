package bootstrap

func Root() {
	Register()
	Setup()
	r.Logger.Fatal(r.Start(":8443"))
}
