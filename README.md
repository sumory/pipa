## pipa

pipe anything.

```
                              XXXXX X    XXXXX    X
                              X   X X    X   X   X X
                              XXXXX X    XXXXX  X X X
                              X     X    X      X   X
                              X     X    X      X   X
              +-------------------------------------------------------+
              |                                                       |
+--------+    |                                                       |    +--------+
|        |    |                   +------+   +------+                 |    |        |
|  file  |    |    +--------+     |      |   |      |     +------+    |    |  file  |
| string +--> |    | source +---> | pipe +-> | pipe +---> | sink |    +--> | string |
|  etc.  |    |    +--------+     |      |   |      |     +------+    |    |  etc.  |
|        |    |                   +------+   +------+                 |    |        |
+--------+    |                                                       |    +--------+
              |                                                       |
              +-------------------------------------------------------+
```


 - inspired by [flume](http://flume.apache.org/)
 - still under heavy development, use at your own risk

### TODO

 - something isn't reasonable in the origin design. the `pipe`, it's not like a real pipe.
 - `pipa.go` should only define the interfaces and design the invoke code.
 - `channel` may be omitted in some cases

