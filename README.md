## pipa

pipe anything.

```

                                 XXXXX X    XXXXX    X
                                 X   X X    X   X   X X
                                 XXXXX X    XXXXX  X X X
                                 X     X    X     X     X
                                 X     X    X     X     X
                 +-------------------------------------------------------+
                 |                                                       |
   +--------+    |                                                       |    +--------+
   |        |    |                     +-----------+                     |    |        |
   |  file  |    |     +--------+      |           |      +--------+     |    |  file  |
   | string +--> |     | source +----> |  channel  +----> |  sink  |     +--> | string |
   |  etc.  |    |     +--------+      |           |      +--------+     |    |  etc.  |
   |        |    |                     +-----------+                     |    |        |
   +--------+    |                                                       |    +--------+
                 |                                                       |
                 +-------------------------------------------------------+

```


 - inspired by [flume](http://flume.apache.org/)
 - still under heavy development, use at your own risk


