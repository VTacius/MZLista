BASE=$(PWD)
BIN=$(BASE)/bin

ZPL:=$(file < private/zimbra_password_ldap) 
ZUL:=$(file < private/zimbra_usuario_ldap) 
ZSL:=$(file < private/zimbra_servidor_ldap) 

LDFLAGS:=-ldflags 
LDFLAGS+="-s -w
LDFLAGS+=-X xibalba.com/vtacius/MZLista/utils.PasswordLdap=${ZPL}
LDFLAGS+=-X xibalba.com/vtacius/MZLista/utils.UsuarioLdap=${ZUL}
LDFLAGS+=-X xibalba.com/vtacius/MZLista/utils.ServidorLdap=${ZSL}"
.PHONY: dep build 

all: build 

dep:
	go get 

build: dep 
	go build -o ./bin/mzlista ${LDFLAGS} 
