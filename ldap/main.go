package main

import (
	"crypto/tls"
	"fmt"
	"github.com/go-ldap/ldap/v3"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) <= 1 {
		log.Printf("Usage: %s <config-file>", os.Args[0])
		os.Exit(1)
	}

	configFile := os.Args[1]
	viper.SetConfigFile(configFile)
	contents, err := os.ReadFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("read config file fail: %s", err.Error()))
	}

	//Replace environment variables
	err = viper.ReadConfig(strings.NewReader(os.ExpandEnv(string(contents))))
	if err != nil {
		panic(fmt.Sprintf("parse config file fail: %s", err.Error()))
	}

	var ldapConfig = viper.Sub("ldap")
	if ldapConfig == nil {
		panic("config not found ldap")
	}

	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", viper.GetString("ldap.host"), viper.GetInt("ldap.port")))
	if err != nil {
		panic(err)
	}

	if viper.GetBool("ldap.tls") {
		err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			panic(err)
		}
	}

	conn.SetTimeout(5 * time.Second)
	defer conn.Close()

	err = conn.Bind(viper.GetString("ldap.bindUserDn"), viper.GetString("ldap.bindPassword"))
    if err != nil {
        panic(err)
    }

	searchRequest := ldap.NewSearchRequest(
		viper.GetString("ldap.baseDn"),
        ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false, 
        fmt.Sprintf("(&(objectClass=organizationalPerson)(uid=%s))", viper.GetString("ldap.username")),
        []string{"dn"},
        nil,
    )

	sr, err := conn.Search(searchRequest)
    if err != nil {
        panic(err)
    }

	if len(sr.Entries) != 1 {
        panic("User does not exist or too many entries returned")
    }

	userDn := sr.Entries[0].DN

	err = conn.Bind(userDn, viper.GetString("ldap.password"))
    if err != nil {
		panic(err)
    }

	log.Printf("user authentication success")
}
