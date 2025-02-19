program main
  use module1
  use module2
  implicit none

  type(person) :: my_person
  integer :: array_sum

  my_person = create_person("Иван Иванов", 30, 1.85)
  call print_person_info(my_person)

  array_sum = sum_array(numbers)
  print *, "Сумма элементов массива:", array_sum
  call select_case_demo(7)

  call do_loop_example(5)
  call do_while_loop_example(3)
  call loop_control_example(10)

end program main


!gfortran -c module1.f90 module2.f90  && gfortran -o main main.f90 module1.o module2.o