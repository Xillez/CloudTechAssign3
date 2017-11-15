package main

import "testing"

func Test_keepAlive(t *testing.T){
	isWorking := false
	keepAlive(&isWorking, false)
	if !isWorking{
		t.Error("Webhooks were not updated")
	}
}

func Test_updateWebhooks(t *testing.T){
	isWorking := false
	updateWebhooks(&isWorking, false)
	if !isWorking{
		t.Error("Webhooks were not updated")
	}
}