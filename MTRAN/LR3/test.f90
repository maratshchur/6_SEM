module Geometry
    implicit none

    type :: Point
      real :: x, y
    end type Point

contains

    function distance(p1,) result(d)

      INTEGER :: a
      REAL :: b
      REAL :: d
      b=3.14
      a=3

      IF (a) THEN
        b=a
      END IF

      DO i = 1, 5
        distance = i 
      END DO
 
    end function distance

end module Geometry
