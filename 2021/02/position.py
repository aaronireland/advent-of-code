import re
import sys


def read_file(file_path):
    movements = []
    with open(file_path, 'r') as input:
        for move in input.readlines():
            parts = move.split()
            if len(parts) > 1:
                if parts[1].isdigit():
                    scalar = int(parts[1])
                    if parts[0] == 'forward':
                        movements.append((scalar, 0,))
                    elif parts[0] == 'up':
                        movements.append((0, -1 * scalar,))
                    elif parts[0] == 'down':
                        movements.append((0, scalar,))
    return movements


def calculate_distance(movements, aimed=False):
    pos_x, pos_z = 0, 0
    aim = 0
    for x, z in movements:
        pos_x += x
        if aimed:
            aim += z
            if x != 0:
                pos_z += aim * x
        else:
            pos_z += z

    return pos_x * pos_z


if __name__ == '__main__':
    aimed = False
    for i, arg in enumerate(sys.argv):
        if i > 0:
            if arg == "aimed":
                aimed = True
    movements = read_file('input.txt')
    dist = calculate_distance(movements, aimed)
    print(dist)
