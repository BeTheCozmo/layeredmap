package layeredmap

import (
	"container/list"
	"time"
)

type valueWithExpiry struct {
  value interface{}
  expiresAt *time.Time
}

type Node struct {
  next   map[byte]*Node
  values *list.List
}

type LayeredMap struct {
  root *Node
}

func New() *LayeredMap {
  return &LayeredMap{
    root: &Node{
      next:   make(map[byte]*Node),
      values: list.New(),
    },
  }
}

func (lm *LayeredMap) Add(key []byte, value interface{}, ttl *time.Duration) {
  current := lm.root
  for _, b := range key {
    if current.next[b] == nil {
      current.next[b] = &Node{
        next:   make(map[byte]*Node),
        values: list.New(),
      }
    }
    current = current.next[b]
  }
  var expiresAt *time.Time
  if ttl != nil {
    t := time.Now().Add(*ttl)
    expiresAt = &t
  }
  current.values.PushBack(&valueWithExpiry{value, expiresAt})
}

func (lm *LayeredMap) GetAll(key []byte) ([]interface{}, bool) {
  current := lm.root
  for _, b := range key {
    if current.next[b] == nil {
      return []interface{}{}, false
    }
    current = current.next[b]
  }
  if current.values.Len() == 0 {
    return []interface{}{}, false
  }

  validValues := make([]interface{}, 0, current.values.Len())
  for e := current.values.Front(); e != nil; {
    next := e.Next()
    v := e.Value.(*valueWithExpiry)
    if v.expiresAt != nil && time.Now().After(*v.expiresAt) {
      current.values.Remove(e)
    } else {
      validValues = append(validValues, v.value)
    }
    e = next
  }
  if len(validValues) == 0 {
    return []interface{}{}, false
  }
  return validValues, true
}

func (lm *LayeredMap) GetLast(key []byte) (interface{}, bool) {
  current := lm.root
  for _, b := range key {
    if current.next[b] == nil {
      return nil, false
    }
    current = current.next[b]
  }
  if current.values.Len() == 0 {
    return nil, false
  }

  for e := current.values.Back(); e != nil; e = e.Prev() {
    v := e.Value.(*valueWithExpiry)
    if v.expiresAt == nil || !time.Now().After(*v.expiresAt) {
      return v.value, true
    }
    current.values.Remove(e)
  }
  return nil, false
}

func (lm *LayeredMap) PopLast(key []byte) (interface{}, bool) {
  current := lm.root
  for _, b := range key {
    if current.next[b] == nil {
      return nil, false
    }
    current = current.next[b]
  }
  if current.values.Len() == 0 {
    return nil, false
  }
  for e := current.values.Back(); e != nil; e = e.Prev() {
    v := e.Value.(*valueWithExpiry)
    if v.expiresAt == nil || !time.Now().After(*v.expiresAt) {
      current.values.Remove(e)
      return v.value, true
    }
    current.values.Remove(e)
  }
  return nil, false
}

func (lm *LayeredMap) PopFirst(key []byte) (interface{}, bool) {
  current := lm.root
  for _, b := range key {
    if current.next[b] == nil {
      return nil, false
    }
    current = current.next[b]
  }
  if current.values.Len() == 0 {
    return nil, false
  }
  for e := current.values.Front(); e != nil; e = e.Next() {
    v := e.Value.(*valueWithExpiry)
    if v.expiresAt == nil || !time.Now().After(*v.expiresAt) {
      current.values.Remove(e)
      return v.value, true
    }
    current.values.Remove(e)
  }
  return nil, false
}