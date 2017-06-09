FLAGS=-ldflags -s

TARGETS=install build

all:	run

run:	install
	render

$(TARGETS):
	go $@ $(FLAGS)
