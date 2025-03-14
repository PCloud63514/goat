package environment

import (
	"reflect"
	"testing"
)

func TestEnvironment_New(t *testing.T) {

	env := New()
	if env == nil {
		t.Fatal("인스턴스 생성에 실패하였습니다.")
	}
}

func TestEnvironment_GetRequiredProperty(t *testing.T) {
	env := New(Option{
		ResPath:  ".",
		Profiles: []string{"test"},
	})
	t.Run("키가 없을 경우 에러를 반환합니다.", func(t *testing.T) {
		// when
		_, err := env.GetRequiredProperty("NOT_EXIST")
		// then
		if err == nil {
			t.Errorf("키가 없을 경우 에러가 발생해야합니다.")
		}
	})
	t.Run("결과를 반환합니다.", func(t *testing.T) {
		v, err := env.GetRequiredProperty("test.value.string")
		if err != nil {
			t.Errorf("키가 있을 경우 결과를 반환해야합니다. %s", err)
		}
		if v != "str" {
			t.Errorf("test.properties의 [key=test.value.string,value=str] 값과 동일하지 않습니다. \nExpected: %v\nActual: %v", "str", v)
		}
	})
}

func TestEnvironment_GetRequiredPropertyInt(t *testing.T) {
	env := New(Option{
		ResPath:  ".",
		Profiles: []string{"test"},
	})
	t.Run("키가 없을 경우 에러를 반환합니다.", func(t *testing.T) {
		// when
		_, err := env.GetRequiredPropertyInt("NOT_EXIST")
		// then
		if err == nil {
			t.Errorf("키가 없을 경우 에러가 발생해야합니다.")
		}
	})
	t.Run("숫자가 아닌 값을 반환할 경우 0을 반환합니다.", func(t *testing.T) {
		v, err := env.GetRequiredPropertyInt("test.value.string")
		if err != nil {
			t.Errorf("키가 있을 경우 결과를 반환해야합니다. %s", err)
		}
		if v != 0 {
			t.Errorf("숫자가 아닐 경우 0을 반환해야합니다. %s", err)
		}
	})
	t.Run("결과를 반환합니다.", func(t *testing.T) {
		v, err := env.GetRequiredPropertyInt("test.value.int")
		if err != nil {
			t.Errorf("키가 있을 경우 결과를 반환해야합니다. %s", err)
		}
		if v != 100 {
			t.Errorf("test.properties의 [key=test.value.int,value=100] 값과 동일하지 않습니다. \nExpected: %v\nActual: %v", 100, v)
		}
	})
}

func TestEnvironment_GetRequiredPropertyBool(t *testing.T) {
	env := New(Option{
		ResPath:  ".",
		Profiles: []string{"test"},
	})
	t.Run("키가 없을 경우 에러를 반환합니다.", func(t *testing.T) {
		// when
		_, err := env.GetRequiredPropertyBool("NOT_EXIST")
		// then
		if err == nil {
			t.Errorf("키가 없을 경우 에러가 발생해야합니다.")
		}
	})
	t.Run("bool이 아닌 값을 반환할 경우 false을 반환합니다.", func(t *testing.T) {
		v, err := env.GetRequiredPropertyBool("test.value.string")
		if err != nil {
			t.Errorf("키가 있을 경우 결과를 반환해야합니다. %s", err)
		}
		if v != false {
			t.Errorf("bool이 아닐 경우 false을 반환해야합니다. %s", err)
		}
	})
	t.Run("결과를 반환합니다.", func(t *testing.T) {
		v, err := env.GetRequiredPropertyBool("test.value.bool")
		if err != nil {
			t.Errorf("키가 있을 경우 결과를 반환해야합니다. %s", err)
		}
		if v != true {
			t.Errorf("test.properties의 [key=test.value.bool,value=true] 값과 동일하지 않습니다. \nExpected: %v\nActual: %v", true, v)
		}
	})
}

func TestEnvironment_GetRequiredPropertySlice(t *testing.T) {
	env := New(Option{
		ResPath:  ".",
		Profiles: []string{"test"},
	})
	t.Run("키가 없을 경우 에러를 반환합니다.", func(t *testing.T) {
		// when
		_, err := env.GetRequiredPropertySlice("NOT_EXIST")
		// then
		if err == nil {
			t.Errorf("키가 없을 경우 에러가 발생해야합니다.")
		}
	})
	t.Run("결과를 반환합니다.", func(t *testing.T) {
		givenValue := []string{"a", "b", "c", "d", "e", "f", "g"}
		v, err := env.GetRequiredPropertySlice("test.value.slice")
		if err != nil {
			t.Errorf("키가 있을 경우 결과를 반환해야합니다. %s", err)
		}
		if !reflect.DeepEqual(v, givenValue) {
			t.Errorf("test.properties의 [key=test.value.list,value=[a,b,c,d,e,f,g]] 값과 동일하지 않습니다. \nExpected: %v\nActual: %v", givenValue, v)
		}
	})
}

func TestEnvironment_ContainsProperty(t *testing.T) {
	env := New(Option{
		ResPath:  ".",
		Profiles: []string{"test"},
	})
	t.Run("키가 없을 경우 false를 반환합니다.", func(t *testing.T) {
		// when
		v := env.ContainsProperty("NOT_EXIST")
		// then
		if v {
			t.Errorf("키가 없을 경우 false를 반환해야합니다.")
		}
	})
	t.Run("키가 있을 경우 true를 반환합니다.", func(t *testing.T) {
		v := env.ContainsProperty("test.value.string")
		if !v {
			t.Errorf("키가 있을 경우 true를 반환해야합니다.")
		}
	})
}
