
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUFFIXES = .go .c
OBJEXT= none

EXTRA_PROGRAMS = builder
builder_SOURCES = builder.go

builder$(EXEEXT): $(builder_SOURCES)
	$(GO) build -o builder$(EXEEXT) $(builder_SOURCES)

bin_PROGRAMS = s2cli
sbin_PROGRAMS = s2srv s2pwd

s2srv_SOURCES = s2srv.go
s2pwd_SOURCES = s2pwd.go
s2cli_SOURCES = s2cli.go

GOFLAGS = -ldflags="-s -w"

EXTRA_s2pwd_SOURCES = \
	config/config.go \
	server/user-model/user_model.go

EXTRA_s2cli_SOURCES = \
	config/config.go \
	client/client.go \
	server/user-model/user_model.go

EXTRA_s2srv_SOURCES = \
	config/config.go \
	daemon/daemon.go \
	server/server.go \
	server/bucket-controller/bucket_controller.go \
	server/file-controller/file_controller.go \
	server/user-controller/user_controller.go \
	server/status-controller/status_controller.go \
	server/user-model/user_model.go \
	tools/file.go

EXTRA_s2srv_SOURCES += \
	bundle/public.go

bundle/public.go: builder$(EXEEXT) public/index.html
	./builder$(EXEEXT) --package=bundle -o bundle/public.go public/

s2srv$(EXEEXT): $(s2srv_SOURCES) $(EXTRA_s2srv_SOURCES)
	$(GO) build $(GOFLAGS) -o s2srv$(EXEEXT) $(s2srv_SOURCES)

s2cli$(EXEEXT): $(s2cli_SOURCES) $(EXTRA_s2cli_SOURCES)
	$(GO) build $(GOFLAGS) -o s2cli$(EXEEXT) $(s2cli_SOURCES)

s2pwd$(EXEEXT): $(s2pwd_SOURCES) $(EXTRA_s2pwd_SOURCES)
	$(GO) build $(GOFLAGS) -o s2pwd$(EXEEXT) $(s2pwd_SOURCES)

EXTRA_DIST = \
	README.md \
	go.mod \
	go.sum \
	\
	front/public \
	front/src \
	public/ \
	\
	front/.gitignore \
	front/babel.config.js \
	front/package-lock.json \
	front/package.json \
	front/postcss.config.js \
	front/README.md \
	front/webpack.config.devel.js \
	front/webpack.config.prod.js


clean-local:
	rm -f builder$(EXEEXT)
	rm -rf autom4te.cache

install-data-local:
	test -z $(DESTDIR)$(APP_LIBDIR) || $(MKDIR_P) $(DESTDIR)$(APP_LIBDIR)
#	cp -r public $(DESTDIR)$(APP_LIBDIR)
	test -z $(DESTDIR)$(APP_CONFDIR) || $(MKDIR_P) $(DESTDIR)$(APP_CONFDIR)
	test -z $(DESTDIR)$(APP_LOGDIR) || $(MKDIR_P) $(DESTDIR)$(APP_LOGDIR)
	test -z $(DESTDIR)$(APP_RUNDIR) || $(MKDIR_P) $(DESTDIR)$(APP_RUNDIR)
	test -z $(DESTDIR)$(APP_DATABASEDIR) || $(MKDIR_P) $(DESTDIR)$(APP_DATABASEDIR)
#EOF
