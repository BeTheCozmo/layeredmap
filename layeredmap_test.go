package layeredmap

import (
	"testing"
	"time"
)

func TestLayeredMapGeneric(t *testing.T) {
  lm := New()

  lm.Add([]byte("abc"), "Hello", nil)
  lm.Add([]byte("abc"), 42, nil)
  lm.Add([]byte("xyz"), true, nil)

  values, found := lm.GetAll([]byte("abc"))
  if !found || len(values) != 2 {
    t.Errorf("Esperado 2 valores, obtido %d", len(values))
  }
  if values[0] != "Hello" || values[1] != 42 {
    t.Errorf("Esperado ['Hello', 42], obtido %v", values)
  }

  last, found := lm.GetLast([]byte("abc"))
  if !found || last != 42 {
    t.Errorf("Esperado 42, obtido %v", last)
  }

  boolResult, found := lm.GetLast([]byte("xyz"))
  if !found || boolResult != true {
    t.Errorf("Esperado true, obtido %v", boolResult)
  }
}

func TestLayeredMapPop(t *testing.T) {
  lm := New()

  lm.Add([]byte("abc"), "Hello", nil)
  lm.Add([]byte("abc"), "World", nil)
  lm.Add([]byte("abc"), 42, nil)

  value, found := lm.PopLast([]byte("abc"))
  if !found || value != 42 {
    t.Errorf("PopLast: esperado 42, obtido %v", value)
  }
  values, _ := lm.GetAll([]byte("abc"))
  if len(values) != 2 {
    t.Errorf("Esperado 2 valores após PopLast, obtido %d", len(values))
  }

  value, found = lm.PopFirst([]byte("abc"))
  if !found || value != "Hello" {
    t.Errorf("PopFirst: esperado 'Hello', obtido %v", value)
  }
  values, _ = lm.GetAll([]byte("abc"))
  if len(values) != 1 || values[0] != "World" {
    t.Errorf("Esperado ['World'] após PopFirst, obtido %v", values)
  }

  value, found = lm.PopLast([]byte("xyz"))
  if found || value != nil {
    t.Errorf("PopLast em chave inexistente deveria retornar nil, obtido %v", value)
  }
}

func TestLayeredMapDeque(t *testing.T) {
  lm := New()

  lm.Add([]byte("abc"), "Hello", nil)
  lm.Add([]byte("abc"), "World", nil)
  lm.Add([]byte("abc"), 42, nil)

  value, found := lm.PopLast([]byte("abc"))
  if !found || value != 42 {
    t.Errorf("PopLast: esperado 42, obtido %v", value)
  }
  values, _ := lm.GetAll([]byte("abc"))
  if len(values) != 2 || values[1] != "World" {
    t.Errorf("Esperado ['Hello', 'World'] após PopLast, obtido %v", values)
  }

  value, found = lm.PopFirst([]byte("abc"))
  if !found || value != "Hello" {
    t.Errorf("PopFirst: esperado 'Hello', obtido %v", value)
  }
  values, _ = lm.GetAll([]byte("abc"))
  if len(values) != 1 || values[0] != "World" {
    t.Errorf("Esperado ['World'] após PopFirst, obtido %v", values)
  }

  last, found := lm.GetLast([]byte("abc"))
  if !found || last != "World" {
    t.Errorf("GetLast: esperado 'World', obtido %v", last)
  }
}

func TestLayeredMapExpiration(t *testing.T) {
  lm := New()

  // Adiciona valores com e sem expiração
  ttlShort := time.Millisecond * 100 // 100ms
  ttlLong := time.Second * 10        // 10s
  lm.Add([]byte("abc"), "Hello", &ttlShort)
  lm.Add([]byte("abc"), "World", &ttlLong)
  lm.Add([]byte("abc"), "Forever", nil) // Sem expiração

  // Espera o primeiro valor expirar
  time.Sleep(time.Millisecond * 150)

  // Testa GetAll
  values, found := lm.GetAll([]byte("abc"))
  if !found || len(values) != 2 || values[0] != "World" || values[1] != "Forever" {
      t.Errorf("Esperado ['World', 'Forever'], obtido %v", values)
  }

  // Testa GetLast
  last, found := lm.GetLast([]byte("abc"))
  if !found || last != "Forever" {
      t.Errorf("GetLast: esperado 'Forever', obtido %v", last)
  }

  // Testa PopFirst
  value, found := lm.PopFirst([]byte("abc"))
  if !found || value != "World" {
      t.Errorf("PopFirst: esperado 'World', obtido %v", value)
  }

  // Testa PopLast
  value, found = lm.PopLast([]byte("abc"))
  if !found || value != "Forever" {
      t.Errorf("PopLast: esperado 'Forever', obtido %v", value)
  }
}