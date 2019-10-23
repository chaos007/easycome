package config

import (
	"fmt"
	"reflect"
	"testing"
)

func TestInit(t *testing.T) {
	_, err := DesignConfigInit("../../config")
	if err != nil {
		fmt.Println("err:", err.Error())
	}

	m, _ := design.Load("item")
	list := m.(Maper).ForEach(func(p Parser) interface{} {
		return p
	})

	for _, item := range list {
		fmt.Println("----1:", item)
	}

	// UploadDesignConfig("/item/20000", `{"id":"20000","name":"dff"}`)

	// list = m.(Maper).ForEach(func(p Parser) interface{} {
	// 	return p
	// })

	// for _, item := range list {
	// 	fmt.Println("----2:", item)
	// }

}

func dddd(list interface{}) {
}

func TestConfig(t *testing.T) {
	// value2 := reflect.ValueOf(designConfig)
	// fmt.Println(value2.NumField())

	// value := reflect.TypeOf(designConfig)
	// num := value.Elem().NumField()
	// for i := 0; i < num; i++ {
	// 	method, ok := value.Elem().Field(i).Type.MethodByName("Init")
	// 	if !ok {
	// 		return
	// 	}
	// 	m := reflect.ValueOf(&TestMap{})
	// 	fmt.Println("----1:", value.Elem().Field(i).Type)
	// 	fmt.Println("----2:", method.Name)
	// 	fmt.Println("----3:", method.Func.Call([]reflect.Value{m}))
	// }

}

func testvalue(v interface{}) {
	ty := reflect.TypeOf(v)
	fmt.Println("-------ty:", ty.Kind())

}

func removeID(tmap map[int64]*Test, c chan int) {
	delete(tmap, 1)
	c <- 1
}

func changeID(t *Test, c chan int) {
	select {
	case <-c:
		t.ID = 2
	}
}
