#include <stdio.h>
#include <stdlib.h>
#include <time.h>
#include <omp.h>

int LEN_DOUBLES = 21;
int LEN_DARTS = 43;
const int DOUBLES[] = {50, 40, 38, 36, 34, 32, 30, 28, 26, 24, 22, 20, 18, 16, 14, 12, 10, 8, 6, 4, 2};
const int DARTS[] = {60, 57, 54, 51, 50, 48, 45, 42, 40, 39, 38, 36, 34, 33, 32, 30, 28, 27, 26, 25, 24, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1};

int target, total_darts, total_solutions;

void print_dart() {
    printf("    __________                                    \n");
    printf("   /M\\\\\\M|||M//.                                  \n");
    printf("  /MMM\\\\\\|||///M:.                                \n");
    printf(" /MMMMM\\\\\\ | //MMMM:.            ______________________ \n");
    printf("(=========+======<]]]]::::::::::<|||_|||_|||_|||_|||_|||>=========-\n");
    printf(" \\#MMMM// | \\\\MMMM:'                              \n");
    printf("  \\#MM///|||\\\\\\M:'                                 \n");
    printf("   \\M///M|||M\\\\'                                  \n\n");
}

void reduce_possible_darts(){
    // If the total score required is high we can skip
    // low scoring darts that couldn't feasibly be used
    // to reach the total.
    int approx_lowest_normal_dart = (target - 50) - ((total_darts - 2) * 60);
    int approx_lowest_finish_dart = target - ((total_darts - 1) * 60);
    if (approx_lowest_normal_dart > 0) {
        for (int i = 0; i < LEN_DARTS; i++){
            if (DARTS[i] < approx_lowest_normal_dart){
                printf("Capping possible darts at %d\n", approx_lowest_normal_dart);
                LEN_DARTS = i;
            }
        }
    }
    if (approx_lowest_finish_dart > 0) {
        for (int i = 0; i < LEN_DOUBLES; i++){
            if (DOUBLES[i] < approx_lowest_finish_dart){
                printf("Capping possible finish dart at %d\n", approx_lowest_finish_dart);
                LEN_DOUBLES = i;
            }
        }
    }
}

void print_checkout(int darts[], int dart_count){
    printf("Checkout found: ");
    for (int i=dart_count; i>0; i--){
        printf("%02d ", darts[i-1]);
    }
    printf("\n");
}

void normal_darts(int points, int darts[], int dart_count){
    // normal darts are any dart which can be used when not
    // checking out.
    if (points == 0){
        //print_checkout(darts, dart_count);
        ++total_solutions;
        return;
    }
    if (points < 0 || dart_count == total_darts) {
        // Invalid state, return
        return;
    }
    for (int i = 0; i < LEN_DARTS; i++) {
        darts[dart_count] = DARTS[i];
        normal_darts(points - DARTS[i], darts, dart_count + 1);
    }
}

int findDartCheckouts(){
    // Loop over the possible checkout darts, then recursively
    // search all of the other darts to find possible combinations
    // that reach the total
    int darts[total_darts];
    total_solutions = 0;

    reduce_possible_darts();

    #pragma omp parallel for num_threads(LEN_DOUBLES) private(darts)
    for (int i = 0; i < LEN_DOUBLES; i++){
        printf("=== Calculating for finishing dart %d ===\n", DOUBLES[i]);
        darts[0] = DOUBLES[i];
        normal_darts(target - DOUBLES[i], darts, 1);
    }
    return total_solutions;
}

int main(int argc, char *argv[]) {
    // takes args from stdin of target and then total_darts, otherwise sets a default value
    double itime, ftime, exec_time;

    if (argc == 3){
        target = atoi(argv[1]);
        total_darts = atoi(argv[2]);
    }
    else{
        target = 125;
        total_darts = 3;
    }
    print_dart();
    printf("Finding %d dart finishes for %d\n", total_darts, target);

    clock_t begin = clock();
    itime = omp_get_wtime();

    int total_solutions = findDartCheckouts();

    clock_t end = clock();
    ftime = omp_get_wtime();

    exec_time = ftime - itime;

    printf("Total solutions: %d\n", total_solutions);
    printf("CPU time: %.2fs\n", (double)(end - begin) / CLOCKS_PER_SEC);
    printf("Time taken is %f\n", exec_time);

    return 0;
}
