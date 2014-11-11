package clipboard

import (
	"fmt"
	"gopkg.in/qml.v1"
	"runtime"
  "time"
	"testing"
)

func setup() {
	qml.SetupTesting()
	qml.SetLogger(nil)
	qml.CollectStats(true)
	qml.ResetStats()
	stats := qml.Stats()
	if stats.EnginesAlive > 0 || stats.ValuesAlive > 0 || stats.ConnectionsAlive > 0 {
		panic(fmt.Sprintf("Test started with values alive: %#v\n", stats))
	}
}

func teardown() {
	retries := 30 // Three seconds top.
	for {
		// Do not call qml.Flush here. It creates a nested event loop
		// that attempts to process the deferred object deletes and cannot,
		// because deferred deletes are only processed at the same loop level.
		// So it *reposts* the deferred deletion event, in practice *preventing*
		// these objects from being deleted.
		runtime.GC()
		stats := qml.Stats()
		if stats.EnginesAlive == 0 && stats.ValuesAlive == 0 && stats.ConnectionsAlive == 0 {
			break
		}
		if retries == 0 {
			panic(fmt.Sprintf("there are values alive:\n%#v\n", stats))
		}
		retries--
		time.Sleep(100 * time.Millisecond)
		if retries%10 == 0 {
			fmt.Printf("There are still objects alive; waiting for them to die: %#v\n", stats)
		}
	}
	qml.SetLogger(nil)
}

func TestCopyAndPaste(t *testing.T) {
	setup()
	defer teardown()
	engine := qml.NewEngine()
	
	clip := New(engine)
	defer engine.Destroy()

	//======================
	a := "Привет, мир!"

	err := clip.WriteAll(a)
	if err != nil {
		t.Fatal(err)
	}

	b, err := clip.ReadAll()
	if err != nil {
		t.Fatal(err)
	}

	if a != b {
		t.Errorf("expected '%s', got '%s'", a, b)
	}
}

// func BenchmarkReadAll(b *testing.B) {
//   //qml.Init(nil)
//   engine := qml.NewEngine()
//   clip := New(engine)
//   b.ResetTimer()
//   for i := 0; i < b.N; i++ {
//     clip.ReadAll()
//   }
// }
//
// func BenchmarkWriteAll(b *testing.B) {
//   //qml.Init(nil)
//   engine := qml.NewEngine()
//   clip := New(engine)
//   b.ResetTimer()
//   text := "Жарю сосиски"
//   for i := 0; i < b.N; i++ {
//     clip.WriteAll(text)
//   }
// }
