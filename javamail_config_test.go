package opennms

// Tests for the Javamail configuration methods – /rest/config/javamail.

import (
	"encoding/json"
	"testing"
)

const javamailDefaultConfigJSON = `{
	"defaultReadConfigName": "localhost",
	"defaultSendConfigName": "localhost",
	"defaultEnd2endConfigName": "localhost"
}`

const javamailReadmailJSON = `{
	"name": "localhost", "host": "localhost", "port": 993, "protocol": "imaps"
}`

const javamailReadmailListJSON = `{"readmail": [` + javamailReadmailJSON + `]}`

const javamailSendmailJSON = `{
	"name": "localhost", "host": "localhost", "port": 25, "protocol": "smtp"
}`

const javamailSendmailListJSON = `{"sendmail": [` + javamailSendmailJSON + `]}`

const javamailEnd2EndJSON = `{
	"name": "localhost",
	"readMailConfigName": "localhost",
	"sendMailConfigName": "localhost"
}`

const javamailEnd2EndListJSON = `{"end2end": [` + javamailEnd2EndJSON + `]}`

func TestGetJavamailDefaultConfig(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/javamail", javamailDefaultConfigJSON)
	result, err := c.GetJavamailDefaultConfig(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["defaultReadConfigName"] != "localhost" {
		t.Errorf("defaultReadConfigName = %v", result["defaultReadConfigName"])
	}
}

func TestSetJavamailDefaultConfig(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "POST "+v1Path+"/config/javamail", 204)
	err := c.SetJavamailDefaultConfig(t.Context(),
		map[string]any{"defaultReadConfigName": "remote"})
	if err != nil {
		t.Fatal(err)
	}
}

// Read-mail configs

func TestGetJavamailReadmails(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/javamail/readmails", javamailReadmailListJSON)
	result, err := c.GetJavamailReadmails(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["readmail"].([]any)[0].(map[string]any)["name"] != "localhost" {
		t.Errorf("result = %v", result)
	}
}

func TestGetJavamailReadmail(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/javamail/readmails/localhost",
		javamailReadmailJSON)
	result, err := c.GetJavamailReadmail(t.Context(), "localhost")
	if err != nil {
		t.Fatal(err)
	}
	if result["protocol"] != "imaps" {
		t.Errorf("protocol = %v", result["protocol"])
	}
}

func TestCreateJavamailReadmail(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/config/javamail/readmails", 201)
	err := c.CreateJavamailReadmail(t.Context(), map[string]any{
		"name": "localhost", "host": "localhost",
		"port": 993, "protocol": "imaps",
	})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "localhost" {
		t.Errorf("name = %v", body["name"])
	}
}

func TestUpdateJavamailReadmail(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "PUT "+v1Path+"/config/javamail/readmails/localhost", 204)
	err := c.UpdateJavamailReadmail(t.Context(), "localhost",
		map[string]string{"port": "995"})
	if err != nil {
		t.Fatal(err)
	}
	if req.body != "port=995" {
		t.Errorf("body = %q, want port=995", req.body)
	}
}

func TestDeleteJavamailReadmail(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/config/javamail/readmails/localhost", 204)
	if err := c.DeleteJavamailReadmail(t.Context(), "localhost"); err != nil {
		t.Fatal(err)
	}
}

// Send-mail configs

func TestGetJavamailSendmails(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/javamail/sendmails", javamailSendmailListJSON)
	result, err := c.GetJavamailSendmails(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["sendmail"].([]any)[0].(map[string]any)["name"] != "localhost" {
		t.Errorf("result = %v", result)
	}
}

func TestGetJavamailSendmail(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/javamail/sendmails/localhost",
		javamailSendmailJSON)
	result, err := c.GetJavamailSendmail(t.Context(), "localhost")
	if err != nil {
		t.Fatal(err)
	}
	if result["protocol"] != "smtp" {
		t.Errorf("protocol = %v", result["protocol"])
	}
}

func TestCreateJavamailSendmail(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/config/javamail/sendmails", 201)
	err := c.CreateJavamailSendmail(t.Context(), map[string]any{
		"name": "localhost", "host": "localhost",
		"port": 25, "protocol": "smtp",
	})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "localhost" {
		t.Errorf("name = %v", body["name"])
	}
}

func TestUpdateJavamailSendmail(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "PUT "+v1Path+"/config/javamail/sendmails/localhost", 204)
	err := c.UpdateJavamailSendmail(t.Context(), "localhost",
		map[string]string{"port": "587"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteJavamailSendmail(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/config/javamail/sendmails/localhost", 204)
	if err := c.DeleteJavamailSendmail(t.Context(), "localhost"); err != nil {
		t.Fatal(err)
	}
}

// End-to-end configs

func TestGetJavamailEnd2ends(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/javamail/end2ends", javamailEnd2EndListJSON)
	result, err := c.GetJavamailEnd2ends(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if result["end2end"].([]any)[0].(map[string]any)["name"] != "localhost" {
		t.Errorf("result = %v", result)
	}
}

func TestGetJavamailEnd2end(t *testing.T) {
	mux, c := newTestClient(t)
	handleJSON(mux, "GET "+v1Path+"/config/javamail/end2ends/localhost",
		javamailEnd2EndJSON)
	result, err := c.GetJavamailEnd2end(t.Context(), "localhost")
	if err != nil {
		t.Fatal(err)
	}
	if result["readMailConfigName"] != "localhost" {
		t.Errorf("readMailConfigName = %v", result["readMailConfigName"])
	}
}

func TestCreateJavamailEnd2end(t *testing.T) {
	mux, c := newTestClient(t)
	req := handleStatus(mux, "POST "+v1Path+"/config/javamail/end2ends", 201)
	err := c.CreateJavamailEnd2end(t.Context(), map[string]any{
		"name":               "localhost",
		"readMailConfigName": "localhost",
		"sendMailConfigName": "localhost",
	})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal([]byte(req.body), &body); err != nil {
		t.Fatal(err)
	}
	if body["name"] != "localhost" {
		t.Errorf("name = %v", body["name"])
	}
}

func TestUpdateJavamailEnd2end(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "PUT "+v1Path+"/config/javamail/end2ends/localhost", 204)
	err := c.UpdateJavamailEnd2end(t.Context(), "localhost",
		map[string]string{"readMailConfigName": "remote"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestDeleteJavamailEnd2end(t *testing.T) {
	mux, c := newTestClient(t)
	handleStatus(mux, "DELETE "+v1Path+"/config/javamail/end2ends/localhost", 204)
	if err := c.DeleteJavamailEnd2end(t.Context(), "localhost"); err != nil {
		t.Fatal(err)
	}
}
