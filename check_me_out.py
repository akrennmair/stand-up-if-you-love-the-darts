"""

Darts Calculator

"""
import sys
from time import time

BULL = 50
OUTER = 25

def min_checkout(target: int, num_darts: int) -> int:
    """Approximate the minimum dart that can be used to checkout given
    the target score"""
    return max(2, target - ((num_darts - 1) * 60))

def min_normal_dart(target: int, num_darts: int) -> int:
    """Approximate the minimum non-checkout dart that can be used given
    the target score"""
    return max(1, target - BULL - (num_darts - 2) * 60)

def check_me_out(target: int, max_darts: int) -> list[list[int]]:
    """
    Return all possible combinations to check out target value
    with in max_darts
    """
    singles = list(reversed(range(1, 21)))
    doubles = list(reversed(range(2, 41, 2)))
    triples = list(reversed(range(3, 61, 3)))
    possibles = [BULL, OUTER] + list(set(singles + doubles + triples))[::-1]

    # Cut off some possibles
    doubles = [dart for dart in doubles if dart >= min_checkout(target, max_darts)]
    possibles = [dart for dart in possibles if dart >= min_normal_dart(target, max_darts)]

    def lets_play_darts(current: list[int] | None=None, dart: int=0):
        """
        Recusively search for dart checkouts

        `current` is a list of darts used so far
        `dart` is the number of darts used so far
        `total` is the sum of current
        """
        if current is None:
            for score in [BULL] + doubles:
                print(f"Looking for scores finishing on {score}")
                lets_play_darts([score], dart + 1)
        elif (total := sum(current)) == target:
            checkouts.append(current)
            return
        elif total > target or dart == max_darts:
            return
        else:
            for score in possibles:
                lets_play_darts(current + [score], dart+1)

    checkouts = []
    lets_play_darts()
    return checkouts

if __name__ == "__main__":
    if len(sys.argv) == 3:
        target = int(sys.argv[1])
        num_darts = int(sys.argv[2])
    else:
        target, num_darts = 125, 3
    print("""            __________                                    
          /M\\\M|||M//.                                  
          /MMM\\\|||///M:.                                
          /MMMMM\\ | //MMMM:.              ______________________
          (=========+======<]]]]::::::::::<|||_|||_|||_|||_|||_|||>=========-
          \#MMMM// | \\MMMM:'                              
          \#MM///|||\\\M:'                                 
          \M///M|||M\\'                                  
          """)
    print(f"Finding {num_darts} dart finishes for {target}")
    start = time()
    checkouts = check_me_out(target, max_darts=num_darts)
    end = time()

    print(f"Found {len(checkouts)} checkouts in {end-start:.2f} seconds")
