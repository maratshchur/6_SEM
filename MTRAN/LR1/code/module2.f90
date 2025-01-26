module module2
  implicit none


contains

  function sum_array(arr) result(sum)
    implicit none
    integer, dimension(:), intent(in) :: arr
    integer :: sum
    integer :: i
    sum = 0
    do i = 1, size(arr)
      sum = sum + arr(i)
    end do
  end function sum_array

  subroutine select_case_demo(value)
    implicit none
    integer, intent(in) :: value
    select case (value)
      case (1:5)
        print *, "Значение находится в диапазоне 1-5"
      case (6:10)
        print *, "Значение находится в диапазоне 6-10"
      case default
        print *, "Значение не находится в указанных диапазонах"
    end select
  end subroutine select_case_demo

  subroutine loop_control_example(n)
    implicit none
    integer, intent(in) :: n
    integer :: i
    do i = 1, n
      if (i == 5) then
        cycle
      end if
      if (i == 8) then
        exit
      end if
      print *, "Итерация с EXIT и CYCLE:", i
    end do
  end subroutine loop_control_example

end module module2