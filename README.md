# Layered Map

It is an algorithm that organizes data by its key, like some map structures.

The differential with `Layered Map` is that the time used to search for some key is equivalent, when not allocating more memory, to the size in bytes of the key.

It means that if you have one key called "four", it will take 4 times to access the value stored in "four".

It supports everything to be a key because the key is one generic slice of bytes, so you can use `strings`, `emojis 😊`, `json` and `media` files if you want, just converting it to slice of bytes.

# How it works

For each byte in the given key, the algorithm will create a way with linked lists to give access to the position or value in the last byte of the key, following the sequence of the key provided. The stored value can be anything, the user should know what is the interface that will handle the bytes.
Ex: I want to store "Hello, World!" with the key "abc"

```
[a] [ ] [ ]
 |___
[ ] [b] [ ]
     |___
[ ] [ ] [c] -> stores "Hello, World!"
```

Explaining: For each byte, if the next layer of chars isn't allocated, the algorithm create the next layer and try to store for the next, making it recursive while it dont reach the end of the key.

# Pros and Cons

This is usefull because if you have a lot of data oriented by documents, for example, and you have milions of data in your storage, with this algorithm you don't need to search in all your database if the document matches taking a lot of time, the time you will need to find the data will be the size of the document provided.

But, for the other side, if you have a lot of different keys, certainly you will need more memory to handle all the memory allocated.
