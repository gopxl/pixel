# IMDraw

<hr>
 IMDraw is an immediate-mode-like shape drawer and BasicTarget. IMDraw supports TrianglesPosition,
 TrianglesColor, TrianglesPicture and PictureColor.

 IMDraw, other than a regular BasicTarget, is used to draw shapes. To draw shapes, you first need
 to Push some points to IMDraw:
```go
   imd := pixel.NewIMDraw(pic)  use nil pic if you only want to draw primitive shapes
   imd.Push(pixel.V(100, 100))
   imd.Push(pixel.V(500, 100))
```
 Once you have Pushed some points, you can use them to draw a shape, such as a line:

   `imd.Line(20)  //draws a 20 units thick line``

 Set exported fields to change properties of Pushed points:

```go
   imd.Color = pixel.RGB(1, 0, 0)
   imd.Push(pixel.V(200, 200))
   imd.Circle(400, 0)
```
 Here is the list of all available point properties (need to be set before Pushing a point):
   - Color     - applies to all
   - Picture   - coordinates, only applies to filled polygons
   - Intensity - picture intensity, only applies to filled polygons
   - Precision - curve drawing precision, only applies to circles and ellipses
   - EndShape  - shape of the end of a line, only applies to lines and outlines

 And here's the list of all shapes that can be drawn (all, except for line, can be filled or
 outlined):
   - Line
   - Polygon
   - Circle
   - Circle arc
   - Ellipse
   - Ellipse arc