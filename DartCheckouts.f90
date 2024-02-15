program DartCheckouts
    implicit none
    integer, parameter :: DOUBLES(21) = [50, 40, 38, 36, 34, 32, 30, 28, 26, 24, 22, 20, 18, 16, 14, 12, 10, 8, 6, 4, 2]
    integer, parameter :: DARTS(43) = [60, 57, 54, 51, 50, 48, 45, 42, 40, 39, 38, 36, 34, 33, 32, 30, 28, 27, 26, 25, &
                                    24, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1]
    integer :: len_doubles, len_darts
    character(4), dimension(2) :: argv
    integer :: argc
    integer :: target_, total_darts, total_solutions, partial_sum
    real :: begin, end, elapsed_time
    integer :: begin2, end2, count_rate

    len_doubles = size(DOUBLES)
    len_darts = size(DARTS)
    
    argc = COMMAND_ARGUMENT_COUNT()

    if (argc == 2) then
        call GET_COMMAND_ARGUMENT(1, argv(1))
        read(argv(1),"(I4)") target_
        call GET_COMMAND_ARGUMENT(2, argv(2))
        read(argv(2),"(I4)") total_darts
    else
        target_ = 125
        total_darts = 3
    end if
    call print_dart()
    print *, "Finding ", total_darts, " dart finishes for ", target_

    call CPU_TIME(begin)
    call system_clock(begin2, count_rate)

    call findDartCheckouts()

    call CPU_TIME(end)
    call system_clock(end2)

    print *, "Total solutions: ", total_solutions
    elapsed_time = end - begin
    print *, "CPU time: ", elapsed_time, " seconds"
    print *, "system time: ", (end2 - begin2) / count_rate, " seconds"

    contains

    subroutine print_dart()
        print *, "    __________          "
        print *, "   /M\\\\\\M|||M//.     "
        print *, "  /MMM\\\\\\|||///M:.   "
        print *, " /MMMMM\\\\\\ | //MMMM:.            ______________________ "
        print *, "(=========+======<]]]]::::::::::<|||_|||_|||_|||_|||_|||>=========-"
        print *, " \\#MMMM// | \\\\MMMM:' "
        print *, "  \\#MM///|||\\\\\\M:'  "
        print *, "   \\M///M|||M\\\\'     "
        print *, ""
    end subroutine

    subroutine reduce_possible_darts()
        !! If the total score required is high we can skip low scoring darts that couldn't feasibly be used to reach the total.
        integer :: approx_lowest_normal_dart
        integer :: approx_lowest_finish_dart
        integer :: i

        approx_lowest_normal_dart = (target_ - 50) - ((total_darts - 2) * 60);
        approx_lowest_finish_dart = target_ - ((total_darts - 1) * 60);

        if (approx_lowest_normal_dart > 0) then
            do i = 1, len_darts
                if (DARTS(i) < approx_lowest_normal_dart) then
                    print *, "Capping possible darts at ", approx_lowest_normal_dart
                    len_darts = i - 1
                    exit
                end if
            end do
        end if

        if (approx_lowest_finish_dart > 0) then
            do i = 1, len_doubles
                if (DOUBLES(i) < approx_lowest_finish_dart) then
                    print *, "Capping possible finish dart at ", approx_lowest_finish_dart
                    len_doubles = i - 1;
                    return
                end if
            end do
        end if
    end subroutine

    recursive subroutine normal_darts(points, darts_, dart_count)
        !! normal darts are any dart which can be used when not checking out.
        integer, intent(in) :: points
        integer, intent(inout) :: darts_(:)
        integer, intent(in) :: dart_count
        integer :: j

        if (points == 0) then
            ! print *, "Checkout found: ", darts_(dart_count:1:-1)
            total_solutions = total_solutions + 1
            return
        end if

        if (points < 0 .OR. dart_count == total_darts) then
            ! Invalid state, return
            return
        end if

        do j = 1, len_darts
            darts_(dart_count+1) = DARTS(j)
            call normal_darts(points - DARTS(j), darts_, dart_count + 1)
        end do
    end subroutine normal_darts

    subroutine findDartCheckouts()
        !! Loop over the possible checkout darts, then recursively search all of the other darts to find possible combinations that reach the total.
        integer :: darts_(total_darts)
        integer :: i
        total_solutions = 0

        call reduce_possible_darts()

        !$OMP PARALLEL DO NUM_THREADS(len_doubles) PRIVATE(darts_)

        do i = 1, len_doubles
            print *, "==== Calculating for finishing dart ", DOUBLES(i), " ==="
            darts_(1) = DOUBLES(i)
            call normal_darts(target_ - DOUBLES(i), darts_, 1)
        end do

        !$OMP END PARALLEL DO

    end subroutine findDartCheckouts

end program DartCheckouts