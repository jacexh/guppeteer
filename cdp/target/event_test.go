package target

import "testing"

func TestEvent(t *testing.T) {
	ev := &EventReceivedMessageFromTarget{}
	if ev.Name() != ReceivedMessageFromTarget {
		t.FailNow()
	}
	if ev.Domain() != "Target" {
		t.FailNow()
	}
	data := []byte(`{"sessionId":"(D2D0CBBEEB35E321B7C8004757E483AC):1","message":"{\"method\":\"Network.dataReceived\",\"params\":{\"requestId\":\"15761.10\",\"timestamp\":3005.34426,\"dataLength\":31623,\"encodedDataLength\":32844}}","targetId":"(D2D0CBBEEB35E321B7C8004757E483AC)"}`)
	params, err := ev.Load(data)
	if err != nil {
		t.FailNow()
	}
	val, ok := params.(*ReceivedMessageFromTargetParams)
	if !ok {
		t.FailNow()
	}
	if val.SessionID != "(D2D0CBBEEB35E321B7C8004757E483AC):1" || val.TargetID != "(D2D0CBBEEB35E321B7C8004757E483AC)" {
		t.FailNow()
	}
}
