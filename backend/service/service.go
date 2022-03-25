package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ontology-tech/ontlogin-sdk-go/did"
	"github.com/ontology-tech/ontlogin-sdk-go/did/ont"
	"github.com/ontology-tech/ontlogin-sdk-go/modules"
	ontloginsdk "github.com/ontology-tech/ontlogin-sdk-go/sdk"

	"github.com/google/uuid"

	"ontlogin-sample/auth"
	"ontlogin-sample/jwt"
)

var loginsdk *ontloginsdk.OntLoginSdk
var mapstore map[string]int

func InitService() {
	mapstore = make(map[string]int)

	vcfilters := make(map[int][]*modules.VCFilter)
	vcfilters[modules.ACTION_AUTHORIZATION] = []*modules.VCFilter{
		{Type: "IdentityCredential", Required: false, TrustRoots: []string{"did:ont:testdid"}},
	}
	conf := &ontloginsdk.SDKConfig{
		Chain: []string{"ONT"},
		Alg:   []string{"ES256"},
		ServerInfo: &modules.ServerInfo{
			Name: "testServcer",
			// Icon:               "http://somepic.jpg",
			Url:                "https://ont.io",
			Did:                "did:ont:sampletest",
			VerificationMethod: "",
		},
		VCFilters: vcfilters,
	}

	resolvers := make(map[string]did.DidProcessor)
	ontresolver, err := ont.NewOntProcessor(false, "http://dappnode3.ont.io:20336", "", "", "")
	if err != nil {
		panic(err)
	}
	resolvers["ont"] = ontresolver
	loginsdk, err = ontloginsdk.NewOntLoginSdk(conf, resolvers, GenUUID, CheckNonce)
	if err != nil {
		panic(err)
	}
}

func RequestChallenge(writer http.ResponseWriter, request *http.Request) {
	cr := &modules.ClientHello{}
	writer.Header().Set("Content-Type", "application/json")
	err := json.NewDecoder(request.Body).Decode(&cr)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	serverHello, err := loginsdk.GenerateChallenge(cr)
	if err != nil {
		writer.Write([]byte(err.Error()))
		return
	}

	bts, _ := json.Marshal(serverHello)

	writer.Write(bts)

}

func Login(writer http.ResponseWriter, request *http.Request) {
	lr := &modules.ClientResponse{}
	log.Printf("login func")
	writer.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(request.Body).Decode(&lr)

	if err != nil {
		log.Printf("err: %s", err.Error())
		writer.Write([]byte(err.Error()))
		return
	}

	log.Printf("ValidateClientResponse")
	err = loginsdk.ValidateClientResponse(lr)

	if err != nil {
		log.Printf("err: %s", err.Error())
		writer.Write([]byte(err.Error()))
		return
	}

	log.Printf("GenerateToken")
	s, err := jwt.GenerateToken(lr.Did)
	writer.Write([]byte(s))

}

func AfterLogin(writer http.ResponseWriter, request *http.Request) {
	if err := auth.CheckLogin(request.Context()); err != nil {
		writer.Write([]byte("please login first"))
		return
	}
	writer.Write([]byte("normal business process"))
}

func GenUUID(action int) string {
	uuid, err := uuid.NewUUID()
	if err != nil {
		return ""
	}
	mapstore[uuid.String()] = action
	return uuid.String()
}

func CheckNonce(nonce string) (int, error) {
	action, ok := mapstore[nonce]
	if !ok {
		return -1, fmt.Errorf("no nonce found")
	}
	return action, nil
}
