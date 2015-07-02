package util

import (
	"errors"
	"testing"
	"time"
)

func TestGenerateRandomECSPassword(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := GenerateRandomECSPassword()

		if len(s) < 8 || len(s) > 30 {
			t.Errorf("Generated ECS password [%v]: bad len", s)
		}

		hasDigit := false
		hasLower := false
		hasUpper := false

		for j := range s {

			switch {
			case '0' <= s[j] && s[j] <= '9':
				hasDigit = true
			case 'a' <= s[j] && s[j] <= 'z':
				hasLower = true
			case 'A' <= s[j] && s[j] <= 'Z':
				hasUpper = true
			}
		}

		if !hasDigit {
			t.Errorf("Generated ECS password [%v]: no digit", s)
		}

		if !hasLower {
			t.Errorf("Generated ECS password [%v]: no lower letter ", s)
		}

		if !hasUpper {
			t.Errorf("Generated ECS password [%v]: no upper letter", s)
		}
	}
}

func TestSimpleLoopCal(t *testing.T) {

	attempts := AttemptStrategy{
		Min:   5,
		Total: 5 * time.Second,
		Delay: 200 * time.Millisecond,
	}

	timeoutFunc := func() bool {
		return false
	}

	successFunc := func() bool {
		return true
	}

	begin := time.Now()

	err := SimpleLoopCall(attempts, timeoutFunc)

	if err != nil {
		t.Logf("timeout func complete successful")
	} else {
		t.Error("Expect timeout result")
	}

	end := time.Now()
	duration := end.Sub(begin).Seconds()

	t.Logf("loop call duration: %f", duration)

	if duration > (float64(attempts.Min) - 1) {
		t.Logf("timeout func duration is enough")
	} else {
		t.Error("timeout func duration is not enough")
	}

	err = SimpleLoopCall(attempts, successFunc)

	if err != nil {
		t.Error("Expect success result")
	} else {
		t.Logf("success func complete successful")
	}
}

func TestLoopCall(t *testing.T) {

	attempts := AttemptStrategy{
		Min:   5,
		Total: 5 * time.Second,
		Delay: 200 * time.Millisecond,
	}

	timeoutFunc := func() (bool, interface{}, error) {
		return false, StatusNotAvailable, nil
	}

	begin := time.Now()

	_, timeoutError := LoopCall(attempts, timeoutFunc)
	if timeoutError != nil {
		t.Logf("timeout func complete successful")
	} else {
		t.Error("Expect timeout result")
	}

	end := time.Now()
	duration := end.Sub(begin).Seconds()
	if duration > (float64(attempts.Min) - 1) {
		t.Logf("timeout func duration is enough")
	} else {
		t.Error("timeout func duration is not enough")
	}

	errorFunc := func() (bool, interface{}, error) {
		err := errors.New("execution failed")
		return false, StatusNotAvailable, err
	}

	_, failedError := LoopCall(attempts, errorFunc)
	if failedError != nil {
		t.Logf("error func complete successful: " + failedError.Error())
	} else {
		t.Error("Expect error result")
	}

	successFunc := func() (bool, interface{}, error) {
		return true, StatusNotAvailable, nil
	}

	_, successError := LoopCall(attempts, successFunc)
	if successError != nil {
		t.Error("Expect success result")
	} else {
		t.Logf("success func complete successful")
	}
}
