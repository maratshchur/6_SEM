module Geometry
    implicit none
    type :: Point
      integer :: age
      real :: height
    end type Point
contains
    subroutine do_loop_example(n)
      implicit none
      integer, intent(in) :: n
      integer :: i
      do while (i <= n)
        print *, "1. Cycle Iteration DO:", i
        n=i
      end do
      do i = 1, n
        print *, "2. Cycle Iteration DO:", i
      end do
      end subroutine do_loop_example
    function distance(name_in, age_in, height_in) result(person_out)
      integer, intent(in) :: age_in
      real, intent(in) :: height_in
      type(person) :: person_out
      person_out%name = name_in
      person_out%age = age_in
      person_out%height = height_in
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
    select case (age_in)
      case (1:5)
        print *, "Value is in 1-5 range"
      case (6:10)
        print *, "Value is in 6-10 range"
      case default
        print *, "Value is not presented in defined range"
    end select
    end function distance
end module Geometry
