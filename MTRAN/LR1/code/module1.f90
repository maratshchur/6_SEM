module module1
  implicit none

  type person
    character(len=50) :: name
    integer :: age
    real :: height
  end type person

contains

  function create_person(name_in, age_in, height_in) result(person_out)
    implicit none
    character(len=*), intent(in) :: name_in
    integer, intent(in) :: age_in
    real, intent(in) :: height_in
    type(person) :: person_out
    person_out%name = name_in
    person_out%age = age_in
    person_out%height = height_in
  end function create_person

  subroutine print_person_info(person_in)
    implicit none
    type(person), intent(in) :: person_in
    print *, "Имя:", person_in%name
    print *, "Возраст:", person_in%age
    print *, "Рост:", person_in%height
  end subroutine print_person_info

  subroutine do_loop_example(n)
    implicit none
    integer, intent(in) :: n
    integer :: i
    do i = 1, n
      print *, "Итерация цикла DO:", i
    end do
  end subroutine do_loop_example

  subroutine do_while_loop_example(n)
    implicit none
    integer, intent(in) :: n
    integer :: i = 1
    do while (i <= n)
      print *, "Итерация цикла DO WHILE:", i
      i = i + 1
    end do
  end subroutine do_while_loop_example

end module module1