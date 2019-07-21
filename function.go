package startdownrec

import (
	cloudkms "cloud.google.com/go/kms/apiv1"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	kmspb "google.golang.org/genproto/googleapis/cloud/kms/v1"
	"io/ioutil"
	"log"
	"net/http"
)

type Function struct {
	DBConn *sql.DB
}

type DbCreds struct {
	Host string `json:"db_host"`
	User string `json:"db_user"`
	Pass string `json:"db_pass"`
}

type Report struct {
	Hostname string `json:"hostname"`
	Status   string `json:"status"`
}

func (f Function) PostReport(r *Report) {
	fmt.Println(r)
}

func (f Function) Exec(w http.ResponseWriter, r *http.Request) {
	var report Report
	decErr := json.NewDecoder(r.Body).Decode(&report)
	if decErr != nil || report.Status == "" || report.Hostname == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, wErr := w.Write([]byte("Invalid JSON body"))
		if wErr != nil {
			log.Fatal("error responding with error")
		}
		return
	}

	fmt.Println("Hostname: " + report.Hostname)
	fmt.Println("Status: " + report.Status)

	_, qErr := f.DBConn.Query("INSERT INTO `reports` (`hostname`, `status`) VALUES (?, ?);", report.Hostname, report.Status)
	if qErr != nil {
		log.Fatal(qErr)
	}

	f.PostReport(&report)

	w.WriteHeader(200)
	_, wErr := w.Write([]byte("Status recorded"))
	if wErr != nil {
		log.Fatal("Error responding with success")
	}
}

// Run is used by google cloud to kick off the function
func Run(w http.ResponseWriter, r *http.Request) {
	encFile, err := ioutil.ReadFile("config/mysql-creds.json.enc")
	if err != nil {
		log.Fatal(err)
	}

	decFile, err := decryptSymmetric("projects/brennon-loveless/locations/global/keyRings/secrets/cryptoKeys/startdownrec-mysql", encFile)
	if err != nil {
		log.Fatal(err)
	}

	var dbCreds DbCreds
	parseCredsErr := json.Unmarshal(decFile, &dbCreds)
	if parseCredsErr != nil {
		log.Fatal(parseCredsErr)
	}

	if dbCreds.Host == "" || dbCreds.User == "" || dbCreds.Pass == "" {
		log.Fatal("invalid credentials for database")
	}

	dbConn, err := sql.Open("mysql", dbCreds.User + ":" + dbCreds.Pass + "@tcp(" + dbCreds.Host + ":3306)/startdownrec")
	if err != nil {
		log.Fatal("Unable to connect to mysql")
	}
	defer dbConn.Close()

	f := Function{
		DBConn: dbConn,
	}

	f.Exec(w, r)
}

// decrypt will decrypt the input ciphertext bytes using the specified symmetric key
// example keyName: "projects/PROJECT_ID/locations/global/keyRings/RING_ID/cryptoKeys/KEY_ID"
func decryptSymmetric(keyName string, ciphertext []byte) ([]byte, error) {
	ctx := context.Background()
	client, err := cloudkms.NewKeyManagementClient(ctx)
	if err != nil {
		return nil, err
	}

	// Build the request.
	req := &kmspb.DecryptRequest{
		Name:       keyName,
		Ciphertext: ciphertext,
	}
	// Call the API.
	resp, err := client.Decrypt(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.Plaintext, nil
}
